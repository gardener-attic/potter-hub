package avcheck

import "context"

type DashboardCheckResult struct {
	Ping *PingResult `json:"ping,omitempty"`
}

func (d *DashboardCheckResult) AllChecksSuccessful() bool {
	return d.Ping.CheckSuccessful
}

type IDashboardChecker interface {
	RunChecks(context.Context) *DashboardCheckResult
}

func NewDashboardChecker(baseURL string) IDashboardChecker {
	instance := &dashboardChecker{
		baseURL: baseURL,
	}

	return instance
}

type dashboardChecker struct {
	baseURL string
}

func (d *dashboardChecker) RunChecks(ctx context.Context) *DashboardCheckResult {
	result := &DashboardCheckResult{}
	result.Ping = checkPing(ctx, d.baseURL)
	return result
}
