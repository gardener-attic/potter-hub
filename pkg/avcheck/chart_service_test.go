package avcheck

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arschles/assert"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"

	logUtils "github.com/gardener/potter-hub/pkg/log"
)

func TestGetCharts(t *testing.T) {
	tests := []struct {
		name            string
		httpReturnCode  int
		getChartsData   map[string]interface{}
		checkSuccessful bool
	}{
		{
			name:           "successful",
			httpReturnCode: http.StatusOK,
			getChartsData: map[string]interface{}{
				"data": []interface{}{
					map[string]interface{}{
						"id": "stable/grafana",
					},
				},
			},
			checkSuccessful: true,
		},
		{
			name:           "no charts in chart data response",
			httpReturnCode: http.StatusOK,
			getChartsData: map[string]interface{}{
				"data": []interface{}{},
			},
			checkSuccessful: false,
		},
		{
			name:            "invalid status code",
			httpReturnCode:  http.StatusInternalServerError,
			getChartsData:   map[string]interface{}{},
			checkSuccessful: false,
		},
		{
			name:           "no data key in chart data response",
			httpReturnCode: http.StatusOK,
			getChartsData: map[string]interface{}{
				"invalidKey": []interface{}{
					map[string]interface{}{
						"id": "stable/grafana",
					},
				},
			},
			checkSuccessful: false,
		},
		{
			name:           "data key in chart data response is not a list",
			httpReturnCode: http.StatusOK,
			getChartsData: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "stable/grafana",
				},
			},
			checkSuccessful: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			marshaledBody, err := json.Marshal(tt.getChartsData)
			assert.NoErr(t, err)
			w.WriteHeader(tt.httpReturnCode)
			_, err = w.Write(marshaledBody)
			assert.NoErr(t, err)
		}))

		nullLogger, _ := test.NewNullLogger()
		getChartsResult := checkChartsAvailable(ts.URL, &logUtils.Logger{Entry: logrus.NewEntry(nullLogger)})

		ts.Close()

		assert.Equal(t, getChartsResult.CheckSuccessful, tt.checkSuccessful, "checkSuccessful")
	}
}

func TestChartServiceCheckResultAllChecksSuccessful(t *testing.T) {
	tests := []struct {
		name            string
		pingResult      *PingResult
		getChartsResult *ChartsAvailableResult
		isOk            bool
	}{
		{
			name: "all successful",
			pingResult: &PingResult{
				CheckSuccessful: true,
			},
			getChartsResult: &ChartsAvailableResult{
				CheckSuccessful: true,
			},
			isOk: true,
		},
		{
			name: "get charts result is nil",
			pingResult: &PingResult{
				CheckSuccessful: true,
			},
			getChartsResult: nil,
			isOk:            true,
		},
		{
			name: "get charts unsuccessful",
			pingResult: &PingResult{
				CheckSuccessful: true,
			},
			getChartsResult: &ChartsAvailableResult{
				CheckSuccessful: false,
			},
			isOk: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		result := ChartServiceCheckResult{
			Ping:            tt.pingResult,
			ChartsAvailable: tt.getChartsResult,
		}
		assert.Equal(t, result.AllChecksSuccessful(), tt.isOk, "isOk")
	}
}
