package avcheck

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	logUtils "github.com/gardener/potter-hub/pkg/log"
)

const completeCheckKey = "completeCheck"

type response struct {
	ChartService *ChartServiceCheckResult `json:"chartService"`
	Dashboard    *DashboardCheckResult    `json:"dashboard"`
	UIBackend    *UIBackendCheckResult    `json:"uiBackend"`
}

func CreateAVCheckHandler(uiBackendChecker IUIBackendChecker, chartServiceChecker IChartServiceChecker, dashboardChecker IDashboardChecker) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		handleAVCheck(w, req, uiBackendChecker, chartServiceChecker, dashboardChecker)
	}
}

func handleAVCheck(w http.ResponseWriter, req *http.Request, uiBackendChecker IUIBackendChecker, chartServiceChecker IChartServiceChecker, dashboardChecker IDashboardChecker) {
	log := logUtils.GetLogger(req.Context())

	// only check if the query param is present, do not check the actual value of it
	_, completeCheck := req.URL.Query()[completeCheckKey]

	body := response{}

	uiBackendCheckResult := uiBackendChecker.RunChecks(req.Context())
	body.UIBackend = uiBackendCheckResult

	dashboardCheckResult := dashboardChecker.RunChecks(req.Context())
	body.Dashboard = dashboardCheckResult

	chartServiceCheckResult := chartServiceChecker.RunChecks(req.Context(), completeCheck)
	body.ChartService = chartServiceCheckResult

	marshaledBody, err := json.Marshal(body)
	if err != nil {
		errMsg := fmt.Sprintf("cannot marshal response body: %s", err.Error())
		log.Error(errors.New(errMsg))
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	var responseCode int
	if uiBackendCheckResult.AllChecksSuccessful() && dashboardCheckResult.AllChecksSuccessful() && chartServiceCheckResult.AllChecksSuccessful() {
		responseCode = http.StatusOK
	} else {
		responseCode = http.StatusInternalServerError
		log.WithField(logUtils.LogKeyResponseBody, body).Warn("avcheck failed")
	}

	w.WriteHeader(responseCode)
	_, err = w.Write(marshaledBody)
	if err != nil {
		log.Error(errors.Wrap(err, "cannot write response"))
	}
}

func HandlePing(w http.ResponseWriter, req *http.Request) {
	log := logUtils.GetLogger(req.Context())
	log.Debugln("Live/Readiness probe called")

	w.WriteHeader(200)
	_, err := w.Write([]byte("pong"))
	if err != nil {
		log.Error(errors.Wrap(err, "cannot write response"))
	}
}
