/*
Copyright (c) 2018 Bitnami

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package chart

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/repo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	appRepov1 "github.com/gardener/potter-hub/cmd/apprepository-controller/pkg/apis/apprepository/v1alpha1"
	appRepoClientSet "github.com/gardener/potter-hub/cmd/apprepository-controller/pkg/client/clientset/versioned"
	logUtils "github.com/gardener/potter-hub/pkg/log"
	"github.com/gardener/potter-hub/pkg/util"
)

const (
	defaultTimeoutSeconds = 180
)

type repoIndex struct {
	checksum string
	index    *repo.IndexFile
}

// nolint
var repoIndexes map[string]*repoIndex

// nolint
func init() {
	repoIndexes = map[string]*repoIndex{}
}

// Details contains the information to retrieve a Chart
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

// HTTPClient Interface to perform HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// LoadChart should return a Chart struct from an IOReader
type LoadChart func(in io.Reader) (*chart.Chart, error)

// Resolver for exposed funcs
type Resolver interface {
	ParseDetails(data []byte) (*Details, error)
	GetChart(details *Details, netClient HTTPClient) (*chart.Chart, error)
	InitNetClient(ctx context.Context, details *Details) (HTTPClient, error)
}

// Client struct contains the clients required to retrieve charts info
type Client struct {
	kubeClient    kubernetes.Interface
	appRepoClient appRepoClientSet.Interface
	load          LoadChart
	userAgent     string
	appRepo       *appRepov1.AppRepository
}

// NewChart returns a new Chart
func NewClient(kubeClient kubernetes.Interface, appRepoClient appRepoClientSet.Interface, load LoadChart, userAgent string) *Client {
	return &Client{
		kubeClient:    kubeClient,
		appRepoClient: appRepoClient,
		load:          load,
		userAgent:     userAgent,
	}
}

func getReq(rawURL string) (*http.Request, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, errors.Wrap(err, "Could not parse URL")
	}

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create request object")
	}

	return req, nil
}

func readResponseBody(res *http.Response) ([]byte, error) {
	if res == nil {
		return nil, errors.New("response must not be nil")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("chart download request failed")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read response body")
	}
	return body, nil
}

func checksum(data []byte) string {
	hasher := sha256.New()
	_, _ = hasher.Write(data)
	return string(hasher.Sum(nil))
}

// Cache the result of parsing the repo index since parsing this YAML
// is an expensive operation. See https://github.com/kubeapps/kubeapps/issues/1052
func getIndexFromCache(repoURL string, data []byte) (*repo.IndexFile, string) {
	sha := checksum(data)
	if repoIndexes[repoURL] == nil || repoIndexes[repoURL].checksum != sha {
		// The repository is not in the cache or the content changed
		return nil, sha
	}
	return repoIndexes[repoURL].index, sha
}

func storeIndexInCache(repoURL string, index *repo.IndexFile, sha string) {
	repoIndexes[repoURL] = &repoIndex{sha, index}
}

func parseIndex(data []byte) (*repo.IndexFile, error) {
	index := &repo.IndexFile{}
	err := yaml.Unmarshal(data, index)
	if err != nil {
		return index, errors.Wrap(err, "Could not unmarshall helm chart repo index")
	}
	index.SortEntries()
	return index, nil
}

// fetchRepoIndex returns a Helm repository
func fetchRepoIndex(netClient HTTPClient, repoURL string) (*repo.IndexFile, error) {
	req, err := getReq(repoURL)
	if err != nil {
		return nil, err
	}

	res, err := (netClient).Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Request failed")
	}
	data, err := readResponseBody(res)
	if err != nil {
		return nil, err
	}

	index, sha := getIndexFromCache(repoURL, data)
	if index == nil {
		// index not found in the cache, parse it
		index, err = parseIndex(data)
		if err != nil {
			return nil, err
		}
		storeIndexInCache(repoURL, index, sha)
	}
	return index, nil
}

func resolveChartURL(index, chartName string) (string, error) {
	indexURL, err := url.Parse(strings.TrimSpace(index))
	if err != nil {
		return "", errors.Wrap(err, "Could not parse chart url")
	}
	chartURL, err := indexURL.Parse(strings.TrimSpace(chartName))
	if err != nil {
		return "", errors.Wrap(err, "Could not parse chart url")
	}
	return chartURL.String(), nil
}

// findChartInRepoIndex returns the URL of a chart given a Helm repository and its name and version
func findChartInRepoIndex(repoIndex *repo.IndexFile, repoURL, chartName, chartVersion string) (string, error) {
	errMsg := fmt.Sprintf("chart %q", chartName)
	if chartVersion != "" {
		errMsg = fmt.Sprintf("%s version %q", errMsg, chartVersion)
	}
	cv, err := repoIndex.Get(chartName, chartVersion)
	if err != nil {
		return "", errors.Errorf("%s not found in repository", errMsg)
	}
	if len(cv.URLs) == 0 {
		return "", errors.Errorf("%s has no downloadable URLs", errMsg)
	}
	return resolveChartURL(repoURL, cv.URLs[0])
}

// fetchChart returns the Chart content given an URL
func fetchChart(netClient HTTPClient, chartURL string, load LoadChart) (*chart.Chart, error) {
	req, err := getReq(chartURL)
	if err != nil {
		return nil, err
	}

	res, err := (netClient).Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Request failed")
	}
	data, err := readResponseBody(res)
	if err != nil {
		return nil, err
	}
	return load(bytes.NewReader(data))
}

// ParseDetails return Chart details
func (c *Client) ParseDetails(data []byte) (*Details, error) {
	details := &Details{}
	err := json.Unmarshal(data, details)
	if err != nil {
		return nil, errors.Wrap(err, "Could not unmarshall chart details")
	}

	if details.AppRepositoryResourceName == "" {
		return nil, errors.New("An AppRepositoryResourceName is required")
	}

	return details, nil
}

// clientWithDefaultHeaders implements chart.HTTPClient interface
// and includes an override of the Do method which injects our default
// headers - User-Agent and Authorization (when present)
type clientWithDefaultHeaders struct {
	client         HTTPClient
	defaultHeaders http.Header
}

// Do HTTP request
func (c *clientWithDefaultHeaders) Do(req *http.Request) (*http.Response, error) {
	for k, v := range c.defaultHeaders {
		// Only add the default header if it's not already set in the request.
		if _, ok := req.Header[k]; !ok {
			req.Header[k] = v
		}
	}
	return c.client.Do(req)
}

// InitNetClient returns an HTTP client based on the chart details loading a
// custom CA if provided (as a secret)
func (c *Client) InitNetClient(ctx context.Context, details *Details) (HTTPClient, error) {
	log := logUtils.GetLogger(ctx)

	// Require the SystemCertPool unless the env var is explicitly set.
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		if _, ok := os.LookupEnv("TILLER_PROXY_ALLOW_EMPTY_CERT_POOL"); !ok {
			return nil, errors.Wrap(err, "Could not create system cert pool object")
		}
		caCertPool = x509.NewCertPool()
	}

	namespace := util.GetPodNamespace()

	// We grab the specified app repository (for later access to the repo URL, as well as any specified
	// auth).
	appRepo, err := c.appRepoClient.KubeappsV1alpha1().AppRepositories(namespace).Get(details.AppRepositoryResourceName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to get app repository %s", details.AppRepositoryResourceName)
	}
	c.appRepo = appRepo
	auth := appRepo.Spec.Auth

	if auth.CustomCA != nil {
		caCertSecret, err := c.kubeClient.CoreV1().Secrets(namespace).Get(context.TODO(), auth.CustomCA.SecretKeyRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to read secret %s in namespace %s", auth.CustomCA.SecretKeyRef.Name, namespace)
		}

		// Append our cert to the system pool
		customData, ok := caCertSecret.Data[auth.CustomCA.SecretKeyRef.Key]
		if !ok {
			return nil, errors.Errorf("Secret %q did not contain key %q", auth.CustomCA.SecretKeyRef.Name, auth.CustomCA.SecretKeyRef.Key)
		}
		if ok := caCertPool.AppendCertsFromPEM(customData); !ok {
			return nil, errors.Errorf("Failed to append %s to RootCAs", auth.CustomCA.SecretKeyRef.Name)
		}
	}

	defaultHeaders := http.Header{"User-Agent": []string{c.userAgent}}
	if auth.Header != nil {
		secret, err := c.kubeClient.CoreV1().Secrets(namespace).Get(context.TODO(), auth.Header.SecretKeyRef.Name, metav1.GetOptions{})
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to read secret %s in namespace %s", auth.Header.SecretKeyRef.Name, namespace)
		}
		authHeader := string(secret.Data[auth.Header.SecretKeyRef.Key])

		if strings.HasPrefix(authHeader, "Basic ") {
			trimmedBasicHeader := strings.TrimPrefix(authHeader, "Basic ")
			username, password, err := util.DecodeBasicAuthCredentials(trimmedBasicHeader)
			if err != nil {
				return nil, err
			}
			if username == "_json_key" {
				log.Info("starting gcloud oauth flow")
				accessToken, err := util.GetGCloudAccessToken(password)
				if err != nil {
					return nil, err
				}
				authHeader = "Bearer " + accessToken
				log.Info("successfully performed gcloud oauth flow and set access token in authorization header")
			}
		}
		defaultHeaders.Set("Authorization", authHeader)
	}

	// Return Transport for testing purposes
	return &clientWithDefaultHeaders{
		client: &http.Client{
			Timeout: time.Second * defaultTimeoutSeconds,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					RootCAs:    caCertPool,
					MinVersion: tls.VersionTLS12,
				},
			},
		},
		defaultHeaders: defaultHeaders,
	}, nil
}

// GetChart retrieves and loads a Chart from a registry
func (c *Client) GetChart(details *Details, netClient HTTPClient) (*chart.Chart, error) {
	repoURL := c.appRepo.Spec.URL
	if repoURL == "" {
		return nil, errors.New("apprepo URL is empty")
	}
	repoURL = strings.TrimSuffix(strings.TrimSpace(repoURL), "/") + "/index.yaml"

	repoIndex, err := fetchRepoIndex(netClient, repoURL)
	if err != nil {
		return nil, err
	}

	chartURL, err := findChartInRepoIndex(repoIndex, repoURL, details.ChartName, details.Version)
	if err != nil {
		return nil, err
	}

	chartRequested, err := fetchChart(netClient, chartURL, c.load)
	if err != nil {
		return nil, err
	}
	return chartRequested, nil
}
