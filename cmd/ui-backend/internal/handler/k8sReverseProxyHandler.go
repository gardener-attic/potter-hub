package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	errorUtils "github.com/gardener/potter-hub/pkg/errors"
	logUtils "github.com/gardener/potter-hub/pkg/log"
	utils "github.com/gardener/potter-hub/pkg/util"
	"github.com/gardener/potter-hub/pkg/wsproxy"
)

func (rp *K8sReverseProxy) ProxyRequestToResourceCluster(w http.ResponseWriter, r *http.Request, params Params) {
	log := logUtils.GetLogger(r.Context())

	watch := r.URL.Query().Get("watch")

	namespace := params["clusterNamespace"]
	accessData := params["accessData"]

	if namespace == "" || accessData == "" {
		err := errorUtils.BadRequest.New(errors.New("namespace and secretName are missing in path"))
		utils.SendErrResponse(r.Context(), w, err)
		return
	}

	log.Infof("Route hit with url: %s", r.URL.String())

	trim := fmt.Sprintf("/%s/%s/reverse-proxy/v1/resourcecluster", namespace, accessData)
	r.RequestURI = strings.TrimPrefix(r.RequestURI, trim)
	r.URL.Path = strings.TrimPrefix(r.URL.Path, trim)

	if watch == utils.StrTrue {
		rp.serveWSRequestToResourceCluster(r, w)
	} else {
		rp.serveRequestToResourceCluster(r, w)
	}
}

func (rp *K8sReverseProxy) serveWSRequestToResourceCluster(r *http.Request, w http.ResponseWriter) {
	token, err := getTokenFromWSRequest(r)
	if err != nil {
		utils.SendErrResponse(r.Context(), w, err)
		return
	}
	tc := TokenCredentials{*token}
	tc.addCredentialsToWSRequest(r)

	rp.proxyWSRequest(r, w, rp.OidcClusterURL, rp.DecodedOidcClusterCA)
}

func (rp *K8sReverseProxy) serveRequestToResourceCluster(r *http.Request, w http.ResponseWriter) {
	token, err := getTokenFromRequest(r)
	if err != nil {
		utils.SendErrResponse(r.Context(), w, err)
		return
	}
	tc := TokenCredentials{*token}
	tc.addCredentialsToRequest(r)

	proxyRequest(r, w, rp.OidcClusterURL, rp.DecodedOidcClusterCA)
}

func (rp *K8sReverseProxy) ProxyRequestToTargetCluster(w http.ResponseWriter, r *http.Request, params Params) {
	log := logUtils.GetLogger(r.Context())

	watch := r.URL.Query().Get("watch")

	namespace := params["clusterNamespace"]
	accessData := params["accessData"]

	if namespace == "" || accessData == "" {
		err := errorUtils.BadRequest.New(errors.New("namespace and secretName are missing in path"))
		utils.SendErrResponse(r.Context(), w, err)
		return
	}

	trim := fmt.Sprintf("/%s/%s/reverse-proxy/v1/", namespace, accessData)
	r.RequestURI = strings.TrimPrefix(r.RequestURI, trim)
	r.URL.Path = strings.TrimPrefix(r.URL.Path, trim)

	log.Infof("Route hit with url: %s", r.URL.String())

	if watch == utils.StrTrue {
		rp.serveWSRequestToTargetCluster(r, w, namespace, accessData)
	} else {
		rp.serveRequestToTargetCluster(r, w, namespace, accessData)
	}
}

func (rp *K8sReverseProxy) serveRequestToTargetCluster(r *http.Request, w http.ResponseWriter, namespace, accessData string) {
	token, err := getTokenFromRequest(r)
	if err != nil {
		utils.SendErrResponse(r.Context(), w, err)
		return
	}

	kubeconfig, err := rp.getKubeconfig(*token, namespace, accessData)
	if err != nil {
		utils.SendErrResponse(r.Context(), w, err)
		return
	}

	kubeconfig.Credentials.addCredentialsToRequest(r)

	proxyRequest(r, w, kubeconfig.APIServer, []byte(kubeconfig.CaCert))
}

func (rp *K8sReverseProxy) serveWSRequestToTargetCluster(r *http.Request, w http.ResponseWriter, namespace, accessData string) {
	token, err := getTokenFromWSRequest(r)
	if err != nil {
		utils.SendErrResponse(r.Context(), w, err)
		return
	}

	kubeconfig, err := rp.getKubeconfig(*token, namespace, accessData)
	if err != nil {
		utils.SendErrResponse(r.Context(), w, err)
		return
	}

	kubeconfig.Credentials.addCredentialsToWSRequest(r)

	rp.proxyWSRequest(r, w, kubeconfig.APIServer, []byte(kubeconfig.CaCert))
}

func proxyRequest(r *http.Request, w http.ResponseWriter, targetURL string, caData []byte) {
	log := logUtils.GetLogger(r.Context())
	proxy, err := newReverseProxy(targetURL, caData)
	if err != nil {
		utils.SendErrResponse(r.Context(), w, err)
		return
	}

	log.Infof("Proxying request to target cluster: %s", r.URL.String())
	log.Infof("Host: %s", r.Host)
	proxy.ServeHTTP(w, r)
}

func (rp *K8sReverseProxy) proxyWSRequest(r *http.Request, w http.ResponseWriter, targetClusterURL string, caData []byte) {
	log := logUtils.GetLogger(r.Context())

	r.Header.Set("Origin", rp.HostURL)
	apiAddress := strings.TrimPrefix(targetClusterURL, "https://")

	log.Infof("Target url: %s", apiAddress)
	targetURL, _ := url.Parse("wss://" + apiAddress)

	wsProxy := wsproxy.NewSecureProxy(targetURL, caData)

	log.Infof("proxying ws request to target cluster: %s", r.URL.String())
	log.Infof("Host: %s", r.Host)
	wsProxy.ServeHTTP(w, r)
}

func LivenessCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// We can implement further checks here. For now, we just want to know the webserver exposing this
		// endpoint is up and running. Therefore we return no error (= success) here.
		logUtils.StandardLogger().Debugf("Liveness probe endpoint called.")
	}
}

func getTokenFromRequest(req *http.Request) (*string, error) {
	authHeader := strings.Split(req.Header.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return nil, errorUtils.BadRequest.New(errors.New("No token found in Authorization header"))
	}
	return &authHeader[1], nil
}
