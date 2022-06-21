package kubeval

import (
	"context"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	errorUtils "github.com/gardener/potter-hub/pkg/errors"
)

func GetKubeconfigFromOidcCluster(token, namespace, secretName, oidcClusterURL string, decodedOidcClusterCA []byte) (*string, error) {
	var config *rest.Config
	var err error

	if len(decodedOidcClusterCA) == 0 {
		decodedOidcClusterCA = nil
	}

	if oidcClusterURL != "" {
		config = &rest.Config{
			Host: oidcClusterURL,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: decodedOidcClusterCA,
			},
			BearerToken: token,
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, errorUtils.InternalServerError.New(errors.Wrap(err, "could not build k8s client in cluster config"))
		}
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errorUtils.InternalServerError.New(errors.Wrap(err, "Could not build k8s client for config"))
	}

	secret, err := k8sClient.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsBadRequest(err) {
			return nil, errorUtils.BadRequest.New(errors.New(err.Error()))
		} else if apierrors.IsUnauthorized(err) {
			return nil, errorUtils.Unauthorized.New(errors.New(err.Error()))
		} else if apierrors.IsForbidden(err) {
			return nil, errorUtils.Forbidden.New(errors.New(err.Error()))
		} else if apierrors.IsNotFound(err) {
			return nil, errorUtils.NotFound.New(errors.New(err.Error()))
		} else {
			return nil, errorUtils.InternalServerError.New(errors.New(err.Error()))
		}
	}
	kubeconfig := string(secret.Data["kubeconfig"])

	return &kubeconfig, nil
}
