package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	errUtils "github.wdf.sap.corp/kubernetes/hub/pkg/errors"
	hubv1 "github.wdf.sap.corp/kubernetes/hub/pkg/external/hubcontroller/api/v1"
	"github.wdf.sap.corp/kubernetes/hub/pkg/log"
	"github.wdf.sap.corp/kubernetes/hub/pkg/util"
)

const clusternameLabel = "hub.k8s.sap.com/cluster-name"

type BomHandler struct {
	ClientFactory  K8sClientFactory
	OidcClusterURL *string
	OidcClusterCA  *[]byte
}

// nolint:gochecknoglobals // performance reasons
var scheme = runtime.NewScheme()

// nolint:gochecknoinits
func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = hubv1.AddToScheme(scheme)
}

func (bomHandler *BomHandler) ListClusterBoms(w http.ResponseWriter, r *http.Request, params Params) {
	token, err := util.GetTokenFromRequest(r)
	if err != nil {
		util.SendErrResponse(r.Context(), w, errUtils.Unauthorized.New(err))
	}

	// create new config with ca and bearer token
	config := &rest.Config{
		Host:        *bomHandler.OidcClusterURL,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: *bomHandler.OidcClusterCA,
		},
	}

	k8sClient, err := bomHandler.ClientFactory(config)
	if err != nil {
		err = errors.Wrap(err, "cannot init K8s client")
		util.CheckAndSendK8sError(r.Context(), w, err)
		return
	}

	var clusterBoMList hubv1.ClusterBomList
	err = k8sClient.List(
		context.TODO(),
		&clusterBoMList,
		client.InNamespace(params["clusterNamespace"]),
		client.MatchingLabelsSelector{
			Selector: labels.SelectorFromSet(
				labels.Set(map[string]string{
					clusternameLabel: params["accessData"],
				}),
			),
		},
	)
	if err != nil {
		err = errors.Wrap(err, "")
		util.CheckAndSendK8sError(r.Context(), w, err)
		return
	}

	response, err := json.Marshal(clusterBoMList)
	if err != nil {
		err = errors.Wrap(err, "could not marshal cluster bom response")
		util.SendErrResponse(r.Context(), w, errUtils.InternalServerError.New(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		err = errors.Wrap(err, "could not write response body")
		log.GetLogger(r.Context()).Error(err)
	}
}

type K8sClientFactory func(*rest.Config) (client.Client, error)

func K8sClientFromConfig(config *rest.Config) (client.Client, error) {
	return client.New(
		config,
		client.Options{
			Scheme: scheme,
		},
	)
}

func (bomHandler *BomHandler) GetClusterBom(w http.ResponseWriter, r *http.Request, params Params) {
	token, err := util.GetTokenFromRequest(r)
	if err != nil {
		util.SendErrResponse(r.Context(), w, errUtils.Unauthorized.New(err))
		return
	}

	// create new config with ca and bearer token
	config := &rest.Config{
		Host:        *bomHandler.OidcClusterURL,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: *bomHandler.OidcClusterCA,
		},
	}

	k8sClient, err := bomHandler.ClientFactory(config)
	if err != nil {
		err = errors.Wrap(err, "cannot init k8sClient")
		util.CheckAndSendK8sError(r.Context(), w, err)
		return
	}
	var clusterBom hubv1.ClusterBom

	var key = types.NamespacedName{
		Name:      params["clusterBomName"],
		Namespace: params["clusterNamespace"],
	}

	err = k8sClient.Get(
		r.Context(),
		key,
		&clusterBom,
	)
	if err != nil {
		err = errors.Wrapf(err, "could not get clusterbom %s for secret name %s", key, params["accessData"])
		util.CheckAndSendK8sError(r.Context(), w, err)
		return
	}

	if clusterBom.Spec.SecretRef != params["accessData"] {
		err = errors.New(fmt.Sprintf("bom %s not found for cluster %s", key, params["accessData"]))
		util.SendErrResponse(r.Context(), w, errUtils.NotFound.New(err))
		return
	}

	response, err := json.Marshal(clusterBom)
	if err != nil {
		err = errors.Wrap(err, "could not marshal cluster bom response")
		util.SendErrResponse(r.Context(), w, errUtils.InternalServerError.New(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		err = errors.Wrap(err, "could not write response body")
		log.GetLogger(r.Context()).Error(err)
	}
}

func (bomHandler *BomHandler) UpdateClusterBom(w http.ResponseWriter, r *http.Request, params Params) {
	decoder := json.NewDecoder(r.Body)
	var clusterBom hubv1.ClusterBom
	decodeErr := decoder.Decode(&clusterBom)
	if decodeErr != nil {
		util.SendErrResponse(r.Context(), w, errUtils.BadRequest.New(decodeErr))
		return
	}

	if params["clusterBomName"] != clusterBom.GetName() {
		util.SendErrResponse(r.Context(), w, errUtils.BadRequest.NewError("BoM name in url and body don't match. Please fix this and try again"))
		return
	}

	token, err := util.GetTokenFromRequest(r)
	if err != nil {
		util.SendErrResponse(r.Context(), w, errUtils.Unauthorized.New(err))
		return
	}

	// create new config with ca and bearer token
	config := &rest.Config{
		Host:        *bomHandler.OidcClusterURL,
		BearerToken: token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: *bomHandler.OidcClusterCA,
		},
	}

	k8sClient, err := bomHandler.ClientFactory(config)
	if err != nil {
		err = errors.Wrap(err, "cannot init k8sClient")
		util.CheckAndSendK8sError(r.Context(), w, err)
		return
	}

	err = k8sClient.Update(r.Context(), &clusterBom, &client.UpdateOptions{})
	if err != nil {
		util.CheckAndSendK8sError(r.Context(), w, err)
		return
	}
	w.WriteHeader(200)
}
