package proxy

import (
	openapi_v2 "github.com/googleapis/gnostic/openapiv2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/discovery"
	restclient "k8s.io/client-go/rest"
)

// fakeCachedDiscoveryClient
// implements the interface CachedDiscoveryInterface without the cache functionality. We need this interface to interact
// with the k8s interface RESTClientGetter (see method ToDiscoveryClient in RemoteRESTClientGetter).
type fakeCachedDiscoveryClient struct {
	delegate discovery.DiscoveryInterface
}

func (fakeCachedDiscoveryClient) Fresh() bool {
	return true
}

func (fakeCachedDiscoveryClient) Invalidate() {

}

func (v fakeCachedDiscoveryClient) ServerGroups() (*metav1.APIGroupList, error) {
	return v.delegate.ServerGroups()
}

func (v fakeCachedDiscoveryClient) ServerResourcesForGroupVersion(groupVersion string) (*metav1.APIResourceList, error) {
	return v.delegate.ServerResourcesForGroupVersion(groupVersion)
}

func (v fakeCachedDiscoveryClient) ServerResources() ([]*metav1.APIResourceList, error) {
	//nolint
	return v.delegate.ServerResources()
}

func (v fakeCachedDiscoveryClient) ServerGroupsAndResources() ([]*metav1.APIGroup, []*metav1.APIResourceList, error) {
	return v.delegate.ServerGroupsAndResources()
}

func (v fakeCachedDiscoveryClient) ServerPreferredResources() ([]*metav1.APIResourceList, error) {
	return v.delegate.ServerPreferredResources()
}

func (v fakeCachedDiscoveryClient) ServerPreferredNamespacedResources() ([]*metav1.APIResourceList, error) {
	return v.delegate.ServerPreferredNamespacedResources()
}

func (v fakeCachedDiscoveryClient) ServerVersion() (*version.Info, error) {
	return v.delegate.ServerVersion()
}

func (v fakeCachedDiscoveryClient) OpenAPISchema() (*openapi_v2.Document, error) {
	return v.delegate.OpenAPISchema()
}

func (v fakeCachedDiscoveryClient) RESTClient() restclient.Interface {
	return v.delegate.RESTClient()
}
