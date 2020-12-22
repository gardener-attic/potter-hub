package avcheck

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/mock"

	logUtils "github.wdf.sap.corp/kubernetes/hub/pkg/log"
)

type chartServiceCheckerMock struct {
	mock.Mock
}

func (c *chartServiceCheckerMock) RunChecks(ctx context.Context, completeCheck bool) *ChartServiceCheckResult {
	args := c.Called(ctx, completeCheck)
	return args.Get(0).(*ChartServiceCheckResult)
}

func (c *chartServiceCheckerMock) StartChartsAvailableCheckBackgroundJob() {
	c.Called()
}

type dashboardCheckerMock struct {
	mock.Mock
}

func (c *dashboardCheckerMock) RunChecks(ctx context.Context) *DashboardCheckResult {
	args := c.Called(ctx)
	return args.Get(0).(*DashboardCheckResult)
}

type uiBackendCheckerMock struct {
	mock.Mock
}

func (c *uiBackendCheckerMock) RunChecks(ctx context.Context) *UIBackendCheckResult {
	args := c.Called(ctx)
	return args.Get(0).(*UIBackendCheckResult)
}

func TestHandleAVCheck(t *testing.T) {
	tests := []struct {
		name                    string
		chartServiceCheckResult *ChartServiceCheckResult
		dashboardCheckResult    *DashboardCheckResult
		uiBackendCheckResult    *UIBackendCheckResult
		completeCheck           bool
		expectedStatus          int
	}{
		{
			name: "basic check with all checks successful",
			chartServiceCheckResult: &ChartServiceCheckResult{
				Ping: &PingResult{
					CheckSuccessful: true,
				},
			},
			dashboardCheckResult: &DashboardCheckResult{
				Ping: &PingResult{
					CheckSuccessful: true,
				},
			},
			uiBackendCheckResult: &UIBackendCheckResult{
				Ping: &PingResult{
					CheckSuccessful: true,
				},
			},
			completeCheck:  false,
			expectedStatus: http.StatusOK,
		},
		{
			name: "complete check with failed GetCharts check",
			chartServiceCheckResult: &ChartServiceCheckResult{
				Ping: &PingResult{
					CheckSuccessful: true,
				},
				ChartsAvailable: &ChartsAvailableResult{
					CheckSuccessful: false,
				},
			},
			dashboardCheckResult: &DashboardCheckResult{
				Ping: &PingResult{
					CheckSuccessful: true,
				},
			},
			uiBackendCheckResult: &UIBackendCheckResult{
				Ping: &PingResult{
					CheckSuccessful: true,
				},
			},
			completeCheck:  true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt

		nullLogger, _ := test.NewNullLogger()
		ctx := context.WithValue(context.TODO(), logUtils.LoggerKey{}, &logUtils.Logger{Entry: logrus.NewEntry(nullLogger)})

		uiBackendChecker := &uiBackendCheckerMock{}
		uiBackendChecker.On("RunChecks", ctx).Return(tt.uiBackendCheckResult)

		chartServiceChecker := &chartServiceCheckerMock{}
		chartServiceChecker.On("RunChecks", ctx, tt.completeCheck).Return(tt.chartServiceCheckResult)

		dashboardChecker := &dashboardCheckerMock{}
		dashboardChecker.On("RunChecks", ctx).Return(tt.dashboardCheckResult)

		url := "/availability"
		if tt.completeCheck {
			url += "?completeCheck=true"
		}
		req := httptest.NewRequest("GET", url, nil)
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		handler := CreateAVCheckHandler(uiBackendChecker, chartServiceChecker, dashboardChecker)
		handler(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		body := map[string]interface{}{}
		err := json.NewDecoder(resp.Body).Decode(&body)
		assert.NoErr(t, err)

		assert.Equal(t, resp.StatusCode, tt.expectedStatus, "http status")
		uiBackendChecker.AssertExpectations(t)
		chartServiceChecker.AssertExpectations(t)
		dashboardChecker.AssertExpectations(t)
	}
}
