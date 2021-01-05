package proxy

import (
	"context"
	"encoding/base64"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/release"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	logUtils "github.com/gardener/potter-hub/pkg/log"
)

// imageSecretManager manages the creation and deletion of image pull secrets for private repositories
type imageSecretManager struct {
	k8sClient     kubernetesClientInterface
	release       *release.Release
	sapSecretName string
	sapSecret     *corev1.Secret
	secretReader  secretConfigurationReader
}

// NewImageSecrets creates a new instance of an image secret.
// Please pass an instance of *ValidationObject* for kubeconfig creation and
// a string *namespace* in which the secret shall be created/updated.
func newImageSecret(ctx context.Context, rel *release.Release, dockerconfigjson string, vo ValidationObject) *imageSecretManager {
	log := logUtils.GetLogger(ctx)

	imageSecret := new(imageSecretManager)

	reader := helmConfigurationReader{}
	imageSecret.secretReader = &reader

	clientset, err := vo.getClientSet(rel.Namespace)
	if err != nil {
		log.Fatalf("Error creating kubernetes client for namespace %s. Error: %v.", rel.Namespace, err)
		return nil
	}

	var k8sClient = &kubernetesClient{clientset: clientset}
	imageSecret.k8sClient = k8sClient

	imageSecret.release = rel

	imageSecret.sapSecretName = "hubsec"
	// The Secret creation encodes the config automatically and it is part of the helm chart also base64 encoded.
	// to prevent double encoding, we decode here, and with this write the secret as-is from the env variable (helm value)
	// into the secret.
	dockerSecretsString, err := base64.StdEncoding.DecodeString(dockerconfigjson)
	if err != nil {
		log.Println("Error decoding docker config image pull secret")
		return nil
	}
	imageSecret.sapSecret = &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      imageSecret.sapSecretName,
			Namespace: rel.Namespace,
		},
		StringData: map[string]string{
			".dockerconfigjson": string(dockerSecretsString),
		},
		Type: corev1.SecretTypeDockerConfigJson,
	}

	return imageSecret
}

// createOrUpdateImageSecret Ensures (creates/updates) the image pull secret.
func (i *imageSecretManager) createOrUpdateImageSecret(ctx context.Context) error {
	log := logUtils.GetLogger(ctx)

	secretEnabled := i.secretReader.isSecretEnabled(ctx, i.release)
	if !secretEnabled {
		log.Printf("sap image secret %s is not required for chart %s", i.sapSecretName, i.release.Name)
		return nil
	}

	log.Printf("Ensuring sap image secret %s exists in namespace %s", i.sapSecretName, i.release.Namespace)

	_, err := i.k8sClient.createSecret(i.release.Namespace, i.sapSecret)
	if k8sErrors.IsAlreadyExists(err) {
		_, err = i.k8sClient.updateSecret(i.release.Namespace, i.sapSecret)
		if err != nil {
			return errors.Wrapf(err, "Unable to update sap image secret %s in namespace %s", i.sapSecretName, i.release.Namespace)
		}
	} else if err != nil {
		return errors.Wrapf(err, "Unable to create/update sap image secret %s in namespace %s", i.sapSecretName, i.release.Namespace)
	}

	log.Printf("sap image secret %s successfully installed in %s", i.sapSecretName, i.release.Namespace)

	return nil
}

// deleteImageSecret Deletes the image pull secret.
func (i *imageSecretManager) deleteImageSecret(ctx context.Context) error {
	log := logUtils.GetLogger(ctx)

	log.Printf("Deleting sap image secret %s from namespace %s", i.sapSecretName, i.release.Namespace)

	err := i.k8sClient.deleteSecret(i.release.Namespace, i.sapSecretName)
	if err != nil && !k8sErrors.IsNotFound(err) {
		return errors.Wrapf(err, "Unable to delete sap image secret %s in namespace %s",
			i.sapSecretName, i.release.Namespace)
	}
	return nil
}

type secretConfigurationReader interface {
	isSecretEnabled(ctx context.Context, release *release.Release) bool
}

type helmConfigurationReader struct{}

func (helmConfigurationReader) isSecretEnabled(ctx context.Context, rel *release.Release) bool {
	// We have to check the secret configuration in both, the chart values and the override values
	log := logUtils.GetLogger(ctx)

	log.Printf("chart config for 'hubsec': %s", rel.Chart.Values["hubsec"])
	if sapSecret, ok := rel.Chart.Values["hubsec"]; ok {
		sapSecretConfiguration := sapSecret.(map[string]interface{})
		// Comparison to with boolean needed to parse to bool type
		return true == sapSecretConfiguration["enabled"]
	}

	log.Printf("override config for 'hubsec': %s", rel.Config["hubsec"])
	if sapSecret, ok := rel.Config["hubsec"]; ok {
		sapSecretConfiguration := sapSecret.(map[string]interface{})
		// Comparison to with boolean needed to parse to bool type
		return true == sapSecretConfiguration["enabled"]
	}
	return false
}

type kubernetesClientInterface interface {
	createSecret(namespace string, secret *corev1.Secret) (*corev1.Secret, error)
	updateSecret(namespace string, secret *corev1.Secret) (*corev1.Secret, error)
	deleteSecret(namespace string, name string) error
}

type kubernetesClient struct {
	clientset *kubernetes.Clientset
}

func (k *kubernetesClient) createSecret(namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	return k.clientset.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
}

func (k *kubernetesClient) updateSecret(namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	return k.clientset.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
}

func (k *kubernetesClient) deleteSecret(namespace, name string) error {
	return k.clientset.CoreV1().Secrets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
