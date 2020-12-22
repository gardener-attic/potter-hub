package proxy

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

type RemoteRESTClientGetter struct {
	kubeconfig []byte
	namespace  string
}

func NewRemoteRESTClientGetter(kubeconfig []byte, namespace string) *RemoteRESTClientGetter {
	return &RemoteRESTClientGetter{
		kubeconfig: kubeconfig,
		namespace:  namespace,
	}
}

func (k *RemoteRESTClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	var config *rest.Config

	// TODO Read in kubeconfig with kube - client and secret name
	if string(k.kubeconfig) != "" {
		config, _ = clientcmd.RESTConfigFromKubeConfig(k.kubeconfig)
	} else {
		panic("Mustn't happen")
	}

	return &ClientConfigGetter{
		config:    config,
		namespace: k.namespace,
	}
}

// ToRESTConfig returns restconfig
func (k *RemoteRESTClientGetter) ToRESTConfig() (*rest.Config, error) {
	return k.ToRawKubeConfigLoader().ClientConfig()
}

// ToDiscoveryClient returns discovery client
func (k *RemoteRESTClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	restConfig, err := k.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	cachedDiscoveryClient := fakeCachedDiscoveryClient{delegate: discoveryClient}

	return cachedDiscoveryClient, err
}

// ToRESTMapper returns a restmapper
func (k *RemoteRESTClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := k.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)
	return expander, nil
}
