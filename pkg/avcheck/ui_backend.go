package avcheck

import (
	"context"
	"fmt"
)

type UIBackendCheckResult struct {
	Ping *PingResult `json:"ping,omitempty"`
}

func (h *UIBackendCheckResult) AllChecksSuccessful() bool {
	return h.Ping.CheckSuccessful
}

type IUIBackendChecker interface {
	RunChecks(context.Context) *UIBackendCheckResult
}

func NewUIBackendChecker(baseURL string) IUIBackendChecker {
	instance := &uiBackendChecker{
		baseURL: baseURL,
	}

	return instance
}

type uiBackendChecker struct {
	baseURL string
}

func (h *uiBackendChecker) RunChecks(ctx context.Context) *UIBackendCheckResult {
	result := &UIBackendCheckResult{}
	livenessURL := fmt.Sprintf("%s/%s", h.baseURL, livenessProbePath)
	result.Ping = checkPing(ctx, livenessURL)
	return result
}
