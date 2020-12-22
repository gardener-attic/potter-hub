package testing

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	hubv1 "github.wdf.sap.corp/kubernetes/hub/pkg/external/hubcontroller/api/v1"

	appRepov1 "github.wdf.sap.corp/kubernetes/hub/cmd/apprepository-controller/pkg/apis/apprepository/v1alpha1"
)

type ReactorFuncs map[string]func() error

type ReactiveMockClient struct {
	FakeClient   client.Client
	StatusWriter *ReactiveMockStatusWriter
	ReactorFuncs
}

func NewReactiveMockClient(funcs map[string]func() error, initObjs ...runtime.Object) ReactiveMockClient {
	scheme := runtime.NewScheme()

	_ = hubv1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
	_ = appRepov1.AddToScheme(scheme)

	fakeClient := fake.NewFakeClientWithScheme(scheme, initObjs...)

	return ReactiveMockClient{
		FakeClient:   fakeClient,
		StatusWriter: &ReactiveMockStatusWriter{funcs, fakeClient.Status()},
		ReactorFuncs: funcs,
	}
}

// Get retrieves an obj for the given object key from the Kubernetes Cluster.
// obj must be a struct pointer so that obj can be updated with the response
// returned by the Server.
func (rmc *ReactiveMockClient) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	fun := rmc.ReactorFuncs[key.String()]
	if fun != nil {
		return fun()
	}

	return rmc.FakeClient.Get(ctx, key, obj)
}

// List retrieves list of objects for a given namespace and list options. On a
// successful call, Items field in the list will be populated with the
// result returned from the server.
func (rmc *ReactiveMockClient) List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
	fun := getReactorFuncForObject(list, &rmc.ReactorFuncs)
	if fun != nil {
		return fun()
	}

	return rmc.FakeClient.List(ctx, list, opts...)
}

func (rmc *ReactiveMockClient) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	fun := getReactorFuncForObject(obj, &rmc.ReactorFuncs)
	if fun != nil {
		return fun()
	}

	return rmc.FakeClient.Create(ctx, obj, opts...)
}

// Delete deletes the given obj from Kubernetes cluster.
func (rmc *ReactiveMockClient) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	fun := getReactorFuncForObject(obj, &rmc.ReactorFuncs)
	if fun != nil {
		return fun()
	}

	return rmc.FakeClient.Delete(ctx, obj, opts...)
}

// Update updates the given obj in the Kubernetes cluster. obj must be a
// struct pointer so that obj can be updated with the content returned by the Server.
func (rmc *ReactiveMockClient) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	fun := getReactorFuncForObject(obj, &rmc.ReactorFuncs)
	if fun != nil {
		return fun()
	}

	return rmc.FakeClient.Update(ctx, obj, opts...)
}

// Patch patches the given obj in the Kubernetes cluster. obj must be a
// struct pointer so that obj can be updated with the content returned by the Server.
func (rmc *ReactiveMockClient) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	fun := getReactorFuncForObject(obj, &rmc.ReactorFuncs)
	if fun != nil {
		return fun()
	}

	return rmc.FakeClient.Patch(ctx, obj, patch, opts...)
}

// DeleteAllOf deletes all objects of the given type matching the given options.
func (rmc *ReactiveMockClient) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error {
	fun := getReactorFuncForObject(obj, &rmc.ReactorFuncs)
	if fun != nil {
		return fun()
	}

	return rmc.FakeClient.DeleteAllOf(ctx, obj, opts...)
}

type ReactiveMockStatusWriter struct {
	ReactorFuncs
	FakeClient client.StatusWriter
}

func (rmc *ReactiveMockClient) Status() client.StatusWriter {
	return rmc.StatusWriter
}

func (rms *ReactiveMockStatusWriter) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	fun := getReactorFuncForObject(obj, &rms.ReactorFuncs)
	if fun != nil {
		return fun()
	}

	return rms.FakeClient.Update(ctx, obj, opts...)
}

// Patch patches the given object's subresource. obj must be a struct
// pointer so that obj can be updated with the content returned by the
// Server.
func (rms *ReactiveMockStatusWriter) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	fun := getReactorFuncForObject(obj, &rms.ReactorFuncs)
	if fun != nil {
		return fun()
	}

	return rms.FakeClient.Patch(ctx, obj, patch, opts...)
}

func getReactorFuncForObject(obj runtime.Object, reactorFuncs *ReactorFuncs) func() error {
	key, err := client.ObjectKeyFromObject(obj)
	if err != nil {
		return nil
	}
	if fun, ok := (*reactorFuncs)[key.String()]; ok {
		return fun
	}
	return nil
}
