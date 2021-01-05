package proxy

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	logUtils "github.com/gardener/potter-hub/pkg/log"
)

type MockedKubernetesClient struct {
	mock.Mock

	createSuccessful bool

	createCalled bool
	updateCalled bool
	deleteCalled bool
}

func (m *MockedKubernetesClient) createSecret(namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	m.createCalled = true

	if m.createSuccessful {
		return &corev1.Secret{
			TypeMeta:   metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{},
			Data:       nil,
			StringData: nil,
			Type:       "",
		}, nil
	}

	return nil, errors.NewAlreadyExists(schema.GroupResource{
		Group:    "mock",
		Resource: "mock",
	}, "mock")
}

func (m *MockedKubernetesClient) updateSecret(namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	m.updateCalled = true
	return &corev1.Secret{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{},
		Data:       nil,
		StringData: nil,
		Type:       "",
	}, nil
}

func (m *MockedKubernetesClient) deleteSecret(namespace, name string) error {
	m.deleteCalled = true
	return nil
}

type MockedSecretConfigurationReader struct {
	mock.Mock
}

func (*MockedSecretConfigurationReader) isSecretEnabled(ctx context.Context, rel *release.Release) bool {
	return true
}

func generateTestData() (*MockedKubernetesClient, imageSecretManager) {
	secretReader := new(MockedSecretConfigurationReader)
	k8sMock := new(MockedKubernetesClient)
	expectedSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "testrelease",
			Namespace: "testnamespace",
		},
	}
	secretManager := imageSecretManager{
		k8sClient: k8sMock,
		release: &release.Release{
			Name:      "testrelease",
			Namespace: "testnamespace",
		},
		sapSecretName: "testsecret",
		sapSecret:     expectedSecret,
		secretReader:  secretReader,
	}
	return k8sMock, secretManager
}

func Test_CreateSecret_Successfully(t *testing.T) {
	k8sMock, secretManager := generateTestData()
	k8sMock.createSuccessful = true
	nullLogger, _ := test.NewNullLogger()
	ctx := context.WithValue(context.TODO(), logUtils.LoggerKey{}, &logUtils.Logger{Entry: logrus.NewEntry(nullLogger)})

	err := secretManager.createOrUpdateImageSecret(ctx)

	assert.True(t, k8sMock.createCalled, "Create should have been called")
	assert.False(t, k8sMock.updateCalled, "Update should not have been called called")
	assert.False(t, k8sMock.deleteCalled, "Delete should not have been called called")
	assert.Nil(t, err, "Secret should have been created without error")
}

func Test_UpdateSecret_Successfully(t *testing.T) {
	k8sMock, secretManager := generateTestData()
	k8sMock.createSuccessful = false
	nullLogger, _ := test.NewNullLogger()
	ctx := context.WithValue(context.TODO(), logUtils.LoggerKey{}, &logUtils.Logger{Entry: logrus.NewEntry(nullLogger)})

	err := secretManager.createOrUpdateImageSecret(ctx)

	assert.True(t, k8sMock.createCalled, "Create should have been called")
	assert.True(t, k8sMock.updateCalled, "Update should have been called")
	assert.False(t, k8sMock.deleteCalled, "Delete should not have been called called")
	assert.Nil(t, err, "Secret should have been updated without error")
}

func Test_DeleteSecret_Successfully(t *testing.T) {
	k8sMock, secretManager := generateTestData()
	nullLogger, _ := test.NewNullLogger()
	ctx := context.WithValue(context.TODO(), logUtils.LoggerKey{}, &logUtils.Logger{Entry: logrus.NewEntry(nullLogger)})

	err := secretManager.deleteImageSecret(ctx)

	assert.False(t, k8sMock.createCalled, "Create should not have been called")
	assert.False(t, k8sMock.updateCalled, "Update should not have been called")
	assert.True(t, k8sMock.deleteCalled, "Delete should have been called called")
	assert.Nil(t, err, "Secret should have been deleted without error")
}
