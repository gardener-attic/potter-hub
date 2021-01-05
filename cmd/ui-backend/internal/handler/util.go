package handler

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/pkg/errors"

	errorUtils "github.com/gardener/potter-hub/pkg/errors"
)

func getClusterCredentialsFromKubeconfig(kubeconfig *UnmarshalKubeconfig) (ClusterCredentials, error) {
	var clusterCredentials ClusterCredentials

	for _, users := range kubeconfig.Users {
		if token, ok := users.UserData["token"]; ok {
			if len(token) > 0 {
				clusterCredentials = &TokenCredentials{Token: token}
				break
			}
		} else if username, ok := users.UserData["username"]; ok {
			if password, ok := users.UserData["password"]; ok {
				if len(password) > 0 && len(username) > 0 {
					clusterCredentials = &BasicAuthCredentials{Username: username, Password: password}
					break
				}
			}
		}
	}

	if clusterCredentials == nil {
		return nil, errorUtils.InternalServerError.New(errors.New("No credentials found in marshaled kubeconfig"))
	}

	return clusterCredentials, nil
}

func getTokenFromWSRequest(r *http.Request) (*string, error) {
	wsHeader := r.Header.Get("Sec-WebSocket-Protocol")
	wsHeaderTrimmedPrefix := strings.TrimPrefix(wsHeader, "base64url.bearer.authorization.k8s.io.")
	token := strings.TrimSuffix(wsHeaderTrimmedPrefix, ", binary.k8s.io")

	if token == "" {
		return nil, errorUtils.BadRequest.NewError("No token found in Sec-WebSocket-Protocol header for ws request")
	}
	return &token, nil
}

func newReverseProxy(targetURL string, encodedCAData []byte) (*httputil.ReverseProxy, error) {
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(encodedCAData)

	if !ok {
		return nil, errorUtils.InternalServerError.NewError("Couldn't add CA to cert pool")
	}

	apiServerURL, err := url.Parse(targetURL)
	if err != nil {
		return nil, errorUtils.InternalServerError.New(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(apiServerURL)
	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs: pool,
		},
	}

	return proxy, nil
}
