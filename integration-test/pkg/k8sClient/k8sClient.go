package k8sClient

import (
	"context"
	"fmt"
	"integration-test/pkg/errors"
	"io/ioutil"
	"log"
	"strings"

	v1apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	C *kubernetes.Clientset
}

type K8sClientInt interface {
	IsServiceHealthy(ns, name string) bool
	IsPVCHealthy(ns, ls string) bool
	IsDeploymentHealthy(ns, ls string) bool
	IsDeployImageCorrect(ns, ls, iVersion string) bool
	CreateNsIfAbsent(namespace string) error
}

func (k *K8sClient) CreateNsIfAbsent(namespace string) error {
	_, err := k.C.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		_, createErr := k.C.CoreV1().Namespaces().Create(
			context.TODO(),
			&v1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Name: namespace,
				},
			},
			metav1.CreateOptions{},
		)
		if createErr != nil {
			return createErr
		}
	}

	return nil
}

func (k *K8sClient) GetIngressURL(ingressName, namespace string) string {
	ing, err := k.C.ExtensionsV1beta1().Ingresses(namespace).Get(
		context.TODO(),
		ingressName,
		metav1.GetOptions{},
	)
	if err != nil {
		log.Fatal("Could not get URL")
	}
	host := ing.Spec.TLS
	return host[0].Hosts[0]
}

func NewK8sClient(kubeconfigPath string) *K8sClient {
	kubeconfig, err := ioutil.ReadFile(kubeconfigPath)
	if err != nil {
		panic(err)
	}

	conf, err := clientcmd.RESTConfigFromKubeConfig(kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset := kubernetes.NewForConfigOrDie(conf)

	k8sClient := K8sClient{C: clientset}
	return &k8sClient
}
func (k *K8sClient) SecretExists(ns, name string) error {
	_, err := k.C.CoreV1().Secrets(ns).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
	if err != nil {
		fmt.Println("Secret " + name + " in namespace " + ns + " could not be accessed. " + err.Error())
		return err
	}
	return nil
}
func (k *K8sClient) NoSecretExists(ns, name string) error {
	_, err := k.C.CoreV1().Secrets(ns).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
	if err != nil && k8sErrors.IsNotFound(err) {
		return nil
	}

	fmt.Println("Secret " + name + " in namespace " + ns + " should not exist.")
	return err
}
func (k *K8sClient) IsDeploymentHealthy(ns, ls string) error {
	deploy, err := k.getDeployments(ns, ls)
	if err != nil {
		fmt.Println("Deployment with LabelSelector" + ls + " in namespace " + ns + " could not be accessed. " + err.Error())
		return err
	}

	for _, d := range deploy {
		fmt.Println("Checking for deployment " + d.Name)
		availableRep := d.Status.AvailableReplicas
		allReps := d.Spec.Replicas

		if availableRep != *allReps {
			return errors.NoType.Newf("Only %v out of %v replicas are up.", availableRep, *allReps)
		} else {
			fmt.Println("All replicas are up.")
		}
	}
	return nil
}

func (k *K8sClient) IsServiceHealthy(ns, ls string) error {
	_, err := k.getServices(ns, ls)
	if err != nil {
		return errors.NotFound.Wrap(err, "Service with LabelSelector"+ls+" in namespace "+ns+" could not be accessed.")
	}
	return nil
}

func (k *K8sClient) IsPVCHealthy(ns, ls string) error {
	pvcs, err := k.getPVCs(ns, ls)
	if err != nil {
		return errors.NotFound.Wrap(err, "PVC with LabelSelector "+ls+" in namespace "+ns+" could not be accessed.")
	}
	for _, pvc := range pvcs {
		fmt.Printf("\nChecking phase for PVC %s.\n", pvc.Name)
		if pvc.Status.Phase != "Bound" {
			return errors.NoType.Newf("PVC %s not in Phase Bound.", ls)
		} else {
			fmt.Printf("PVC Bound\n")
		}
	}
	return nil
}

func (k *K8sClient) IsDeployImageCorrect(ns, ls, imageV string) error {
	deploy, err := k.getDeployments(ns, ls)
	if err != nil {
		return errors.NoType.Wrapf(err, "Deployment with LabelSelector "+ls+" in namespace "+ns+" could not be accessed.")
	}
	for _, deploy := range deploy {
		fmt.Printf("\nChecking Image version for deployment %s.\n", deploy.Name)
		containers := deploy.Spec.Template.Spec.Containers
		for _, c := range containers {
			imageSplit := strings.Split(c.Image, ":")
			if imageV != imageSplit[1] {
				return errors.NoType.Newf("Image version of container %s (%v) does not equal the expected version(%v).", c.Name, imageSplit[1], imageV)
			}
			fmt.Printf("Image version of container %s (%v) matches the expected version(%v).", c.Name, imageSplit[1], imageV)
			fmt.Println()
		}
	}
	return nil
}

func (k K8sClient) getServices(ns, ls string) ([]v1.Service, error) {
	services, err := k.C.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{LabelSelector: ls})
	if err != nil {
		return nil, errors.NoType.Wrap(err, "Services could not be accessed.")
	}
	if len(services.Items) == 0 {
		return nil, errors.NotFound.New("Service could not be found resources found.")
	}
	return services.Items, nil
}
func (k K8sClient) getPVCs(ns, ls string) ([]v1.PersistentVolumeClaim, error) {
	pvcs, err := k.C.CoreV1().PersistentVolumeClaims(ns).List(context.TODO(), metav1.ListOptions{LabelSelector: ls})
	if err != nil {
		return nil, errors.NoType.Wrap(err, "PVC could not be accessed.")
	}
	if len(pvcs.Items) == 0 {
		return nil, errors.NotFound.New("PVC could not be found.")
	}

	return pvcs.Items, nil
}

func (k K8sClient) getDeployments(ns, ls string) ([]v1apps.Deployment, error) {
	deployments, err := k.C.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{LabelSelector: ls})
	if err != nil {
		return nil, errors.NoType.Wrap(err, "Deployment could not be accessed.")
	}
	if len(deployments.Items) == 0 {
		return nil, errors.NotFound.New("Deployment could not be found.")
	}
	return deployments.Items, nil
}
