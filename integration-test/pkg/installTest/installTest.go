package installTest

import (
	"integration-test/pkg/helmProxyWrapper"
	"integration-test/pkg/k8sClient"
	"integration-test/pkg/util"

	log "github.com/sirupsen/logrus"
)

func Run(hpw *helmProxyWrapper.HelmProxyWrapper) {
	runTestsForGrafana(hpw)
	runTestsForEchoServer(hpw)
	runTestsForMongoDB(hpw)
}

func runTestsForGrafana(helmProxyWrapper *helmProxyWrapper.HelmProxyWrapper) {
	const releaseName = "muh-grafana"
	const chartName = "grafana"
	const repo = "bitnami"
	const labelSelector = "app.kubernetes.io/component=" + chartName + ",app.kubernetes.io/instance=" + releaseName
	namespace := helmProxyWrapper.Namespace

	helmProxyWrapper.InstallRelease(releaseName, chartName, "7.9.11", repo)
	log.Printf("\n == Checking Chart " + chartName + " == \n\n")
	err := util.PollFunction(30, checkGrafana, helmProxyWrapper.KubeClient, namespace, labelSelector, "8.5.3-debian-10-r5", false)
	helmProxyWrapper.CheckForErrors(err, chartName, releaseName)

	helmProxyWrapper.GetRelease(releaseName)
	helmProxyWrapper.UpgradeRelease(releaseName, chartName, "7.9.10", repo)
	err = util.PollFunction(30, checkGrafana, helmProxyWrapper.KubeClient, namespace, labelSelector, "8.5.4-debian-10-r0", false)
	helmProxyWrapper.CheckForErrors(err, chartName, releaseName)

	helmProxyWrapper.DeleteRelease(releaseName)
	err = util.PollFunction(30, checkGrafana, helmProxyWrapper.KubeClient, namespace, labelSelector, "8.5.4-debian-10-r0", true)
	helmProxyWrapper.CheckForErrors(err, chartName, releaseName)
}

func runTestsForMongoDB(helmProxyWrapper *helmProxyWrapper.HelmProxyWrapper) {
	const releaseName = "muh-mongo"
	const chartName = "mongodb"
	const repo = "bitnami"
	const labelSelector = "app.kubernetes.io/component=" + chartName + ",app.kubernetes.io/instance=" + releaseName
	namespace := helmProxyWrapper.Namespace

	chartVersion := "12.1.11"
	imageVersion := "5.0.8-debian-10-r24"
	helmProxyWrapper.InstallRelease(releaseName, chartName, chartVersion, repo)
	log.Printf("\n == Checking Chart mongodb %s, image %s == \n\n", chartVersion, imageVersion)
	err := util.PollFunction(30, checkMongo, helmProxyWrapper.KubeClient, namespace, labelSelector, imageVersion, false)
	helmProxyWrapper.CheckForErrors(err, chartName, releaseName)

	helmProxyWrapper.GetRelease(releaseName)
	helmProxyWrapper.UpgradeRelease(releaseName, chartName, "12.1.15", repo)
	err = util.PollFunction(30, checkMongo, helmProxyWrapper.KubeClient, namespace, labelSelector, "5.0.9-debian-10-r0", false)
	helmProxyWrapper.CheckForErrors(err, chartName, releaseName)

	helmProxyWrapper.DeleteRelease(releaseName)
	err = util.PollFunction(30, checkMongo, helmProxyWrapper.KubeClient, namespace, labelSelector, "5.0.9-debian-10-r0", true)
	helmProxyWrapper.CheckForErrors(err, chartName, releaseName)
}

func runTestsForEchoServer(helmProxyWrapper *helmProxyWrapper.HelmProxyWrapper) {
	const releaseName = "muh-echo-private-image"
	const chartName = "echo-server-private-image"
	const labelSelector = "app=" + chartName + ",release=" + releaseName
	namespace := helmProxyWrapper.Namespace

	log.Printf("\n == Checking secret does not exists before chart installation == \n\n")
	err := helmProxyWrapper.KubeClient.NoSecretExists(namespace, "hubsec")
	helmProxyWrapper.CheckForErrors(err, chartName, releaseName)

	helmProxyWrapper.InstallRelease(releaseName, "echo-server-private-image", "1.0.4", "sap-incubator")
	log.Printf("\n == Checking Chart echo-server-private-image == \n\n")
	err = util.PollFunction(30, checkEchoServerPrivate, helmProxyWrapper.KubeClient, namespace, labelSelector, "1.10", false)
	helmProxyWrapper.CheckForErrors(err, chartName, releaseName)

	log.Printf("\n == Checking secret exists before after installing chart == \n\n")
	err = helmProxyWrapper.KubeClient.SecretExists(namespace, "hubsec")
	helmProxyWrapper.CheckForErrors(err, chartName, releaseName)

	helmProxyWrapper.GetRelease(releaseName)
	helmProxyWrapper.UpgradeRelease(releaseName, "echo-server-private-image", "1.0.5", "sap-incubator")
	err = util.PollFunction(30, checkEchoServerPrivate, helmProxyWrapper.KubeClient, namespace, labelSelector, "1.10", false)
	helmProxyWrapper.CheckForErrors(err, chartName, releaseName)

	helmProxyWrapper.DeleteRelease(releaseName)
	err = util.PollFunction(30, checkEchoServerPrivate, helmProxyWrapper.KubeClient, namespace, labelSelector, "1.10", true)
	helmProxyWrapper.CheckForErrors(err, chartName, releaseName)

	log.Printf("\n == Checking secret does not exists after chart deletion == \n\n")
	err = helmProxyWrapper.KubeClient.NoSecretExists(namespace, "hubsec")
	helmProxyWrapper.CheckForErrors(err, chartName, releaseName)
}

func checkMongo(kc *k8sClient.K8sClient, ns, ls, imageVersion string) error {

	err := kc.IsDeploymentHealthy(ns, ls)
	if err != nil {
		return err
	}

	err = kc.IsDeployImageCorrect(ns, ls, imageVersion)
	if err != nil {
		return err
	}
	err = kc.IsPVCHealthy(ns, ls)
	if err != nil {
		return err
	}

	// output formatting
	log.Println()

	return nil
}

func checkGrafana(kc *k8sClient.K8sClient, ns, ls, imageVersion string) error {
	err := kc.IsDeploymentHealthy(ns, ls)
	if err != nil {
		return err
	}

	err = kc.IsDeployImageCorrect(ns, ls, imageVersion)
	if err != nil {
		return err
	}
	return nil
}

func checkEchoServerPrivate(kc *k8sClient.K8sClient, namespace string, labelSelector string, imageVersion string) error {
	err := kc.IsDeploymentHealthy(namespace, labelSelector)
	if err != nil {
		return err
	}

	err = kc.IsDeployImageCorrect(namespace, labelSelector, imageVersion)
	if err != nil {
		return err
	}
	return nil
}
