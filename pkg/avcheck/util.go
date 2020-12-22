package avcheck

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	logUtils "github.wdf.sap.corp/kubernetes/hub/pkg/log"
)

const (
	livenessProbePath = "live"
)

type PingResult struct {
	CheckSuccessful bool `json:"checkSuccessful"`
}

func checkPing(ctx context.Context, url string) *PingResult {
	logger := logUtils.GetLogger(ctx)
	logger.Infof("initiating request to url %s", url)

	result := &PingResult{
		CheckSuccessful: false,
	}

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		logger.Error(errors.Wrapf(err, "request to url %s failed", url))
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		errMsg := fmt.Sprintf("ping failed for url %s: returned with status code %d", url, resp.StatusCode)
		logger.Error(errors.New(errMsg))
		return result
	}

	logger.Infof("request to url %s was successful", url)

	result.CheckSuccessful = true
	return result
}
