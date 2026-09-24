package cortex

// Values of the Jira "type" property. The API picks the request shape from it.
const (
	JiraTypeCloudBasic  = "CLOUD_BASIC"
	JiraTypeCloudScoped = "CLOUD_SCOPED"
	JiraTypeOnPremBasic = "ON_PREM_BASIC"
)

// JiraConfigurationRequest is the create and update body. Each type uses a subset of the fields; the others stay
// empty and are not sent. BaseUrl takes the wire values "jira.com", "atlassian.net", or "api.atlassian.com/ex/jira".
type JiraConfigurationRequest struct {
	Type         string `json:"type"`
	Alias        string `json:"alias"`
	IsDefault    bool   `json:"isDefault"`
	Subdomain    string `json:"subdomain,omitempty"`
	Email        string `json:"email,omitempty"`
	ApiToken     string `json:"apiToken,omitempty"`
	BaseUrl      string `json:"baseUrl,omitempty"`
	CloudId      string `json:"cloudId,omitempty"`
	Host         string `json:"host,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	FrontendHost string `json:"frontendHost,omitempty"`
}

// JiraConfiguration is a Jira configuration. The API never returns the API token, the password, or the cloud ID.
type JiraConfiguration struct {
	Alias        string `json:"alias"`
	IsDefault    bool   `json:"isDefault"`
	Type         string `json:"type"`
	Subdomain    string `json:"subdomain,omitempty"`
	BaseUrl      string `json:"baseUrl,omitempty"`
	Host         string `json:"host,omitempty"`
	FrontendHost string `json:"frontendHost,omitempty"`
	LastFour     string `json:"lastFour,omitempty"`
	Email        string `json:"email,omitempty"`
	Username     string `json:"username,omitempty"`
}

func (c JiraConfiguration) GetAlias() string { return c.Alias }
