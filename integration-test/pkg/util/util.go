package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"integration-test/pkg/errors"
	"integration-test/pkg/k8sClient"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	log "github.com/sirupsen/logrus"
)

const PROJECT_NAME = "garden-hubtest"
const SECRET_NAME = "int-test.kubeconfig"

// APIKube const for kube api reverse proxying
const APIKube = "/api/ui-backend/reverse-proxy/v1"

const DEPLOYMENT_API = "/apis/apps/v1"
const CORE_V1_API = "/api/v1"

type CheckFunction func(*k8sClient.K8sClient, string, string, string) error

func PollFunction(retries int, checkFunction CheckFunction, kubeClient *k8sClient.K8sClient, namespace, labelSelector, imageVersion string, pollForDeletion bool) error {
	//func pollFunction(retries int, checkFunction checkFunction, KubeClient *k8sClient.K8sClient, namespace, labelSelector string, pollForDeletion bool) error {
	for i := 0; i < retries; i++ {
		err := checkFunction(kubeClient, namespace, labelSelector, imageVersion)
		if pollForDeletion {
			cerr := err.(errors.CustomError)
			if cerr.ErrorType == errors.NotFound {
				fmt.Printf("Deletion successfull.\n\n")
				return nil
			} else {
				fmt.Println("Resource still exists.")
				fmt.Printf("Number of retries: %v. Polling for %v seconds.\n", i, i*20)
				time.Sleep(time.Second * 20)
			}

		} else {
			if err == nil {
				fmt.Printf("Installation successfull.\n\n")
				return nil
			} else {
				fmt.Println(err)
				fmt.Printf("Number of retries: %v. Polling for %v seconds.\n", i, i*20)
				time.Sleep(time.Second * 20)
			}
		}
	}
	return fmt.Errorf("polling exceeded the number of retries")
}

func CloseBody(resp *http.Response) {
	if resp == nil {
		return
	}

	err := resp.Body.Close()
	if err != nil {
		log.Fatal("Could not close response body: " + err.Error())
	}
}

func CloseWriter(writer io.WriteCloser) {
	if writer == nil {
		return
	}

	err := writer.Close()
	if err != nil {
		log.Fatal("Could not close writer: " + err.Error())
	}
}

func RequestBuilder(method, url string, payload interface{}) *http.Request {
	var request *http.Request

	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			panic("Could not marshall post body")
		}
		request, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			panic("Could not create post request")
		}
	} else {
		var err error
		request, err = http.NewRequest(method, url, nil)
		if err != nil {
			panic("Could not create post request")
		}
	}
	return request
}

func ConvertResponseToMap(resp *http.Response, err error) map[string]map[string]interface{} {
	body, _ := ioutil.ReadAll(resp.Body)
	var responseMap map[string]map[string]interface{}
	err = json.Unmarshal(body, &responseMap)
	if err != nil {
		log.Fatal("omg")
	}
	return responseMap
}

func SetIntegrationTestHeader(r *http.Request, encodedKubeconfig []byte) {
	r.Header.Add("disableAuth", "true")
	r.Header.Add("targetKubeconfig", string(encodedKubeconfig))
}

func EncodeKubeconfig(targetKubeconfigPath string) []byte {
	targetKubeconfig, err := ioutil.ReadFile(targetKubeconfigPath)
	if err != nil {
		panic(err)
	}

	encodedKubeconfig := base64.StdEncoding.EncodeToString(targetKubeconfig)

	return []byte(encodedKubeconfig)
}

func BuildApiProxyRequestWithBearerAuth(method, url, token string, payload interface{}) (req *http.Request) {
	req = RequestBuilder(method, url, payload)
	req.Header.Add("Authorization", "Bearer "+token)
	setNamespaceAndSecretNameHeader(req)
	return
}

func BuildApiProxyRequestWithoutAuth(method, url string, payload interface{}) *http.Request {
	req := RequestBuilder(method, url, payload)
	setNamespaceAndSecretNameHeader(req)
	return req
}

func setNamespaceAndSecretNameHeader(req *http.Request) {
	req.Header.Add("SecretName", SECRET_NAME)
	req.Header.Add("Namespace", PROJECT_NAME)
}

func IsResponseLoginForm(body string) bool {
	if strings.HasPrefix(body, "<!DOCTYPE html>") {
		return true
	} else {
		return false
	}
}
