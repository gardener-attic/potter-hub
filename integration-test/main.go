package main

import (
	"integration-test/pkg/apiProxyTest"
	"integration-test/pkg/helmProxyWrapper"
	"integration-test/pkg/installTest"
	"integration-test/pkg/k8sClient"
	"integration-test/pkg/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var (
	kubeconfigPath       string
	targetKubeconfigPath string
	baseURL              string
	token                string
	namespace            string
)

func init() {
	pflag.StringVar(&kubeconfigPath, "kubeconfig", "", "Kubeconfig of central dev cluster")
	pflag.StringVar(&targetKubeconfigPath, "target-kubeconfig", "", "Kubeconfig of target cluster")
	pflag.StringVar(&token, "token", "", "A valid token for the oidc cluster")
	pflag.StringVar(&namespace, "namespace", "app-test", "Namespace the integration-test will run")
}

func main() {
	helmProxy := setupTestEnvironment()

	installTest.Run(helmProxy)
	apiProxyTest.Run(helmProxy)
}

func setupTestEnvironment() *helmProxyWrapper.HelmProxyWrapper {
	pflag.Parse()

	centralKubeClient := k8sClient.NewK8sClient(kubeconfigPath)

	ingressURL := centralKubeClient.GetIngressURL("hub-k8s-potter-hub", "hub")
	log.Infof("Using ingress: %s", ingressURL)
	baseURL = ingressURL
	helmURL := "https://" + ingressURL + "/" + util.PROJECT_NAME + "/" + util.SECRET_NAME + "/api/ui-backend/helm/v1/namespaces/" + namespace + "/releases"
	httpsIngressURL := "https://" + baseURL + "/" + util.PROJECT_NAME + "/" + util.SECRET_NAME
	log.Println("Using baseURL: " + baseURL)
	log.Println("Installing in namespace " + namespace)

	targetKubeClient := k8sClient.NewK8sClient(targetKubeconfigPath)

	err := targetKubeClient.CreateNsIfAbsent(namespace)
	if err != nil {
		log.Error("Can't init namespace")
		panic(err)
	}

	encodedKubeconfig := util.EncodeKubeconfig(targetKubeconfigPath)

	return &helmProxyWrapper.HelmProxyWrapper{
		EncodedKubeconfig: encodedKubeconfig,
		KubeClient:        targetKubeClient,
		HelmURL:           helmURL,
		BaseURL:           baseURL,
		HTTPSBaseURL:      httpsIngressURL,
		Token:             token,
		Namespace:         namespace,
	}
}
