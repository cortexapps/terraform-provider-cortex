package integrations

import "github.com/cortexapps/terraform-provider-cortex/internal/cortex"

// PagerDutyConfigurationRequest is the body of POST and PUT /api/v1/pagerduty/configuration. PUT replaces the whole
// configuration, so it takes the same body as POST.
type PagerDutyConfigurationRequest struct {
	Token           string `json:"token"`
	IsTokenReadonly bool   `json:"isTokenReadonly"`
}

// PagerDutyConfiguration is the PagerDuty configuration. The API returns only the last four characters of the token.
type PagerDutyConfiguration struct {
	LastFour        string `json:"lastFour"`
	IsTokenReadonly bool   `json:"isTokenReadonly"`
}

// PagerDuty returns the client for the PagerDuty configuration.
func PagerDuty(c *cortex.HttpClient) SingleInstanceClientInterface[PagerDutyConfigurationRequest, PagerDutyConfiguration] {
	return &SingleInstanceClient[PagerDutyConfigurationRequest, PagerDutyConfiguration]{client: c, base: "/api/v1/pagerduty/", name: "pagerduty"}
}
