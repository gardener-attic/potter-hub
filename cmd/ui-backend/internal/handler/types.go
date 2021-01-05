package handler

import (
	"encoding/base64"
	"net/http"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/gardener/potter-hub/pkg/kubeval"
)

type UnmarshalKubeconfig struct {
	Clusters []Cluster `yaml:"clusters"`
	Users    []User    `yaml:"users"`
}

type Kubeconfig struct {
	Name        string
	APIServer   string
	CaCert      string
	Credentials ClusterCredentials
}

type Cluster struct {
	Name    string      `yaml:"name"`
	Cluster ClusterData `yaml:"cluster"`
}
type ClusterData struct {
	CaCert string `yaml:"certificate-authority-data"`
	Server string `yaml:"server"`
}

type User struct {
	Name     string            `yaml:"name"`
	UserData map[string]string `yaml:"user"`
}

func NewK8sReverseProxy(oidcClusterURL, hostURL string, decodedClusterCA []byte) *K8sReverseProxy {
	return &K8sReverseProxy{
		OidcClusterURL:       oidcClusterURL,
		DecodedOidcClusterCA: decodedClusterCA,
		HostURL:              hostURL,
	}
}

type K8sReverseProxy struct {
	OidcClusterURL       string
	DecodedOidcClusterCA []byte
	HostURL              string
}

func (rp *K8sReverseProxy) getKubeconfig(token, namespace, accessData string) (*Kubeconfig, error) {
	marshaled, err := kubeval.GetKubeconfigFromOidcCluster(token, namespace, accessData, rp.OidcClusterURL, rp.DecodedOidcClusterCA)
	if err != nil {
		return nil, err
	}

	var unmarshaledKubeconfig UnmarshalKubeconfig
	err = yaml.Unmarshal([]byte(*marshaled), &unmarshaledKubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "Could not unmarshal kubeconfig")
	}

	decodedCaCert, err := base64.StdEncoding.DecodeString(unmarshaledKubeconfig.Clusters[0].Cluster.CaCert)
	if err != nil {
		return nil, errors.Wrap(err, "Could not decode caCert")
	}

	clusterCredentials, err := getClusterCredentialsFromKubeconfig(&unmarshaledKubeconfig)
	if err != nil {
		return nil, err
	}

	kubeconfig := Kubeconfig{
		Name:        unmarshaledKubeconfig.Clusters[0].Name,
		APIServer:   unmarshaledKubeconfig.Clusters[0].Cluster.Server,
		CaCert:      string(decodedCaCert),
		Credentials: clusterCredentials,
	}

	return &kubeconfig, nil
}

type ClusterCredentials interface {
	addCredentialsToRequest(r *http.Request)
	addCredentialsToWSRequest(r *http.Request)
}

type BasicAuthCredentials struct {
	Username, Password string
}

type TokenCredentials struct {
	Token string
}

func (tc *TokenCredentials) addCredentialsToWSRequest(r *http.Request) {
	tc.addCredentialsToRequest(r)
	r.Header.Set("Sec-WebSocket-Protocol", "binary.k8s.io")
}
func (tc *TokenCredentials) addCredentialsToRequest(r *http.Request) {
	r.Header.Set("Authorization", "Bearer "+tc.Token)
}

func (bac *BasicAuthCredentials) addCredentialsToWSRequest(r *http.Request) {
	bac.addCredentialsToRequest(r)
	r.Header.Set("Sec-WebSocket-Protocol", "binary.k8s.io")
}

func (bac *BasicAuthCredentials) addCredentialsToRequest(r *http.Request) {
	r.SetBasicAuth(bac.Username, bac.Password)
}
