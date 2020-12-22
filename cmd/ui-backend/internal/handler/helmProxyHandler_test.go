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

package handler

import (
	"context"
	"fmt"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	logrusTest "github.com/sirupsen/logrus/hooks/test"
	"helm.sh/helm/v3/pkg/release"

	"github.wdf.sap.corp/kubernetes/hub/pkg/auth"
	authFake "github.wdf.sap.corp/kubernetes/hub/pkg/auth/fake"
	chartFake "github.wdf.sap.corp/kubernetes/hub/pkg/chart/fake"
	errorUtils "github.wdf.sap.corp/kubernetes/hub/pkg/errors"
	logUtils "github.wdf.sap.corp/kubernetes/hub/pkg/log"
	proxy2 "github.wdf.sap.corp/kubernetes/hub/pkg/proxy"
	proxyFake "github.wdf.sap.corp/kubernetes/hub/pkg/proxy/fake"
)

func TestErrorCodeWithDefault(t *testing.T) {
	type test struct {
		err          error
		defaultCode  errorUtils.HTTPErrorType
		expectedCode errorUtils.HTTPErrorType
	}
	tests := []test{
		{fmt.Errorf("a release named foo already exists"), errorUtils.InternalServerError, errorUtils.Conflict},
		{fmt.Errorf("release foo not found"), errorUtils.InternalServerError, errorUtils.NotFound},
		{fmt.Errorf("unauthorized to get release foo"), errorUtils.InternalServerError, errorUtils.Forbidden},
		{fmt.Errorf("release \"Foo \" failed"), errorUtils.InternalServerError, errorUtils.UnprocessableEntity},
		{fmt.Errorf("this is an unexpected error"), errorUtils.InternalServerError, errorUtils.InternalServerError},
		{fmt.Errorf("this is an unexpected error"), errorUtils.UnprocessableEntity, errorUtils.UnprocessableEntity},
	}
	for _, s := range tests {
		err := errorCodeWithDefault(s.err, s.defaultCode)
		code, _ := errorUtils.GetHTTPErrorType(err)
		if code != s.expectedCode {
			t.Errorf("Expected '%v' to return code %v got %v", s.err, s.expectedCode, code)
		}
	}
}

type helmProxyTestScenario struct {
	// Scenario params
	Description      string
	ExistingReleases []release.Release
	DisableAuth      bool
	ForbiddenActions []auth.Action
	// Request params
	RequestBody  string
	RequestQuery string
	Action       string
	Params       map[string]string
	// Expected result
	StatusCode        int
	RemainingReleases []release.Release
	ResponseBody      string // Optional
}

func TestCreateWithoutAuth(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Create a simple release without auth",
		ExistingReleases: []release.Release{},
		DisableAuth:      true,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
		RequestQuery: "",
		Action:       "create",
		Params:       map[string]string{"namespace": "default"},
		// Expected result
		StatusCode: 200,
		RemainingReleases: []release.Release{
			{Name: "foobar", Namespace: "default"},
		},
		ResponseBody: "",
	}

	executeHelmProxyTest(test, t)
}

func TestCreateWithAuth(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Create a simple release with auth",
		ExistingReleases: []release.Release{},
		DisableAuth:      false,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
		RequestQuery: "",
		Action:       "create",
		Params:       map[string]string{"namespace": "default"},
		// Expected result
		StatusCode: 200,
		RemainingReleases: []release.Release{
			{Name: "foobar", Namespace: "default"},
		},
		ResponseBody: "",
	}

	executeHelmProxyTest(test, t)
}

func TestConflictingCreate(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Create a conflicting release",
		ExistingReleases: []release.Release{{Name: "foobar", Namespace: "default"}},
		DisableAuth:      false,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
		RequestQuery: "",
		Action:       "create",
		Params:       map[string]string{"namespace": "default"},
		// Expected result
		StatusCode: 409,
		RemainingReleases: []release.Release{
			{Name: "foobar", Namespace: "default"},
		},
		ResponseBody: "",
	}

	executeHelmProxyTest(test, t)
}

func TestCreateWithForbiddenActions(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Create a simple release with forbidden actions",
		ExistingReleases: []release.Release{},
		DisableAuth:      false,
		ForbiddenActions: []auth.Action{
			{APIVersion: "v1", Resource: "pods", Namespace: "default", ClusterWide: false, Verbs: []string{"create"}},
		},
		// Request params
		RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
		RequestQuery: "",
		Action:       "create",
		Params:       map[string]string{"namespace": "default"},
		// Expected result
		StatusCode:        403,
		RemainingReleases: []release.Release{},
		ResponseBody:      `{"code":403,"message":"[{\"apiGroup\":\"v1\",\"resource\":\"pods\",\"namespace\":\"default\",\"clusterWide\":false,\"verbs\":[\"create\"]}]"}`,
	}

	executeHelmProxyTest(test, t)
}

func TestSimpleUpgrade(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Upgrade a simple release",
		ExistingReleases: []release.Release{{Name: "foobar", Namespace: "default"}},
		DisableAuth:      true,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
		RequestQuery: "",
		Action:       "upgrade",
		Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
		// Expected result
		StatusCode:        200,
		RemainingReleases: []release.Release{{Name: "foobar", Namespace: "default"}},
		ResponseBody:      "",
	}
	executeHelmProxyTest(test, t)
}

func TestUpgradeOfMissingRelease(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Upgrade a missing release",
		ExistingReleases: []release.Release{},
		DisableAuth:      true,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
		RequestQuery: "",
		Action:       "upgrade",
		Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
		// Expected result
		StatusCode:        404,
		RemainingReleases: []release.Release{},
		ResponseBody:      "",
	}
	executeHelmProxyTest(test, t)
}

func TestUpgradeWithForbiddenActions(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Upgrade a simple release with forbidden actions",
		ExistingReleases: []release.Release{{Name: "foobar", Namespace: "default"}},
		DisableAuth:      false,
		ForbiddenActions: []auth.Action{
			{APIVersion: "v1", Resource: "pods", Namespace: "default", ClusterWide: false, Verbs: []string{"upgrade"}},
		},
		// Request params
		RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
		RequestQuery: "",
		Action:       "upgrade",
		Params:       map[string]string{"namespace": "default"},
		// Expected result
		StatusCode:        403,
		RemainingReleases: []release.Release{{Name: "foobar", Namespace: "default"}},
		ResponseBody:      `{"code":403,"message":"[{\"apiGroup\":\"v1\",\"resource\":\"pods\",\"namespace\":\"default\",\"clusterWide\":false,\"verbs\":[\"upgrade\"]}]"}`,
	}
	executeHelmProxyTest(test, t)
}

func TestSimpleDelete(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Delete a simple release",
		ExistingReleases: []release.Release{{Name: "foobar", Namespace: "default"}},
		DisableAuth:      true,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody:  "",
		RequestQuery: "?keepHistory=true",
		Action:       "delete",
		Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
		// Expected result
		StatusCode:        200,
		RemainingReleases: []release.Release{{Name: "foobar", Namespace: "default", Info: &release.Info{Status: release.StatusUninstalled}}},
		ResponseBody:      "",
	}

	executeHelmProxyTest(test, t)
}

func TestDeleteWithPurge(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Delete and purge a simple release",
		ExistingReleases: []release.Release{{Name: "foobar", Namespace: "default"}},
		DisableAuth:      true,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody:  "",
		RequestQuery: "",
		Action:       "delete",
		Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
		// Expected result
		StatusCode:        200,
		RemainingReleases: []release.Release{},
		ResponseBody:      "",
	}

	executeHelmProxyTest(test, t)
}

func TestDeleteMissingRelease(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Delete a missing release",
		ExistingReleases: []release.Release{},
		DisableAuth:      true,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
		RequestQuery: "",
		Action:       "delete",
		Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
		// Expected result
		StatusCode:        404,
		RemainingReleases: []release.Release{},
		ResponseBody:      "",
	}

	executeHelmProxyTest(test, t)
}

func TestDeleteWithForbiddenActions(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Delete a release with forbidden actions",
		ExistingReleases: []release.Release{{Name: "foobar", Namespace: "default", Config: nil}},
		DisableAuth:      false,
		ForbiddenActions: []auth.Action{
			{APIVersion: "v1", Resource: "pods", Namespace: "default", ClusterWide: false, Verbs: []string{"delete"}},
		},
		// Request params
		RequestBody: `{"chartName": "foo", "releaseName": "foobar",	"version": "1.0.0"}`,
		RequestQuery: "",
		Action:       "delete",
		Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
		// Expected result
		StatusCode:        403,
		RemainingReleases: []release.Release{{Name: "foobar", Namespace: "default", Config: nil}},
		ResponseBody:      `{"code":403,"message":"[{\"apiGroup\":\"v1\",\"resource\":\"pods\",\"namespace\":\"default\",\"clusterWide\":false,\"verbs\":[\"delete\"]}]"}`,
	}
	executeHelmProxyTest(test, t)
}

func TestGetRelease(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Get a simple release",
		ExistingReleases: []release.Release{{Name: "foobar", Namespace: "default"}},
		DisableAuth:      true,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody:  "",
		RequestQuery: "",
		Action:       "get",
		Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
		// Expected result
		StatusCode:        200,
		RemainingReleases: []release.Release{{Name: "foobar", Namespace: "default"}},
		ResponseBody:      `{"data":{"name":"foobar","config":{},"namespace":"default"}}`,
	}
	executeHelmProxyTest(test, t)
}

func TestGetMissingRelease(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Get a missing release",
		ExistingReleases: []release.Release{},
		DisableAuth:      true,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody:  "",
		RequestQuery: "",
		Action:       "get",
		Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
		// Expected result
		StatusCode:        404,
		RemainingReleases: []release.Release{},
		ResponseBody:      "",
	}
	executeHelmProxyTest(test, t)
}

func TestGetWithForbiddenActions(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Get a release with forbidden actions",
		ExistingReleases: []release.Release{{Name: "foobar", Namespace: "default", Config: nil}},
		DisableAuth:      false,
		ForbiddenActions: []auth.Action{
			{APIVersion: "v1", Resource: "pods", Namespace: "default", ClusterWide: false, Verbs: []string{"get"}},
		},
		// Request params
		RequestBody:  "",
		RequestQuery: "",
		Action:       "get",
		Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
		// Expected result
		StatusCode:        403,
		RemainingReleases: []release.Release{{Name: "foobar", Namespace: "default", Config: nil}},
		ResponseBody:      `{"code":403,"message":"[{\"apiGroup\":\"v1\",\"resource\":\"pods\",\"namespace\":\"default\",\"clusterWide\":false,\"verbs\":[\"get\"]}]"}`,
	}
	executeHelmProxyTest(test, t)
}

func TestListAllReleases(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description: "List all releases",
		ExistingReleases: []release.Release{
			{Name: "foobar", Namespace: "default"},
			{Name: "foo", Namespace: "not-default"},
		},
		DisableAuth:      false,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody:  "",
		RequestQuery: "",
		Action:       "listall",
		Params:       map[string]string{},
		// Expected result
		StatusCode: 200,
		RemainingReleases: []release.Release{
			{Name: "foobar", Namespace: "default"},
			{Name: "foo", Namespace: "not-default"},
		},
		ResponseBody: `{"data":[{"releaseName":"foobar","description":"","version":"","namespace":"default","status":"DEPLOYED","chart":"","chartMetadata":{}},{"releaseName":"foo","description":"","version":"","namespace":"not-default","status":"DEPLOYED","chart":"","chartMetadata":{}}]}`,
	}
	executeHelmProxyTest(test, t)
}

func TestListReleasesInNamespace(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description: "List releases in a namespace",
		ExistingReleases: []release.Release{
			{Name: "foobar", Namespace: "default"},
			{Name: "foo", Namespace: "not-default"},
		},
		DisableAuth:      false,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody:  "",
		RequestQuery: "",
		Action:       "list",
		Params:       map[string]string{"namespace": "default"},
		// Expected result
		StatusCode: 200,
		RemainingReleases: []release.Release{
			{Name: "foobar", Namespace: "default"},
			{Name: "foo", Namespace: "not-default"},
		},
		ResponseBody: `{"data":[{"releaseName":"foobar","description":"","version":"","namespace":"default","status":"DEPLOYED","chart":"","chartMetadata":{}}]}`,
	}
	executeHelmProxyTest(test, t)
}

func TestFilterReleaseOnStatus(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description: "Filter releases based on status when listing",
		ExistingReleases: []release.Release{
			{Name: "foobar", Namespace: "default", Info: &release.Info{Status: release.StatusDeployed}},
			{Name: "foo", Namespace: "default", Info: &release.Info{Status: release.StatusUninstalled}},
		},
		DisableAuth:      false,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody:  "",
		RequestQuery: "?statuses=deployed",
		Action:       "list",
		Params:       map[string]string{"namespace": "default"},
		// Expected result
		StatusCode: 200,
		RemainingReleases: []release.Release{
			{Name: "foobar", Namespace: "default", Info: &release.Info{Status: release.StatusDeployed}},
			{Name: "foo", Namespace: "default", Info: &release.Info{Status: release.StatusUninstalled}},
		},
		ResponseBody: `{"data":[{"releaseName":"foobar","description":"","version":"","namespace":"default","status":"deployed","chart":"","chartMetadata":{}}]}`,
	}

	executeHelmProxyTest(test, t)
}

func TestRollback(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description: "Rolls back a release",
		ExistingReleases: []release.Release{
			{Name: "foo", Namespace: "default", Info: &release.Info{Status: release.StatusDeployed}},
		},
		DisableAuth:      false,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody:  "",
		RequestQuery: "?revision=1",
		Action:       "rollback",
		Params:       map[string]string{"namespace": "default", "releaseName": "foo"},
		// Expected result
		StatusCode: 200,
		RemainingReleases: []release.Release{
			{Name: "foo", Namespace: "default", Info: &release.Info{Status: release.StatusDeployed}},
		},
		ResponseBody: `{"data":{"name":"foo","info":{"first_deployed":"","last_deployed":"","deleted":"","status":"deployed"},"namespace":"default"}}`,
	}
	executeHelmProxyTest(test, t)
}

func TestMissingRollback(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Rollsback a missing release",
		ExistingReleases: []release.Release{},
		DisableAuth:      true,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody:  "",
		RequestQuery: "?revision=1",
		Action:       "rollback",
		Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
		// Expected result
		StatusCode:        404,
		RemainingReleases: []release.Release{},
		ResponseBody:      "",
	}
	executeHelmProxyTest(test, t)
}

func TestRollbackWithoutRevision(t *testing.T) {
	test := &helmProxyTestScenario{
		// Scenario params
		Description:      "Rollback without a revision",
		ExistingReleases: []release.Release{},
		DisableAuth:      true,
		ForbiddenActions: []auth.Action{},
		// Request params
		RequestBody:  "",
		RequestQuery: "",
		Action:       "rollback",
		Params:       map[string]string{"namespace": "default", "releaseName": "foobar"},
		// Expected result
		StatusCode:        422,
		RemainingReleases: []release.Release{},
		ResponseBody:      "",
	}
	executeHelmProxyTest(test, t)
}

func executeHelmProxyTest(test *helmProxyTestScenario, t *testing.T) {
	// Prepare environment
	proxy := &proxyFake.Proxy{
		Releases: test.ExistingReleases,
	}
	handler := HelmProxy{
		DisableAuth: test.DisableAuth,
		ListLimit:   255,
		ChartClient: &chartFake.Chart{},
		ProxyClient: proxy,
	}
	req := httptest.NewRequest("GET", fmt.Sprintf("http://foo.bar%s", test.RequestQuery), strings.NewReader(test.RequestBody))
	ctx := context.WithValue(req.Context(), validationObjectKey{}, &proxy2.TokenValidation{Token: "desu"})
	nullLogger, _ := logrusTest.NewNullLogger()
	ctx = context.WithValue(ctx, logUtils.LoggerKey{}, &logUtils.Logger{Entry: logrus.NewEntry(nullLogger)})
	req = req.WithContext(ctx)
	if !test.DisableAuth {
		fauth := &authFake.Auth{
			ForbiddenActions: test.ForbiddenActions,
		}
		ctx := context.WithValue(req.Context(), userKey{}, fauth)
		ctx = context.WithValue(ctx, logUtils.LoggerKey{}, &logUtils.Logger{Entry: logrus.NewEntry(logrus.New())})
		req = req.WithContext(ctx)
	}
	response := httptest.NewRecorder()
	// Perform request
	t.Log(test.Description)
	switch test.Action {
	case "create":
		handler.CreateRelease(response, req, test.Params)
	case "upgrade":
		handler.UpgradeRelease(response, req, test.Params)
	case "delete":
		handler.DeleteRelease(response, req, test.Params)
	case "get":
		handler.GetRelease(response, req, test.Params)
	case "rollback":
		handler.RollbackRelease(response, req, test.Params)
	case "list":
		handler.ListReleases(response, req, test.Params)
	case "listall":
		handler.ListAllReleases(response, req)
	default:
		t.Errorf("Unexpected action %s", test.Action)
	}
	// Check result
	if response.Code != test.StatusCode {
		t.Errorf("Expecting a StatusCode %d, received %d", test.StatusCode, response.Code)
	}
	if !reflect.DeepEqual(proxy.Releases, test.RemainingReleases) {
		t.Errorf("Unexpected remaining releases. Expecting %v, found %v", test.RemainingReleases, proxy.Releases)
	}
	if test.ResponseBody != "" {
		if test.ResponseBody != response.Body.String() {
			t.Errorf("Unexpected body response. Expecting %s, found %s", test.ResponseBody, response.Body)
		}
	}
}
