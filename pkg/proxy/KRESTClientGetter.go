package proxy

import (
	"os"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type KRESTClientGetter struct {
	delegate    genericclioptions.RESTClientGetter
	bearerToken string
	namespace   string
}

func NewKRESTClientGetter(bearerToken, namespace string) *KRESTClientGetter {
	return &KRESTClientGetter{
		delegate:    genericclioptions.NewConfigFlags(false),
		bearerToken: bearerToken,
		namespace:   namespace,
	}
}

// ToRESTConfig returns restconfig
func (k *KRESTClientGetter) ToRESTConfig() (*rest.Config, error) {
	return k.ToRawKubeConfigLoader().ClientConfig()
}

// ToDiscoveryClient returns discovery client
func (k *KRESTClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	return k.delegate.ToDiscoveryClient()
}

// ToRESTMapper returns a restmapper
func (k *KRESTClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	return k.delegate.ToRESTMapper()
}

// ToRawKubeConfigLoader return kubeconfig loader as-is
func (k *KRESTClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	var config *rest.Config

	remoteClusterKubeconfig := os.Getenv("KUBECONFIG")
	if remoteClusterKubeconfig != "" {
		config, _ = clientcmd.BuildConfigFromFlags("", remoteClusterKubeconfig)
	} else {
		config, _ = rest.InClusterConfig()
	}

	// Kubeconfigs usually contain KeyData to authenticate a service account. We do not want this, because we
	// want to authenticate with the Bearer Token. Therefore we remove the KeyData and set the Bearer Token.
	config.KeyData = []byte{}
	config.BearerTokenFile = ""
	config.BearerToken = k.bearerToken

	return &ClientConfigGetter{
		config:    config,
		namespace: k.namespace,
	}
}

type ClientConfigGetter struct {
	config    *rest.Config
	namespace string
}

func (c *ClientConfigGetter) ClientConfig() (*rest.Config, error) {
	return c.config, nil
}

func (c *ClientConfigGetter) RawConfig() (clientcmdapi.Config, error) {
	panic("Not implemented")
}

func (c *ClientConfigGetter) Namespace() (string, bool, error) {
	return c.namespace, false, nil
}

func (c *ClientConfigGetter) ConfigAccess() clientcmd.ConfigAccess {
	panic("Not implemented")
}
