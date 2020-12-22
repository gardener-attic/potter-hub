/*
Copyright (c) 2018 Bitnami

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kubeapps/common/response"
	"github.com/pkg/errors"
	"github.com/urfave/negroni"
	"helm.sh/helm/v3/pkg/chart"

	"github.wdf.sap.corp/kubernetes/hub/pkg/auth"
	chartUtils "github.wdf.sap.corp/kubernetes/hub/pkg/chart"
	errorUtils "github.wdf.sap.corp/kubernetes/hub/pkg/errors"
	"github.wdf.sap.corp/kubernetes/hub/pkg/kubeval"
	logUtils "github.wdf.sap.corp/kubernetes/hub/pkg/log"
	"github.wdf.sap.corp/kubernetes/hub/pkg/proxy"
	utils "github.wdf.sap.corp/kubernetes/hub/pkg/util"
)

// userKey is the context key for the user data in the request context
type userKey struct{}
type validationObjectKey struct{}

// AuthGate implements middleware to check if the user is logged in before continuing
func TokenAuthorization() negroni.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		token, err := utils.GetTokenFromRequest(req)
		if err != nil {
			utils.SendErrResponse(req.Context(), w, err)
			return
		}
		ctx := context.WithValue(req.Context(), validationObjectKey{}, proxy.TokenValidation{Token: token})
		userAuth, err := auth.NewAuth(token)
		if err != nil {
			wrappedErr := errors.Wrap(err, "")
			utils.SendErrResponse(req.Context(), w, wrappedErr)
			return
		}
		err = userAuth.Validate()
		if err != nil {
			wrappedErr := errorUtils.Unauthorized.New(errors.New(err.Error()))
			utils.SendErrResponse(req.Context(), w, wrappedErr)
			return
		}
		ctx = context.WithValue(ctx, userKey{}, userAuth)
		next(w, req.WithContext(ctx))
	}
}

func KubeconfigAuthorization(oidcClusterURL string, decodedOidcClusterCA []byte) negroni.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		var byteKube []byte
		var disableAuth bool
		var err error

		disableAuthHeader := req.Header.Get("disableAuth")
		if disableAuthHeader != "" {
			disableAuth, err = strconv.ParseBool(disableAuthHeader)
			if err != nil {
				wrappedErr := errors.New("Disable auth header could not be parsed")
				utils.SendErrResponse(req.Context(), w, wrappedErr)
				return
			}
		}

		if disableAuth {
			encodedKubeconfig := req.Header.Get("targetKubeconfig")
			if encodedKubeconfig == "" {
				resErr := errorUtils.BadRequest.New(errors.New("Kubeconfig header is missing"))
				utils.SendErrResponse(req.Context(), w, resErr)
				return
			}

			byteKube, err = base64.StdEncoding.DecodeString(encodedKubeconfig)
			if err != nil {
				err := errors.New("Kubeconfig could not be base64 decoded")
				utils.SendErrResponse(req.Context(), w, err)
				return
			}
		} else {
			token, err := utils.GetTokenFromRequest(req)
			if err != nil {
				err = errorUtils.Unauthorized.New(err)
				utils.SendErrResponse(req.Context(), w, err)
				return
			}

			params := mux.Vars(req)

			namespace := params["clusterNamespace"]
			accessData := params["accessData"]

			if namespace == "" || accessData == "" {
				err = errorUtils.BadRequest.New(errors.New("URL parameters namespace and secretName are missing"))
				utils.SendErrResponse(req.Context(), w, err)
				return
			}

			kubeconfig, err := kubeval.GetKubeconfigFromOidcCluster(token, namespace, accessData, oidcClusterURL, decodedOidcClusterCA)
			if err != nil {
				utils.SendErrResponse(req.Context(), w, err)
				return
			}
			byteKube = []byte(*kubeconfig)
		}

		ctx := context.WithValue(req.Context(), validationObjectKey{}, proxy.KubeconfigValidation{Kubeconfig: byteKube})

		next(w, req.WithContext(ctx))
	}
}

// Params a key-value map of path params
type Params map[string]string

// WithParams can be used to wrap handlers to take an extra arg for path params
type WithParams func(http.ResponseWriter, *http.Request, Params)

func (h WithParams) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	h(w, req, vars)
}

// WithoutParams can be used to wrap handlers that doesn't take params
type WithoutParams func(http.ResponseWriter, *http.Request)

func (h WithoutParams) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h(w, req)
}

func isNotFound(err error) bool {
	errorAsLowerCase := strings.ToLower(err.Error())
	return strings.Contains(errorAsLowerCase, "not found")
}

func isAlreadyExists(err error) bool {
	errorAsLowerCase := strings.ToLower(err.Error())
	return strings.Contains(errorAsLowerCase, "is still in use") || strings.Contains(errorAsLowerCase, "already exists")
}

func isForbidden(err error) bool {
	errorAsLowerCase := strings.ToLower(err.Error())
	return strings.Contains(errorAsLowerCase, "unauthorized")
}

func isUnprocessable(err error) bool {
	errorAsLowerCase := strings.ToLower(err.Error())
	re := regexp.MustCompile(`release.*failed`)
	return re.MatchString(errorAsLowerCase)
}

func errorCode(err error) error {
	return errorCodeWithDefault(err, errorUtils.InternalServerError)
}

func errorCodeWithDefault(err error, defaultCode errorUtils.HTTPErrorType) error {
	var httperr error

	if isAlreadyExists(err) {
		httperr = errorUtils.Conflict.New(err)
	} else if isNotFound(err) {
		httperr = errorUtils.NotFound.New(err)
	} else if isForbidden(err) {
		httperr = errorUtils.Forbidden.New(err)
	} else if isUnprocessable(err) {
		httperr = errorUtils.UnprocessableEntity.New(err)
	} else {
		httperr = defaultCode.New(err)
	}

	return httperr
}

func getChart(req *http.Request, cu chartUtils.Resolver) (*chartUtils.Details, *chart.Chart, error) {
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Could not read response body")
	}

	chartDetails, err := cu.ParseDetails(body)
	if err != nil {
		return nil, nil, err
	}

	netClient, err := cu.InitNetClient(req.Context(), chartDetails)
	if err != nil {
		return nil, nil, err
	}

	ch, err := cu.GetChart(chartDetails, netClient)
	if err != nil {
		return nil, nil, err
	}

	return chartDetails, ch, nil
}

func returnForbiddenActions(ctx context.Context, w http.ResponseWriter, forbiddenActions []auth.Action) {
	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(forbiddenActions)
	if err != nil {
		utils.SendErrResponse(ctx, w, errorCode(err))
		return
	}
	utils.SendErrResponse(ctx, w, errorUtils.Forbidden.New(errors.New(string(body))))
}

// HelmProxy client and configuration
type HelmProxy struct {
	DisableAuth bool
	ListLimit   int
	ChartClient chartUtils.Resolver
	ProxyClient proxy.TillerClient
}

func (h *HelmProxy) logStatus(ctx context.Context, namespace, name string, vo proxy.ValidationObject) {
	log := logUtils.GetLogger(ctx)
	status, err := h.ProxyClient.GetReleaseStatus(ctx, namespace, name, vo)
	if err != nil {
		log.Errorf("Unable to fetch release status of %s: %v", name, err)
	} else {
		log.Infof("Release status: %s", status)
	}
}

// CreateRelease creates a new release in the namespace given as Param
func (h *HelmProxy) CreateRelease(w http.ResponseWriter, req *http.Request, params Params) {
	vo := req.Context().Value(validationObjectKey{}).(proxy.ValidationObject)
	log := logUtils.GetLogger(req.Context())

	chartDetails, ch, err := getChart(req, h.ChartClient)
	if err != nil {
		utils.SendErrResponse(req.Context(), w, errorCode(err))
		return
	}

	if !h.DisableAuth {
		manifest, manifestErr := h.ProxyClient.ResolveManifest(req.Context(), params["namespace"], chartDetails.Values, ch, vo)
		if manifestErr != nil {
			utils.SendErrResponse(req.Context(), w, errorCode(manifestErr))
			return
		}
		userAuth := req.Context().Value(userKey{}).(auth.Checker)
		forbiddenActions, actionsErr := userAuth.GetForbiddenActions(params["namespace"], "create", manifest)
		if actionsErr != nil {
			utils.SendErrResponse(req.Context(), w, errorCode(actionsErr))
			return
		}
		if len(forbiddenActions) > 0 {
			returnForbiddenActions(req.Context(), w, forbiddenActions)
			return
		}
	}

	rel, err := h.ProxyClient.CreateRelease(req.Context(), chartDetails.ReleaseName, params["namespace"], chartDetails.Values, ch, vo)
	if err != nil {
		utils.SendErrResponse(req.Context(), w, errorCodeWithDefault(err, errorUtils.UnprocessableEntity))
		return
	}

	log.Infof("Installed release %s", rel.Name)
	h.logStatus(req.Context(), params["namespace"], rel.Name, vo)
	response.NewDataResponse(*rel).Write(w)
}

// OperateRelease decides which method to call depending in the "action" query param
func (h *HelmProxy) OperateRelease(w http.ResponseWriter, req *http.Request, params Params) {
	switch req.FormValue("action") {
	case "upgrade":
		h.UpgradeRelease(w, req, params)
	case "rollback":
		h.RollbackRelease(w, req, params)
	default:
		// By default, for maintaining compatibility, we call upgrade
		h.UpgradeRelease(w, req, params)
	}
}

// RollbackRelease performs an action over a release
func (h *HelmProxy) RollbackRelease(w http.ResponseWriter, req *http.Request, params Params) {
	log := logUtils.GetLogger(req.Context())
	log.Infof("Rolling back %s", params["releaseName"])

	vo := req.Context().Value(validationObjectKey{}).(proxy.ValidationObject)

	revision := req.FormValue("revision")
	if revision == "" {
		err := errorUtils.UnprocessableEntity.New(errors.New("Missing revision to rollback in request"))
		utils.SendErrResponse(req.Context(), w, err)
		return
	}
	revisionInt, err := strconv.ParseInt(revision, 10, 64)
	if err != nil {
		err = errors.Wrap(err, "Could not parse revision")
		utils.SendErrResponse(req.Context(), w, err)
		return
	}

	if !h.DisableAuth {
		manifest, manifestErr := h.ProxyClient.ResolveManifestFromRelease(req.Context(), params["namespace"], params["releaseName"], int32(revisionInt), vo)
		if manifestErr != nil {
			utils.SendErrResponse(req.Context(), w, errorCode(manifestErr))
			return
		}
		userAuth := req.Context().Value(userKey{}).(auth.Checker)
		// Using "upgrade" action since the concept is the same
		forbiddenActions, actionsErr := userAuth.GetForbiddenActions(params["namespace"], "upgrade", manifest)
		if actionsErr != nil {
			utils.SendErrResponse(req.Context(), w, errorCode(actionsErr))
			return
		}
		if len(forbiddenActions) > 0 {
			returnForbiddenActions(req.Context(), w, forbiddenActions)
			return
		}
	}
	rel, err := h.ProxyClient.RollbackRelease(req.Context(), params["releaseName"], params["namespace"], int32(revisionInt), vo)
	if err != nil {
		utils.SendErrResponse(req.Context(), w, errorCodeWithDefault(err, errorUtils.UnprocessableEntity))
		return
	}
	log.Infof("Rollback release for %s to %d", rel.Name, revisionInt)
	h.logStatus(req.Context(), params["namespace"], rel.Name, vo)
	response.NewDataResponse(*rel).Write(w)
}

// UpgradeRelease upgrades a release in the namespace given as Param
func (h *HelmProxy) UpgradeRelease(w http.ResponseWriter, req *http.Request, params Params) {
	log := logUtils.GetLogger(req.Context())
	log.Infof("Upgrading Helm Release")
	chartDetails, ch, err := getChart(req, h.ChartClient)
	if err != nil {
		utils.SendErrResponse(req.Context(), w, errorCode(err))
		return
	}
	vo := req.Context().Value(validationObjectKey{}).(proxy.ValidationObject)

	if !h.DisableAuth {
		manifest, manifestErr := h.ProxyClient.ResolveManifest(req.Context(), params["namespace"], chartDetails.Values, ch, vo)
		if manifestErr != nil {
			utils.SendErrResponse(req.Context(), w, errorCode(manifestErr))
			return
		}
		userAuth := req.Context().Value(userKey{}).(auth.Checker)
		forbiddenActions, actionsErr := userAuth.GetForbiddenActions(params["namespace"], "upgrade", manifest)
		if actionsErr != nil {
			utils.SendErrResponse(req.Context(), w, errorCode(actionsErr))
			return
		}
		if len(forbiddenActions) > 0 {
			returnForbiddenActions(req.Context(), w, forbiddenActions)
			return
		}
	}

	rel, err := h.ProxyClient.UpdateRelease(req.Context(), params["releaseName"], params["namespace"], chartDetails.Values, ch, vo)
	if err != nil {
		utils.SendErrResponse(req.Context(), w, errorCodeWithDefault(err, errorUtils.UnprocessableEntity))
		return
	}
	log.Infof("Upgraded release %s", rel.Name)
	h.logStatus(req.Context(), params["namespace"], rel.Name, vo)
	response.NewDataResponse(*rel).Write(w)
}

// ListAllReleases list all releases that Tiller stores
func (h *HelmProxy) ListAllReleases(w http.ResponseWriter, req *http.Request) {
	vo := req.Context().Value(validationObjectKey{}).(proxy.ValidationObject)

	apps, err := h.ProxyClient.ListReleases(req.Context(), "", h.ListLimit, req.URL.Query().Get("statuses"), vo)
	if err != nil {
		utils.SendErrResponse(req.Context(), w, errorCode(err))
		return
	}
	response.NewDataResponse(apps).Write(w)
}

// ListReleases in the namespace given as Param
func (h *HelmProxy) ListReleases(w http.ResponseWriter, req *http.Request, params Params) {
	vo := req.Context().Value(validationObjectKey{}).(proxy.ValidationObject)
	apps, err := h.ProxyClient.ListReleases(req.Context(), params["namespace"], h.ListLimit, req.URL.Query().Get("statuses"), vo)
	if err != nil {
		utils.SendErrResponse(req.Context(), w, errorCode(err))
		return
	}
	response.NewDataResponse(apps).Write(w)
}

// GetRelease returns the release info
func (h *HelmProxy) GetRelease(w http.ResponseWriter, req *http.Request, params Params) {
	vo := req.Context().Value(validationObjectKey{}).(proxy.ValidationObject)
	rel, err := h.ProxyClient.GetRelease(req.Context(), params["releaseName"], params["namespace"], vo)
	if err != nil {
		utils.SendErrResponse(req.Context(), w, errorCode(err))
		return
	}

	if !h.DisableAuth {
		manifest, err := h.ProxyClient.ResolveManifest(req.Context(), params["namespace"], "", rel.Chart, vo)
		if err != nil {
			utils.SendErrResponse(req.Context(), w, errorCode(err))
			return
		}
		userAuth := req.Context().Value(userKey{}).(auth.Checker)
		forbiddenActions, err := userAuth.GetForbiddenActions(params["namespace"], "get", manifest)
		if err != nil {
			utils.SendErrResponse(req.Context(), w, errorCode(err))
			return
		}
		if len(forbiddenActions) > 0 {
			returnForbiddenActions(req.Context(), w, forbiddenActions)
			return
		}
	}

	var myRelease = chartUtils.KubeappsRelease(*rel)

	response.NewDataResponse(&myRelease).Write(w)
}

// DeleteRelease removes a release from a namespace
func (h *HelmProxy) DeleteRelease(w http.ResponseWriter, req *http.Request, params Params) {
	vo := req.Context().Value(validationObjectKey{}).(proxy.ValidationObject)
	log := logUtils.GetLogger(req.Context())

	if !h.DisableAuth {
		rel, err := h.ProxyClient.GetRelease(req.Context(), params["releaseName"], params["namespace"], vo)
		if err != nil {
			utils.SendErrResponse(req.Context(), w, errorCode(err))
			return
		}
		manifest, err := h.ProxyClient.ResolveManifest(req.Context(), params["namespace"], "", rel.Chart, vo)
		if err != nil {
			utils.SendErrResponse(req.Context(), w, errorCode(err))
			return
		}
		userAuth := req.Context().Value(userKey{}).(auth.Checker)
		forbiddenActions, err := userAuth.GetForbiddenActions(params["namespace"], "delete", manifest)
		if err != nil {
			utils.SendErrResponse(req.Context(), w, errorCode(err))
			return
		}
		if len(forbiddenActions) > 0 {
			returnForbiddenActions(req.Context(), w, forbiddenActions)
			return
		}
	}

	keepHistory := false
	if req.URL.Query().Get("keepHistory") == "1" || req.URL.Query().Get("keepHistory") == utils.StrTrue {
		keepHistory = true
	}
	err := h.ProxyClient.DeleteRelease(req.Context(), params["releaseName"], params["namespace"], keepHistory, vo)
	if err != nil {
		utils.SendErrResponse(req.Context(), w, errorCode(err))
		return
	}
	w.Header().Set("Status-Code", "200")

	_, err = w.Write([]byte("OK"))
	if err != nil {
		err = errors.Wrap(err, "Could not write status code")
		log.Error(err)
	}
}
