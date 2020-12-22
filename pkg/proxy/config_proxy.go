package proxy

import (
	"os"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth" // Import to initialize client auth plugins.
)

func logf(format string, v ...interface{}) {
	log.Infof(format, v...)
}

func (kv KubeconfigValidation) initActionConfig(namespace string) *action.Configuration {
	actionConfig := new(action.Configuration)

	restClientGetter := NewRemoteRESTClientGetter(kv.Kubeconfig, namespace)
	kc := kube.New(restClientGetter)
	kc.Log = logf

	clientset, err := kc.Factory.KubernetesClientSet()
	if err != nil {
		// TODO return error
		log.Fatal(err)
	}

	store := getStorageType(clientset, namespace)

	actionConfig.RESTClientGetter = restClientGetter
	actionConfig.KubeClient = kc
	actionConfig.Releases = store
	actionConfig.Log = logf

	return actionConfig
}

func getStorageType(clientset *kubernetes.Clientset, namespace string) *storage.Storage {
	var store *storage.Storage
	switch os.Getenv("HELM_DRIVER") {
	case "secret", "secrets", "":
		d := driver.NewSecrets(clientset.CoreV1().Secrets(namespace))
		d.Log = logf
		store = storage.Init(d)
	case "configmap", "configmaps":
		d := driver.NewConfigMaps(clientset.CoreV1().ConfigMaps(namespace))
		d.Log = logf
		store = storage.Init(d)
	case "memory":
		d := driver.NewMemory()
		store = storage.Init(d)
	default:
		// Not sure what to do here.
		panic("Unknown driver in HELM_DRIVER: " + os.Getenv("HELM_DRIVER"))
	}
	return store
}

func (kv KubeconfigValidation) getClientSet(namespace string) (*kubernetes.Clientset, error) {
	restClientGetter := NewRemoteRESTClientGetter(kv.Kubeconfig, namespace)

	kc := kube.New(restClientGetter)
	kc.Log = logf

	return kc.Factory.KubernetesClientSet()
}

func (tv TokenValidation) initActionConfig(namespace string) *action.Configuration {
	actionConfig := new(action.Configuration)

	restClientGetter := NewKRESTClientGetter(tv.Token, namespace)

	kc := kube.New(restClientGetter)
	kc.Log = logf

	clientset, err := kc.Factory.KubernetesClientSet()
	if err != nil {
		// TODO return error
		log.Fatal(err)
	}

	store := getStorageType(clientset, namespace)

	actionConfig.RESTClientGetter = restClientGetter
	actionConfig.KubeClient = kc
	actionConfig.Releases = store
	actionConfig.Log = logf

	return actionConfig
}

func (tv TokenValidation) getClientSet(namespace string) (*kubernetes.Clientset, error) {
	restClientGetter := NewKRESTClientGetter(tv.Token, namespace)

	kc := kube.New(restClientGetter)
	kc.Log = logf

	return kc.Factory.KubernetesClientSet()
}

// ValidationObject can be used to initiate an helm action configuration or a kubernetes client set.
type ValidationObject interface {
	initActionConfig(string) *action.Configuration
	getClientSet(namespace string) (*kubernetes.Clientset, error)
}

// TokenValidation stores a token, which can be used to create access via *ValidationObject* methods
type TokenValidation struct {
	Token string
}

// KubeconfigValidation stores a kubeconfig, which can be used to create access via *ValidationObject* methods
type KubeconfigValidation struct {
	Kubeconfig []byte
}
