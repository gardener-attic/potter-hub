package handler

import (
	"reflect"
	"testing"
)

// nolint
var testBearerTokenKubeconfig, testBasicAuthKubeconfig, testBasicAndTokenKubeconfig UnmarshalKubeconfig

// nolint
func init() {
	clusterName := "test-01"
	kubeconfigWithoutUsers := UnmarshalKubeconfig{
		Clusters: []Cluster{
			{Name: clusterName,
				Cluster: ClusterData{
					CaCert: "someCert",
					Server: "https://bestServerEuWest.com/foo/bar"},
			}},
		Users: nil,
	}
	bearerTokenUser := User{
		Name: clusterName + "-token",
		UserData: map[string]string{
			"token": "someToken",
		},
	}
	basicUser := User{
		Name: clusterName + "-token",
		UserData: map[string]string{
			"username": "foo",
			"password": "bar",
		},
	}

	testBearerTokenKubeconfig = kubeconfigWithoutUsers
	testBearerTokenKubeconfig.Users = []User{
		bearerTokenUser,
	}
	testBasicAuthKubeconfig = kubeconfigWithoutUsers
	testBasicAuthKubeconfig.Users = []User{
		basicUser,
	}
	testBasicAndTokenKubeconfig = kubeconfigWithoutUsers
	testBasicAndTokenKubeconfig.Users = []User{
		bearerTokenUser,
		basicUser,
	}
}

func TestBasicAuthKubeconfig(t *testing.T) {
	test := &utilTestScenario{
		Description: "Test with basic auth kubeconfig",
		Kubeconfig:  testBasicAuthKubeconfig,
		Expected:    true,
		CC:          &BasicAuthCredentials{Username: "foo", Password: "bar"},
	}
	executeUtilTest(test, t)
}

func TestNegativBasicAuthKubeconfig(t *testing.T) {
	test := &utilTestScenario{
		Description: "Test with basic auth kubeconfig",
		Kubeconfig:  testBasicAuthKubeconfig,
		Expected:    false,
		CC:          &BasicAuthCredentials{Username: "bar", Password: "foo"},
	}
	executeUtilTest(test, t)
}

func TestBearerTokenKubeconfig(t *testing.T) {
	test := &utilTestScenario{
		Description: "Test with basic auth kubeconfig",
		Kubeconfig:  testBearerTokenKubeconfig,
		Expected:    true,
		CC:          &TokenCredentials{Token: "someToken"},
	}
	executeUtilTest(test, t)
}

func TestNegativBearerTokenKubeconfig(t *testing.T) {
	test := &utilTestScenario{
		Description: "Test with basic auth kubeconfig",
		Kubeconfig:  testBearerTokenKubeconfig,
		Expected:    false,
		CC:          &TokenCredentials{Token: "definitely the right token"},
	}
	executeUtilTest(test, t)
}

func TestBasicAndTokenKubeconfig(t *testing.T) {
	test := &utilTestScenario{
		Description: "Test with basic auth kubeconfig",
		Kubeconfig:  testBasicAndTokenKubeconfig,
		Expected:    true,
		CC:          &TokenCredentials{Token: "someToken"},
	}
	executeUtilTest(test, t)
}

func TestNegativBasicAndTokenKubeconfig(t *testing.T) {
	test := &utilTestScenario{
		Description: "Test with basic auth kubeconfig",
		Kubeconfig:  testBasicAndTokenKubeconfig,
		Expected:    false,
		CC:          &BasicAuthCredentials{Username: "foo", Password: "bar"},
	}
	executeUtilTest(test, t)
}

func executeUtilTest(test *utilTestScenario, t *testing.T) {
	testClusterCredentials, err := getClusterCredentialsFromKubeconfig(&test.Kubeconfig)
	if err != nil {
		t.Errorf("Unexpected error. %s", err.Error())
	}

	if reflect.DeepEqual(testClusterCredentials, test.CC) != test.Expected {
		t.Errorf("Unexpected cluster credentials. Expected %v, got %v", test.CC, testClusterCredentials)
	}
}

type utilTestScenario struct {
	// Scenario params
	Description string
	Kubeconfig  UnmarshalKubeconfig
	Expected    bool
	CC          ClusterCredentials
}
