package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	hubv1 "github.com/gardener/potter-hub/pkg/external/hubcontroller/api/v1"
	testUtils "github.com/gardener/potter-hub/pkg/external/hubcontroller/pkg/testing"
	logUtils "github.com/gardener/potter-hub/pkg/log"
)

func createFakeClientFactory(initObjects ...runtime.Object) K8sClientFactory {
	return func(config *rest.Config) (client.Client, error) {
		fakeClient := fake.NewFakeClientWithScheme(scheme, initObjects...)
		return fakeClient, nil
	}
}

// nolint
var clusterNamespace = "bom-test"

// nolint
var kubeconfigName = "cluster.kubeconfig"

// nolint
var kubeconfigName2 = "cluster2.kubeconfig"

// nolint
var testBom1Secret = corev1.Secret{
	ObjectMeta: v1.ObjectMeta{
		Name:      kubeconfigName,
		Namespace: clusterNamespace,
	},
	Data: map[string][]byte{
		"kubeconfig": []byte("some real nice kubeconfig"),
	},
}

// nolint
var testBom2Secret = corev1.Secret{
	ObjectMeta: v1.ObjectMeta{
		Name:      kubeconfigName2,
		Namespace: clusterNamespace,
	},
	Data: map[string][]byte{
		"kubeconfig": []byte("some real nice kubeconfig"),
	},
}

// nolint
var testBom1 = hubv1.ClusterBom{
	TypeMeta: v1.TypeMeta{
		Kind:       "ClusterBom",
		APIVersion: "hub.k8s.sap.com/v1",
	},
	ObjectMeta: v1.ObjectMeta{
		Name:      "test-bom-1",
		Namespace: clusterNamespace,
		Labels: map[string]string{
			clusternameLabel: kubeconfigName,
		},
	},
	Spec: hubv1.ClusterBomSpec{
		SecretRef: kubeconfigName,
		ApplicationConfigs: []hubv1.ApplicationConfig{
			{
				ID:               "test-app-id-1",
				ConfigType:       "helm",
				TypeSpecificData: *testUtils.FakeRawExtensionWithProperty("existing-value"),
				Values:           testUtils.FakeRawExtensionWithProperty("existing-value"),
			},
		},
	},
	Status: hubv1.ClusterBomStatus{},
}

// nolint
var testBom2 = hubv1.ClusterBom{
	ObjectMeta: v1.ObjectMeta{
		Name:      "test-bom-2",
		Namespace: clusterNamespace,
		Labels: map[string]string{
			clusternameLabel: kubeconfigName,
		},
	},
	Spec: hubv1.ClusterBomSpec{
		SecretRef: kubeconfigName,
		ApplicationConfigs: []hubv1.ApplicationConfig{
			{
				ID:               "test-app-id-1",
				ConfigType:       "helm",
				TypeSpecificData: *testUtils.FakeRawExtensionWithProperty("existing-value"),
				Values:           testUtils.FakeRawExtensionWithProperty("existing-value"),
			},
		},
	},
	Status: hubv1.ClusterBomStatus{},
}

// nolint
var testBom3 = hubv1.ClusterBom{
	ObjectMeta: v1.ObjectMeta{
		Name:      "test-bom-3",
		Namespace: clusterNamespace,
		Labels: map[string]string{
			clusternameLabel: "some-other-cluster.kubeconfig",
		},
	},
	Spec: hubv1.ClusterBomSpec{
		SecretRef: "some-other-cluster.kubeconfig",
		ApplicationConfigs: []hubv1.ApplicationConfig{
			{
				ID:               "test-app-id-1",
				ConfigType:       "helm",
				TypeSpecificData: *testUtils.FakeRawExtensionWithProperty("existing-value"),
				Values:           testUtils.FakeRawExtensionWithProperty("existing-value"),
			},
		},
	},
	Status: hubv1.ClusterBomStatus{},
}

// nolint
var testBom4 = hubv1.ClusterBom{
	ObjectMeta: v1.ObjectMeta{
		Name:      "test-bom-4",
		Namespace: clusterNamespace,
		Labels: map[string]string{
			clusternameLabel: "another-one.kubeconfig",
		},
	},
	Spec: hubv1.ClusterBomSpec{
		SecretRef: "another-one.kubeconfig",
		ApplicationConfigs: []hubv1.ApplicationConfig{
			{
				ID:               "test-app-id-1",
				ConfigType:       "helm",
				TypeSpecificData: *testUtils.FakeRawExtensionWithProperty("existing-value"),
				Values:           testUtils.FakeRawExtensionWithProperty("existing-value"),
			},
		},
	},
	Status: hubv1.ClusterBomStatus{},
}

func TestListClusterBoms(t *testing.T) {
	caData := []byte("this_is_some_Ca_Stuff_5.1246!$@365")
	url := "https://some.random.url/some/cool/path"
	tst := []struct {
		name               string
		bomHandler         BomHandler
		expectedBoms       []*hubv1.ClusterBom
		expectedHTTPStatus int
	}{
		{
			name: "test with two valid boms",
			bomHandler: BomHandler{
				OidcClusterCA:  &caData,
				OidcClusterURL: &url,
				ClientFactory:  createFakeClientFactory(&testBom1Secret, &testBom1, &testBom2, &testBom3, &testBom4),
			},
			expectedBoms:       []*hubv1.ClusterBom{&testBom1, &testBom2},
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "test with no valid boms",
			bomHandler: BomHandler{
				OidcClusterCA:  &caData,
				OidcClusterURL: &url,
				ClientFactory:  createFakeClientFactory(&testBom1Secret, &testBom3, &testBom4),
			},
			expectedBoms:       []*hubv1.ClusterBom{},
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "test with one valid bom",
			bomHandler: BomHandler{
				OidcClusterCA:  &caData,
				OidcClusterURL: &url,
				ClientFactory:  createFakeClientFactory(&testBom1Secret, &testBom1, &testBom3, &testBom4),
			},
			expectedBoms:       []*hubv1.ClusterBom{&testBom1},
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "test with secret not found",
			bomHandler: BomHandler{
				OidcClusterCA:  &caData,
				OidcClusterURL: &url,
				ClientFactory:  createFakeClientFactory(&testBom1, &testBom3, &testBom4),
			},
			expectedBoms:       []*hubv1.ClusterBom{},
			expectedHTTPStatus: http.StatusOK,
		},
	}

	for _, tt := range tst {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"GET",
				fmt.Sprintf("/%s/%s/v1/boms", clusterNamespace, kubeconfigName),
				nil,
			)
			nullLogger, _ := test.NewNullLogger()
			ctx := context.WithValue(context.TODO(), logUtils.LoggerKey{}, &logUtils.Logger{Entry: logrus.NewEntry(nullLogger)})
			request = request.WithContext(ctx)
			request.Header.Add("Authorization", "Bearer token")
			params := Params{
				"clusterNamespace": clusterNamespace,
				"accessData":       kubeconfigName,
			}

			tt.bomHandler.ListClusterBoms(recorder, request, params)
			resp := recorder.Result()
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			var clusterBoMList hubv1.ClusterBomList
			err = json.Unmarshal(body, &clusterBoMList)
			assert.NoError(t, err)

			assertBoMs(t, &clusterBoMList, tt.expectedBoms...)

			assert.Equal(t, resp.StatusCode, tt.expectedHTTPStatus, "statuscode")
		})
	}
}

func assertBoMs(t *testing.T, clusterBoMList *hubv1.ClusterBomList, boms ...*hubv1.ClusterBom) {
	for i := range boms {
		var found bool
		for k := range clusterBoMList.Items {
			if reflect.DeepEqual(clusterBoMList.Items[k], *boms[i]) {
				found = true
			}
		}
		assert.True(t, found, "clusterbom %s not found in response", boms[i].GetName())
		found = false
	}
}

func TestGetClusterBom(t *testing.T) {
	caData := []byte("this_is_some_Ca_Stuff_5.1246!$@365")
	url := "https://some.random.url/some/cool/path"
	tst := []struct {
		name               string
		bomHandler         BomHandler
		expectedBom        *hubv1.ClusterBom
		expectedHTTPStatus int
	}{
		{
			name: "test with one valid bom",
			bomHandler: BomHandler{
				OidcClusterCA:  &caData,
				OidcClusterURL: &url,
				ClientFactory:  createFakeClientFactory(&testBom1Secret, &testBom1),
			},
			expectedBom:        &testBom1,
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "test with two valid boms",
			bomHandler: BomHandler{
				OidcClusterCA:  &caData,
				OidcClusterURL: &url,
				ClientFactory:  createFakeClientFactory(&testBom1Secret, &testBom1, &testBom2Secret, &testBom2),
			},
			expectedBom:        &testBom1,
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "test valid get with many objects",
			bomHandler: BomHandler{
				OidcClusterCA:  &caData,
				OidcClusterURL: &url,
				ClientFactory:  createFakeClientFactory(&testBom1Secret, &testBom1, &testBom2Secret, &testBom2, &testBom3, &testBom4),
			},
			expectedBom:        &testBom1,
			expectedHTTPStatus: http.StatusOK,
		},
		{
			name: "test with present secret but bom missing",
			bomHandler: BomHandler{
				OidcClusterCA:  &caData,
				OidcClusterURL: &url,
				ClientFactory:  createFakeClientFactory(&testBom1Secret, &testBom3, &testBom4),
			},
			expectedBom:        &hubv1.ClusterBom{},
			expectedHTTPStatus: http.StatusNotFound,
		},
		{
			name: "test with no boms",
			bomHandler: BomHandler{
				OidcClusterCA:  &caData,
				OidcClusterURL: &url,
				ClientFactory:  createFakeClientFactory(&testBom1Secret),
			},
			expectedBom:        &hubv1.ClusterBom{},
			expectedHTTPStatus: http.StatusNotFound,
		},
		{
			name: "test with secret missing",
			bomHandler: BomHandler{
				OidcClusterCA:  &caData,
				OidcClusterURL: &url,
				ClientFactory:  createFakeClientFactory(&testBom1, &testBom2Secret, &testBom2, &testBom3),
			},
			expectedBom:        &hubv1.ClusterBom{},
			expectedHTTPStatus: http.StatusNotFound,
		},
		{
			name: "test secret does't match bom",
			bomHandler: BomHandler{
				OidcClusterCA:  &caData,
				OidcClusterURL: &url,
				ClientFactory:  createFakeClientFactory(&testBom1Secret, &testBom1, &testBom2Secret, &testBom2, &testBom3, &testBom4),
			},
			expectedBom:        &hubv1.ClusterBom{},
			expectedHTTPStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tst {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(
				"GET",
				fmt.Sprintf("/%s/%s/v1/boms", clusterNamespace, kubeconfigName),
				nil,
			)
			nullLogger, _ := test.NewNullLogger()
			ctx := context.WithValue(context.TODO(), logUtils.LoggerKey{}, &logUtils.Logger{Entry: logrus.NewEntry(nullLogger)})
			request = request.WithContext(ctx)
			request.Header.Add("Authorization", "Bearer token")

			params := Params{
				"clusterNamespace": clusterNamespace,
				"accessData":       kubeconfigName,
				"clusterBomName":   tt.expectedBom.GetName(),
			}

			tt.bomHandler.GetClusterBom(recorder, request, params)
			resp := recorder.Result()
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			var receivedClusterBom hubv1.ClusterBom
			err = json.Unmarshal(body, &receivedClusterBom)
			assert.NoError(t, err)
			assert.Equal(t, resp.StatusCode, tt.expectedHTTPStatus, "statuscode")
			assert.Equal(t, *tt.expectedBom, receivedClusterBom, "bom compare")
		})
	}
}
