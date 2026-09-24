package cortex

// PagerdutyConfigurationRequest is the body of POST and PUT /api/v1/pagerduty/configuration. PUT replaces the whole
// configuration, so it takes the same body as POST.
type PagerdutyConfigurationRequest struct {
	Token           string `json:"token"`
	IsTokenReadonly bool   `json:"isTokenReadonly"`
}

// PagerdutyConfiguration is the PagerDuty configuration. The API returns only the last four characters of the token.
type PagerdutyConfiguration struct {
	LastFour        string `json:"lastFour"`
	IsTokenReadonly bool   `json:"isTokenReadonly"`
}

func (c *HttpClient) PagerdutyConfiguration() *SingleInstanceClient[PagerdutyConfigurationRequest, PagerdutyConfiguration] {
	return &SingleInstanceClient[PagerdutyConfigurationRequest, PagerdutyConfiguration]{client: c, domain: "pagerduty", name: "pagerduty"}
}
