package testing

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	corev1 "k8s.io/api/core/v1"

	errors2 "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	hubv1 "github.com/gardener/potter-hub/pkg/external/hubcontroller/api/v1"
)

type TestClient struct {
}

func (t TestClient) Status() client.StatusWriter {
	return nil
}

func (t TestClient) Create(context.Context, runtime.Object, ...client.CreateOption) error {
	return nil
}

func (t TestClient) Delete(context.Context, runtime.Object, ...client.DeleteOption) error {
	return nil
}

func (t TestClient) Update(context.Context, runtime.Object, ...client.UpdateOption) error {
	return nil
}

func (t TestClient) Patch(context.Context, runtime.Object, client.Patch, ...client.PatchOption) error {
	return nil
}

func (t TestClient) DeleteAllOf(context.Context, runtime.Object, ...client.DeleteAllOfOption) error {
	return nil
}

func (t TestClient) Get(context.Context, client.ObjectKey, runtime.Object) error {
	return nil
}

func (t TestClient) List(context.Context, runtime.Object, ...client.ListOption) error {
	return nil
}

// Begin Unit Test Client that holds object in memory to be used as client mock
type UnitTestClient struct {
	TestClient
	ClusterBoms     map[string]*hubv1.ClusterBom
	HubDeployments  map[string]*hubv1.HubDeploymentConfig
	Secrets         map[string]*corev1.Secret
	ClusterBomSyncs map[string]*hubv1.ClusterBomSync
}

func NewUnitTestClient() *UnitTestClient {
	cli := new(UnitTestClient)
	cli.ClusterBoms = make(map[string]*hubv1.ClusterBom)
	cli.HubDeployments = make(map[string]*hubv1.HubDeploymentConfig)
	cli.Secrets = make(map[string]*corev1.Secret)
	cli.ClusterBomSyncs = make(map[string]*hubv1.ClusterBomSync)
	return cli
}

func NewUnitTestClientWithCB(bom *hubv1.ClusterBom) *UnitTestClient {
	cli := NewUnitTestClient()
	cli.AddClusterBom(bom)
	return cli
}

func NewUnitTestClientWithCBandHDC(bom *hubv1.ClusterBom, hdc *hubv1.HubDeploymentConfig) *UnitTestClient {
	cli := NewUnitTestClient()
	cli.AddClusterBom(bom)
	cli.AddHubDeploymentConfig(hdc)
	return cli
}

func NewUnitTestClientWithCBandHDCs(bom *hubv1.ClusterBom, hdcs ...*hubv1.HubDeploymentConfig) *UnitTestClient {
	cli := NewUnitTestClient()
	cli.AddClusterBom(bom)
	for i := range hdcs {
		cli.AddHubDeploymentConfig(hdcs[i])
	}
	return cli
}

func NewUnitTestClientWithHDC(hdc *hubv1.HubDeploymentConfig) *UnitTestClient {
	cli := NewUnitTestClient()
	cli.AddHubDeploymentConfig(hdc)
	return cli
}

func (t UnitTestClient) AddClusterBom(bom *hubv1.ClusterBom) {
	if bom != nil {
		t.ClusterBoms[bom.Name] = bom
	}
}

func (t UnitTestClient) AddHubDeploymentConfig(hdc *hubv1.HubDeploymentConfig) {
	if hdc != nil {
		t.HubDeployments[hdc.Name] = hdc
	}
}

func (t UnitTestClient) AddSecret(secret *corev1.Secret) {
	if secret != nil {
		t.Secrets[secret.Name] = secret
	}
}

func (t UnitTestClient) AddClusterBomSync(sync *hubv1.ClusterBomSync) {
	if sync != nil {
		t.ClusterBomSyncs[sync.Name] = sync
	}
}

func (t UnitTestClient) List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
	switch typedList := list.(type) {
	case *hubv1.HubDeploymentConfigList:
		for _, hd := range t.HubDeployments {
			typedList.Items = append(typedList.Items, *hd)
		}
		return nil

	case *hubv1.ClusterBomList:
		for _, cb := range t.ClusterBoms {
			typedList.Items = append(typedList.Items, *cb)
		}
		return nil

	case *corev1.SecretList:
		for _, s := range t.Secrets {
			typedList.Items = append(typedList.Items, *s)
		}
		return nil

	default:
		return errors2.NewNotFound(schema.GroupResource{}, "CLIENT")
	}
}

func (t UnitTestClient) GetUncached(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	return t.Get(ctx, key, obj)
}

func (t UnitTestClient) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	switch typedObj := obj.(type) {
	case *hubv1.ClusterBom:
		clusterBom := t.ClusterBoms[key.Name]
		if clusterBom == nil {
			return errors2.NewNotFound(schema.GroupResource{}, "CLIENT")
		}
		clusterBom.DeepCopyInto(typedObj)
		return nil

	case *hubv1.ClusterBomSync:
		clusterBomSync := t.ClusterBomSyncs[key.Name]
		if clusterBomSync == nil {
			return errors2.NewNotFound(schema.GroupResource{}, "CLIENT")
		}
		clusterBomSync.DeepCopyInto(typedObj)
		return nil

	case *hubv1.HubDeploymentConfig:
		hubDeploymentConfig := t.HubDeployments[key.Name]
		if hubDeploymentConfig == nil {
			return errors2.NewNotFound(schema.GroupResource{}, "CLIENT")
		}
		hubDeploymentConfig.DeepCopyInto(typedObj)
		return nil

	case *corev1.Secret:
		secret := t.Secrets[key.Name]
		if secret == nil {
			return errors2.NewNotFound(schema.GroupResource{}, "CLIENT")
		}
		secret.DeepCopyInto(typedObj)
		return nil

	default:
		return errors2.NewNotFound(schema.GroupResource{}, "CLIENT")
	}
}

func (t UnitTestClient) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	switch typedObj := obj.(type) {
	case *hubv1.ClusterBom:
		key := typedObj.Name

		if _, ok := t.ClusterBoms[key]; ok {
			return errors.New("FAKE CLIENT - could not create object")
		}

		typedObj.Status = hubv1.ClusterBomStatus{}

		t.ClusterBoms[key] = typedObj
		return nil

	case *hubv1.ClusterBomSync:
		key := typedObj.Name

		if _, ok := t.ClusterBomSyncs[key]; ok {
			return errors.New("FAKE CLIENT - could not create object")
		}

		typedObj.Status = hubv1.ClusterBomSyncStatus{}

		t.ClusterBomSyncs[key] = typedObj
		return nil

	case *hubv1.HubDeploymentConfig:
		key := typedObj.Name

		if _, ok := t.HubDeployments[key]; ok {
			return errors.New("FAKE CLIENT - could not create object")
		}

		typedObj.Status = hubv1.HubDeploymentConfigStatus{}

		t.HubDeployments[key] = typedObj
		return nil

	default:
		return errors.New("FAKE CLIENT - could not create object")
	}
}

func (t UnitTestClient) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	switch typedObj := obj.(type) {
	case *hubv1.ClusterBom:
		delete(t.ClusterBoms, typedObj.Name)
		return nil

	case *hubv1.ClusterBomSync:
		delete(t.ClusterBomSyncs, typedObj.Name)
		return nil

	case *hubv1.HubDeploymentConfig:
		delete(t.HubDeployments, typedObj.Name)
		return nil

	default:
		return errors.New("FAKE CLIENT - could not delete object")
	}
}

func (t UnitTestClient) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	switch typedObj := obj.(type) {
	case *hubv1.ClusterBom:
		key := typedObj.Name

		oldClusterBom, ok := t.ClusterBoms[key]
		if !ok {
			return errors.New("FAKE CLIENT - could not update object")
		}

		typedObj.Status = oldClusterBom.Status

		t.ClusterBoms[key] = typedObj
		return nil

	case *hubv1.ClusterBomSync:
		key := typedObj.Name

		oldClusterBomSync, ok := t.ClusterBomSyncs[key]
		if !ok {
			return errors.New("FAKE CLIENT - could not update object")
		}

		typedObj.Status = oldClusterBomSync.Status

		t.ClusterBomSyncs[key] = typedObj
		return nil

	case *hubv1.HubDeploymentConfig:
		key := typedObj.Name

		oldDeploymentConfig, ok := t.HubDeployments[key]
		if !ok {
			return errors.New("FAKE CLIENT - could not update object")
		}

		typedObj.Status = oldDeploymentConfig.Status

		t.HubDeployments[key] = typedObj
		return nil

	default:
		return errors.New("FAKE CLIENT - could not update object")
	}
}

func (t UnitTestClient) Status() client.StatusWriter {
	return &UnitTestStatusWriter{t}
}

// End Unit Test Client

type UnitTestStatusWriter struct {
	unitTestClient UnitTestClient
}

func (w *UnitTestStatusWriter) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	switch typedObj := obj.(type) {
	case *hubv1.ClusterBom:
		key := typedObj.Name

		oldClusterBom, ok := w.unitTestClient.ClusterBoms[key]
		if !ok {
			return errors.New("FAKE CLIENT - could not update status")
		}

		oldClusterBom.Status = typedObj.Status
		return nil

	case *hubv1.HubDeploymentConfig:
		key := typedObj.Name

		oldDeploymentConfig, ok := w.unitTestClient.HubDeployments[key]
		if !ok {
			return errors.New("FAKE CLIENT - could not update status")
		}

		oldDeploymentConfig.Status = typedObj.Status
		return nil

	default:
		return errors.New("FAKE CLIENT - could not update status")
	}
}

func (w *UnitTestStatusWriter) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return nil
}

type fakeTypeSpecificData struct {
	FakeString string `json:"fakeString"`
}

func FakeRawExtensionSample() *runtime.RawExtension {
	return FakeRawExtensionWithProperty("test-type-specific-1")
}

func FakeRawExtensionWithProperty(name string) *runtime.RawExtension {
	fake := fakeTypeSpecificData{FakeString: name}
	rawData, err := json.Marshal(fake)
	if err != nil {
		return nil
	}

	object := runtime.RawExtension{Raw: rawData}
	return &object
}

// UnitTestListErrorClient is a UnitTestClient whose List method fails.
type UnitTestListErrorClient struct {
	UnitTestClient
}

func NewUnitTestListErrorClient() *UnitTestListErrorClient {
	cli := NewUnitTestClient()
	return &UnitTestListErrorClient{*cli}
}

func (t UnitTestListErrorClient) List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
	return errors.New("dummy error")
}

// UnitTestGetErrorClient is a UnitTestClient whose Get method fails.
type UnitTestGetErrorClient struct {
	UnitTestClient
}

func NewUnitTestGetErrorClient() *UnitTestGetErrorClient {
	cli := NewUnitTestClient()
	return &UnitTestGetErrorClient{*cli}
}

func (t UnitTestGetErrorClient) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	return errors.New("dummy error")
}

// UnitTestDeleteErrorClient is a UnitTestClient whose Delete method fails.
type UnitTestDeleteErrorClient struct {
	UnitTestClient
}

func NewUnitTestDeleteErrorClient() *UnitTestDeleteErrorClient {
	cli := NewUnitTestClient()
	return &UnitTestDeleteErrorClient{*cli}
}

func (t UnitTestDeleteErrorClient) Delete(context.Context, runtime.Object, ...client.DeleteOption) error {
	return errors.New("dummy error")
}

// UnitTestStatusErrorWriter is a StatusWriter whose Update method fails.
type UnitTestStatusErrorWriter struct {
}

func (w *UnitTestStatusErrorWriter) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	return errors.New("dummy error")
}

func (w *UnitTestStatusErrorWriter) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return nil
}

// UnitTestStatusErrorClient is a UnitTestClient whose Status method return a UnitTestStatusErrorWriter.
type UnitTestStatusErrorClient struct {
	UnitTestClient
}

func NewUnitTestStatusErrorClient() *UnitTestStatusErrorClient {
	cli := NewUnitTestClient()
	return &UnitTestStatusErrorClient{*cli}
}

func (t UnitTestStatusErrorClient) Status() client.StatusWriter {
	return &UnitTestStatusErrorWriter{}
}

type HubControllerTestClient struct {
	TestClient
	ReconcileMap *corev1.ConfigMap
}

func (t HubControllerTestClient) Create(ctx context.Context, obj runtime.Object, options ...client.CreateOption) error {
	switch typedObj := obj.(type) {
	case *corev1.ConfigMap:
		if t.ReconcileMap != nil {
			return errors.New("FAKE CLIENT - reconcilemap does already exist")
		}

		t.ReconcileMap = typedObj
	default:
		return errors.New("FAKE CLIENT - unsupported type")
	}
	return nil
}

func (t HubControllerTestClient) Update(ctx context.Context, obj runtime.Object, options ...client.UpdateOption) error {
	switch typedObj := obj.(type) {
	case *corev1.ConfigMap:
		if t.ReconcileMap == nil {
			return errors.New("FAKE CLIENT - reconcilemap does not exist")
		}

		t.ReconcileMap = typedObj
	default:
		return errors.New("FAKE CLIENT - unsupported type")
	}
	return nil
}

func (t HubControllerTestClient) GetUncached(ctx context.Context, cli client.ObjectKey, obj runtime.Object) error {
	return t.Get(ctx, cli, obj)
}

func (t HubControllerTestClient) Get(ctx context.Context, cli client.ObjectKey, obj runtime.Object) error {
	switch typedObj := obj.(type) {
	case *corev1.ConfigMap:
		t.ReconcileMap.DeepCopyInto(typedObj)
	default:
		return errors.New("FAKE CLIENT - unsupported type")
	}
	return nil
}

type FakeReconcileClock struct {
	Time time.Time
}

func (c *FakeReconcileClock) Now() time.Time {
	return c.Time
}

func (c *FakeReconcileClock) Sleep(d time.Duration) {
}
