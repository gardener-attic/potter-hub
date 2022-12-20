package helmProxyWrapper

import (
	"fmt"
	"integration-test/pkg/k8sClient"
	"integration-test/pkg/util"
	"log"
	"net/http"
	"time"
)

type Details struct {
	// AppRepositoryResourceName specifies an app repository resource to use
	// for the request.
	AppRepositoryResourceName string `json:"appRepositoryResourceName,omitempty"`
	// ChartName is the name of the chart within the repo.
	ChartName string `json:"chartName"`
	// ReleaseName is the Name of the release given to Tiller.
	ReleaseName string `json:"releaseName"`
	// Version is the chart version.
	Version string `json:"version"`
	// Values is a string containing (unparsed) YAML values.
	Values string `json:"values,omitempty"`
}

type HelmProxyWrapper struct {
	EncodedKubeconfig []byte
	KubeClient        *k8sClient.K8sClient
	HelmURL           string
	HTTPSBaseURL      string
	BaseURL           string
	TestNS            string
	Token             string
	Namespace         string
}

func (helmProxyWrapper *HelmProxyWrapper) InstallRequestFactory(payload interface{}) *http.Request {
	request := util.RequestBuilder("POST", helmProxyWrapper.HelmURL, payload)
	util.SetIntegrationTestHeader(request, helmProxyWrapper.EncodedKubeconfig)
	return request
}

func (helmProxyWrapper *HelmProxyWrapper) getRequestFactory(releaseName string) *http.Request {
	request := util.RequestBuilder("GET", helmProxyWrapper.HelmURL+"/"+releaseName, nil)
	util.SetIntegrationTestHeader(request, helmProxyWrapper.EncodedKubeconfig)
	return request
}

func (helmProxyWrapper *HelmProxyWrapper) upgradeRequestFactory(releaseName string, payload interface{}) *http.Request {
	request := util.RequestBuilder("PUT", helmProxyWrapper.HelmURL+"/"+releaseName, payload)
	util.SetIntegrationTestHeader(request, helmProxyWrapper.EncodedKubeconfig)
	return request
}

func (helmProxyWrapper *HelmProxyWrapper) deleteRequestFactory(releaseName string) *http.Request {
	request := util.RequestBuilder("DELETE", helmProxyWrapper.HelmURL+"/"+releaseName, nil)
	util.SetIntegrationTestHeader(request, helmProxyWrapper.EncodedKubeconfig)
	return request
}

func (helmProxyWrapper *HelmProxyWrapper) GetRelease(release string) {
	fmt.Println("Looking for installed chart " + release)

	request := helmProxyWrapper.getRequestFactory(release)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Error on sending request: ", err)
	}
	defer util.CloseBody(resp)
	responseMap := util.ConvertResponseToMap(resp, err)

	actualReleaseName := responseMap["data"]["name"].(string)
	if actualReleaseName == release {
		fmt.Println("Got the right chart: " + actualReleaseName)
	} else {
		log.Fatal("Got the wrong chart!" + actualReleaseName)
	}
}

func (helmProxyWrapper *HelmProxyWrapper) handleConflict(releaseDetails Details) {
	for i := 0; i < 15; i++ {
		helmProxyWrapper.DeleteRelease(releaseDetails.ReleaseName)
		time.Sleep(20 * time.Second)

		installationRequest := helmProxyWrapper.InstallRequestFactory(releaseDetails)
		client := &http.Client{}
		resp, err := client.Do(installationRequest)
		if err != nil {
			message := "Error deleting " + releaseDetails.ReleaseName + ": " + err.Error()
			panic(message)
		}

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("Installation of chart %v successful\n", releaseDetails.ReleaseName)
			return
		}
		util.CloseBody(resp)
		fmt.Printf("Installtion of chart %v not successful.\n", releaseDetails.ReleaseName)

		fmt.Printf("Sleeping since %v seconds\n", 20*i)
	}

	panic("Conflict could not be resolved")
}

func (helmProxyWrapper *HelmProxyWrapper) CheckForErrors(err error, chartName, releaseName string) {
	if err != nil {
		fmt.Printf("\nChart %s encountered an error. Cleaning up now...", chartName)
		helmProxyWrapper.DeleteRelease(releaseName)
		log.Fatal(err)
	}
}

func (helmProxyWrapper *HelmProxyWrapper) DeleteRelease(releaseName string) {
	log.Println("Trying to delete release " + releaseName)

	request := helmProxyWrapper.deleteRequestFactory(releaseName)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		message := "Error deleting " + releaseName + ": " + err.Error()
		panic(message)
	}
	defer util.CloseBody(resp)
	log.Println("Deletion call returned with status: " + resp.Status)
}

func (helmProxyWrapper *HelmProxyWrapper) InstallRelease(releaseName string, chartName string, version string, apprepository string) {
	log.Println("Trying to install new chart " + chartName + ", release " + releaseName + ", version " + version + " from apprepository " + apprepository)

	releaseDetails := Details{
		ChartName:                 chartName,
		ReleaseName:               releaseName,
		Version:                   version,
		AppRepositoryResourceName: apprepository,
	}

	request := helmProxyWrapper.InstallRequestFactory(releaseDetails)

	client := &http.Client{}
	resp, err := client.Do(request)
	chartDisolayText := chartName + ":" + version + " from apprepository " + apprepository + ":" + err.Error() + ", URL=" + request.URL.RequestURI()

	if err != nil {
		message := "Error installing chart " + chartDisolayText
		panic(message)
	}
	defer util.CloseBody(resp)
	log.Println("Installation call for chart " + chartDisolayText + " returned with status: " + resp.Status)

	if resp.Status == "409 Conflict" {
		log.Println(releaseName + " has already been installed, removing it...")
		helmProxyWrapper.handleConflict(releaseDetails)
	}
}

func (helmProxyWrapper *HelmProxyWrapper) UpgradeRelease(releaseName string, chartName string, version string, apprepository string) {
	log.Println("Trying to upgrade chart " + chartName + " to version " + version)

	payload := Details{
		ChartName:                 chartName,
		ReleaseName:               releaseName,
		Version:                   version,
		AppRepositoryResourceName: apprepository,
	}

	request := helmProxyWrapper.upgradeRequestFactory(releaseName, payload)

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		message := "Error upgrading chart: " + chartName + err.Error()
		panic(message)
	}
	defer util.CloseBody(resp)
	log.Println("Upgrade call returned with status: " + resp.Status)
}
