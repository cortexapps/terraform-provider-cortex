package cortex

type CreateGitlabConfigurationRequest struct {
	Alias                string   `json:"alias"`
	IsDefault            bool     `json:"isDefault"`
	Host                 string   `json:"host,omitempty"`
	PersonalAccessToken  string   `json:"personalAccessToken"`
	HidePersonalProjects bool     `json:"hidePersonalProjects"`
	GroupNames           []string `json:"groupNames"`
}

// UpdateGitlabConfigurationRequest has no host: the API ignores a host on update.
type UpdateGitlabConfigurationRequest struct {
	Alias                string   `json:"alias"`
	IsDefault            bool     `json:"isDefault"`
	HidePersonalProjects bool     `json:"hidePersonalProjects"`
	GroupNames           []string `json:"groupNames"`
	PersonalAccessToken  string   `json:"personalAccessToken,omitempty"`
}

// GitlabConfiguration is a GitLab configuration. The API returns only the last four characters of the token.
type GitlabConfiguration struct {
	Alias                string   `json:"alias"`
	Host                 string   `json:"host,omitempty"`
	IsDefault            bool     `json:"isDefault"`
	LastFour             string   `json:"lastFour"`
	HidePersonalProjects bool     `json:"hidePersonalProjects"`
	GroupNames           []string `json:"groupNames"`
}

func (c GitlabConfiguration) GetAlias() string { return c.Alias }

func (c *HttpClient) GitlabConfigurations() *MultiInstanceClient[CreateGitlabConfigurationRequest, UpdateGitlabConfigurationRequest, GitlabConfiguration] {
	return &MultiInstanceClient[CreateGitlabConfigurationRequest, UpdateGitlabConfigurationRequest, GitlabConfiguration]{client: c, domain: "gitlab", name: "gitlab"}
}
