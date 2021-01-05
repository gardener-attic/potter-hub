package avcheck

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	logUtils "github.com/gardener/potter-hub/pkg/log"
	"github.com/gardener/potter-hub/pkg/util"
)

const (
	allChartsPath = "v1/charts"
)

type ChartServiceCheckResult struct {
	Ping            *PingResult            `json:"ping,omitempty"`
	ChartsAvailable *ChartsAvailableResult `json:"chartsAvailable,omitempty"`
}

func (result *ChartServiceCheckResult) AllChecksSuccessful() bool {
	checksSuccessful := []bool{}

	checksSuccessful = append(checksSuccessful, result.Ping.CheckSuccessful)

	if result.ChartsAvailable != nil {
		checksSuccessful = append(checksSuccessful, result.ChartsAvailable.CheckSuccessful)
	}

	return util.AreAllItemsTrue(checksSuccessful)
}

type ChartsAvailableResult struct {
	LastCheckTime   time.Time `json:"lastCheckTime"`
	CheckSuccessful bool      `json:"checkSuccessful"`
}

type IChartServiceChecker interface {
	RunChecks(context.Context, bool) *ChartServiceCheckResult
	StartChartsAvailableCheckBackgroundJob()
}

func NewChartServiceChecker(baseURL string, config *Configuration) IChartServiceChecker {
	instance := &chartServiceChecker{
		baseURL: util.GetEnvOrPanic("CHARTSVC_URL"),
		config:  config,
		chartsAvailableResult: &ChartsAvailableResult{
			LastCheckTime:   time.Now(),
			CheckSuccessful: false,
		},
	}

	return instance
}

type chartServiceChecker struct {
	baseURL                  string
	config                   *Configuration
	chartsAvailableResult    *ChartsAvailableResult
	chartsAvailableResultMux sync.Mutex
}

func (c *chartServiceChecker) RunChecks(ctx context.Context, completeCheck bool) *ChartServiceCheckResult {
	result := &ChartServiceCheckResult{}

	livenessURL := fmt.Sprintf("%s/%s", c.baseURL, livenessProbePath)
	result.Ping = checkPing(ctx, livenessURL)

	if !completeCheck {
		return result
	}

	result.ChartsAvailable = c.getChartsAvailableResult()

	return result
}

func (c *chartServiceChecker) setChartsAvailableResult(newResult *ChartsAvailableResult) {
	c.chartsAvailableResultMux.Lock()
	c.chartsAvailableResult = newResult
	c.chartsAvailableResultMux.Unlock()
}

func (c *chartServiceChecker) getChartsAvailableResult() *ChartsAvailableResult {
	c.chartsAvailableResultMux.Lock()
	result := c.chartsAvailableResult
	c.chartsAvailableResultMux.Unlock()
	return result
}

func (c *chartServiceChecker) StartChartsAvailableCheckBackgroundJob() {
	innerLogger := logUtils.StandardLogger().WithFields(logrus.Fields{
		logUtils.LogKeyLoggerName: "AVCheck.ChartsAvailableCheckBackgroundJob",
	})
	logger := &logUtils.Logger{Entry: innerLogger}
	allChartsURL := fmt.Sprintf("%s/%s", c.baseURL, allChartsPath)

	for {
		result := checkChartsAvailable(allChartsURL, logger)
		c.setChartsAvailableResult(result)
		time.Sleep(c.config.ChartsAvailableCheckInterval)
	}
}

func checkChartsAvailable(allChartsURL string, logger *logUtils.Logger) *ChartsAvailableResult {
	result := &ChartsAvailableResult{
		LastCheckTime:   time.Now(),
		CheckSuccessful: false,
	}

	req, err := http.NewRequest("GET", allChartsURL, nil)
	if err != nil {
		logger.Error(errors.Wrap(err, "cannot create http request"))
		return result
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error(errors.Wrap(err, "http request failed"))
		return result
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errMsg := fmt.Sprintf("get charts failed: returned with status code %d", resp.StatusCode)
		logger.Error(errors.New(errMsg))
		return result
	}

	getChartsResponse := map[string]interface{}{}

	err = json.NewDecoder(resp.Body).Decode(&getChartsResponse)
	if err != nil {
		logger.Error(errors.Wrap(err, "cannot unmarshal get charts response"))
		return result
	}

	data, ok := getChartsResponse["data"]
	if !ok {
		logger.Error(errors.New("get charts response does not contain the data attribute"))
		return result
	}

	charts, ok := data.([]interface{})
	if !ok {
		logger.Error(errors.New("data attribute of the get charts response is not of type list"))
		return result
	}

	if len(charts) == 0 {
		logger.Error(errors.New("get charts response doesn't contain any charts"))
		return result
	}

	result.CheckSuccessful = true
	return result
}
