package cortex

/***********************************************************************************************************************
 * Integrations - Generic structs for integrations
 **********************************************************************************************************************/

// CatalogEntityGit represents the Git metadata around a catalog entity
// @see https://docs.cortex.io/docs/reference/basics/entities#example-cortexyaml-for-service-entity
type CatalogEntityGit struct {
	Github    CatalogEntityGitGithub      `json:"github,omitempty" yaml:"github,omitempty"`
	Gitlab    CatalogEntityGitGitlab      `json:"gitlab,omitempty" yaml:"gitlab,omitempty"`
	Azure     CatalogEntityGitAzureDevOps `json:"azure,omitempty" yaml:"azure,omitempty"`
	BitBucket CatalogEntityGitBitBucket   `json:"bitbucket,omitempty" yaml:"bitbucket,omitempty"`
}

type CatalogEntityOnCall struct {
	Pagerduty CatalogEntityOnCallPagerduty `json:"pagerduty,omitempty" yaml:"pagerduty,omitempty"`
	OpsGenie  CatalogEntityOnCallOpsGenie  `json:"opsgenie,omitempty" yaml:"opsgenie,omitempty"`
	VictorOps CatalogEntityOnCallVictorOps `json:"victorops,omitempty" yaml:"victorops,omitempty"`
}

/***********************************************************************************************************************
 * GitHub - https://docs.cortex.io/docs/reference/integrations/github
 **********************************************************************************************************************/

type CatalogEntityGitGithub struct {
	Repository string `json:"repository" yaml:"repository"`
	BasePath   string `json:"basepath,omitempty" yaml:"basepath,omitempty"`
}

/***********************************************************************************************************************
 * GitLab - https://docs.cortex.io/docs/reference/integrations/gitlab
 **********************************************************************************************************************/

type CatalogEntityGitGitlab struct {
	Repository string `json:"repository" yaml:"repository"`
	BasePath   string `json:"basepath,omitempty" yaml:"basepath,omitempty"`
}

/***********************************************************************************************************************
 * Azure DevOps - https://docs.cortex.io/docs/reference/integrations/azuredevops
 **********************************************************************************************************************/

type CatalogEntityGitAzureDevOps struct {
	Project    string `json:"project" yaml:"project"`
	Repository string `json:"repository" yaml:"repository"`
	BasePath   string `json:"basepath,omitempty" yaml:"basepath,omitempty"`
}

/***********************************************************************************************************************
 * BitBucket - https://docs.cortex.io/docs/reference/integrations/bitbucket
 **********************************************************************************************************************/

type CatalogEntityGitBitBucket struct {
	Repository string `json:"repository" yaml:"repository"`
}

/***********************************************************************************************************************
 * PagerDuty - https://docs.cortex.io/docs/reference/integrations/pagerduty
 **********************************************************************************************************************/

type CatalogEntityOnCallPagerduty struct {
	ID   string `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}

/***********************************************************************************************************************
 * OpsGenie -  https://docs.cortex.io/docs/reference/integrations/opsgenie
 **********************************************************************************************************************/

type CatalogEntityOnCallOpsGenie struct {
	ID   string `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}

/***********************************************************************************************************************
 * VictorOps -  https://docs.cortex.io/docs/reference/integrations/victorops
 **********************************************************************************************************************/

type CatalogEntityOnCallVictorOps struct {
	ID   string `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}
