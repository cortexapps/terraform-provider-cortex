package cortex

/***********************************************************************************************************************
 * Integrations - Generic structs for integrations
 **********************************************************************************************************************/

type CatalogEntityAlert struct {
	Type  string `json:"type" yaml:"type"`
	Tag   string `json:"tag" yaml:"tag"`
	Value string `json:"value" yaml:"value"`
}

type CatalogEntityApm struct {
	DataDog   CatalogEntityApmDataDog   `json:"datadog,omitempty" yaml:"datadog,omitempty"`
	Dynatrace CatalogEntityApmDynatrace `json:"dynatrace,omitempty" yaml:"dynatrace,omitempty"`
	NewRelic  CatalogEntityApmNewRelic  `json:"newrelic,omitempty" yaml:"newrelic,omitempty"`
}

type CatalogEntityDashboards struct {
	Embeds []CatalogEntityDashboardsEmbed `json:"embeds,omitempty" yaml:"embeds,omitempty"`
}

type CatalogEntityDashboardsEmbed struct {
	Type string `json:"type" yaml:"type"`
	URL  string `json:"url" yaml:"url"`
}

// CatalogEntityGit represents the Git metadata around a catalog entity
// @see https://docs.cortex.io/docs/reference/basics/entities#example-cortexyaml-for-service-entity
type CatalogEntityGit struct {
	Github    CatalogEntityGitGithub      `json:"github,omitempty" yaml:"github,omitempty"`
	Gitlab    CatalogEntityGitGitlab      `json:"gitlab,omitempty" yaml:"gitlab,omitempty"`
	Azure     CatalogEntityGitAzureDevOps `json:"azure,omitempty" yaml:"azure,omitempty"`
	BitBucket CatalogEntityGitBitBucket   `json:"bitbucket,omitempty" yaml:"bitbucket,omitempty"`
}

type CatalogEntityIssues struct {
	Jira CatalogEntityIssuesJira `json:"jira,omitempty" yaml:"jira,omitempty"`
}

type CatalogEntityOnCall struct {
	Pagerduty CatalogEntityOnCallPagerduty `json:"pagerduty,omitempty" yaml:"pagerduty,omitempty"`
	OpsGenie  CatalogEntityOnCallOpsGenie  `json:"opsgenie,omitempty" yaml:"opsgenie,omitempty"`
	VictorOps CatalogEntityOnCallVictorOps `json:"victorops,omitempty" yaml:"victorops,omitempty"`
}

type CatalogEntityOwner struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Type        string `json:"type" yaml:"type"`
	Provider    string `json:"provider" yaml:"provider"`
}

type CatalogEntitySLOs struct {
	DataDog    []CatalogEntitySLODataDog    `json:"datadog,omitempty" yaml:"datadog,omitempty"`
	Dynatrace  []CatalogEntitySLODynatrace  `json:"dynatrace,omitempty" yaml:"dynatrace,omitempty"`
	LightStep  []CatalogEntitySLOLightStep  `json:"lightstep,omitempty" yaml:"lightstep,omitempty"`
	Prometheus []CatalogEntitySLOPrometheus `json:"prometheus,omitempty" yaml:"prometheus,omitempty"`
	SignalFX   []CatalogEntitySLOSignalFX   `json:"signalfx,omitempty" yaml:"signalfx,omitempty"`
	SumoLogic  []CatalogEntitySLOSumoLogic  `json:"sumologic,omitempty" yaml:"sumologic,omitempty"`
}

type CatalogEntityStaticAnalysis struct {
	CodeCov   CatalogEntityStaticAnalysisCodeCov   `json:"codecov,omitempty" yaml:"codecov,omitempty"`
	Mend      CatalogEntityStaticAnalysisMend      `json:"mend,omitempty" yaml:"mend,omitempty"`
	SonarQube CatalogEntityStaticAnalysisSonarQube `json:"sonarqube,omitempty" yaml:"sonarqube,omitempty"`
	Veracode  CatalogEntityStaticAnalysisVeracode  `json:"veracode,omitempty" yaml:"veracode,omitempty"`
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
 * BugSnag - https://docs.cortex.io/docs/reference/integrations/bugsnag
 **********************************************************************************************************************/

type CatalogEntityBugSnag struct {
	Project string `json:"project" yaml:"project"`
}

/***********************************************************************************************************************
 * Checkmarx - https://docs.cortex.io/docs/reference/integrations/checkmarx
 **********************************************************************************************************************/

type CatalogEntityCheckmarx struct {
	Projects []string `json:"projects" yaml:"projects"`
}

type CatalogEntityCheckmarxProject struct {
	ProjectID   string `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	ProjectName string `json:"projectName,omitempty" yaml:"projectName,omitempty"`
}

/***********************************************************************************************************************
 * CodeCov - https://docs.cortex.io/docs/reference/integrations/codecov
 **********************************************************************************************************************/

type CatalogEntityStaticAnalysisCodeCov struct {
	Repository string `json:"repo" yaml:"repo"`
	Provider   string `json:"provider" yaml:"provider"`
}

/***********************************************************************************************************************
 * DataDog - https://docs.cortex.io/docs/reference/integrations/datadog
 **********************************************************************************************************************/

type CatalogEntityApmDataDog struct {
	Monitors []int `json:"monitors" yaml:"monitors"`
}

type CatalogEntitySLODataDog struct {
	ID string `json:"id" yaml:"id"`
}

/***********************************************************************************************************************
 * Dynatrace - https://docs.cortex.io/docs/reference/integrations/dynatrace
 **********************************************************************************************************************/

type CatalogEntityApmDynatrace struct {
	EntityIDs          []string `json:"entityIds,omitempty" yaml:"entityIds,omitempty"`
	EntityNameMatchers []string `json:"entityNameMatchers,omitempty" yaml:"entityNameMatchers,omitempty"`
}

type CatalogEntitySLODynatrace struct {
	ID string `json:"id" yaml:"id"`
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
 * JIRA - https://docs.cortex.io/docs/reference/integrations/jira
 **********************************************************************************************************************/

type CatalogEntityIssuesJira struct {
	DefaultJQL string   `json:"defaultJql" yaml:"defaultJql"`
	Projects   []string `json:"projects,omitempty" yaml:"projects,omitempty"`
	Components []string `json:"components,omitempty" yaml:"components,omitempty"`
}

/***********************************************************************************************************************
 * LightStep - https://docs.cortex.io/docs/reference/integrations/lightstep
 **********************************************************************************************************************/

type CatalogEntitySLOLightStep struct {
	StreamID string                             `json:"streamId" yaml:"streamId"`
	Latency  []CatalogEntitySLOLightStepLatency `json:"latency" yaml:"latency"`
}

type CatalogEntitySLOLightStepLatency struct {
	Percentile float64 `json:"percentile" yaml:"percentile"`
	Target     int     `json:"target" yaml:"target"`
	SLO        float64 `json:"slo" yaml:"slo"`
}

/***********************************************************************************************************************
 * Mend - https://docs.cortex.io/docs/reference/integrations/mend
 **********************************************************************************************************************/

type CatalogEntityStaticAnalysisMend struct {
	ApplicationIDs []string `json:"applicationIds,omitempty" yaml:"applicationIds,omitempty"`
	ProjectIDs     []string `json:"projectIds,omitempty" yaml:"projectIds,omitempty"`
}

/***********************************************************************************************************************
 * New Relic - https://docs.cortex.io/docs/reference/integrations/newrelic
 **********************************************************************************************************************/

type CatalogEntityApmNewRelic struct {
	ApplicationID int    `json:"applicationId" yaml:"applicationId"`
	Alias         string `json:"alias,omitempty" yaml:"alias,omitempty"`
}

/***********************************************************************************************************************
 * PagerDuty - https://docs.cortex.io/docs/reference/integrations/pagerduty
 **********************************************************************************************************************/

type CatalogEntityOnCallPagerduty struct {
	ID   string `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}

/***********************************************************************************************************************
 * Prometheus - https://docs.cortex.io/docs/reference/integrations/prometheus
 **********************************************************************************************************************/

type CatalogEntitySLOPrometheus struct {
	ErrorQuery string  `json:"errorQuery" yaml:"errorQuery"`
	TotalQuery string  `json:"totalQuery" yaml:"totalQuery"`
	SLO        float64 `json:"slo" yaml:"slo"`
	Name       string  `json:"name,omitempty" yaml:"name,omitempty"`
	Alias      string  `json:"alias,omitempty" yaml:"alias,omitempty"`
}

/***********************************************************************************************************************
 * Rollbar - https://docs.cortex.io/docs/reference/integrations/rollbar
 **********************************************************************************************************************/

type CatalogEntityRollbar struct {
	Project string `json:"project" yaml:"project"`
}

/***********************************************************************************************************************
 * SignalFX - https://docs.cortex.io/docs/reference/integrations/signalfx
 **********************************************************************************************************************/

type CatalogEntitySLOSignalFX struct {
	Query     string `json:"query" yaml:"query"`
	Rollup    string `json:"rollup" yaml:"rollup"`
	Target    int    `json:"target" yaml:"target"`
	LookBack  string `json:"lookback" yaml:"lookback"`
	Operation string `json:"operation" yaml:"operation"`
}

/***********************************************************************************************************************
 * Sentry - https://docs.cortex.io/docs/reference/integrations/sentry
 **********************************************************************************************************************/

type CatalogEntitySentry struct {
	Project string `json:"project" yaml:"project"`
}

/***********************************************************************************************************************
 * Snyk - https://docs.cortex.io/docs/reference/integrations/snyk
 **********************************************************************************************************************/

type CatalogEntitySnyk struct {
	Projects []CatalogEntitySnykProject `json:"projects,omitempty" yaml:"projects,omitempty"`
}

type CatalogEntitySnykProject struct {
	Organization string `json:"organization" yaml:"organization"`
	ProjectID    string `json:"projectId" yaml:"projectId"`
}

/***********************************************************************************************************************
 * SonarQube - https://docs.cortex.io/docs/reference/integrations/sonarqube
 **********************************************************************************************************************/

type CatalogEntityStaticAnalysisSonarQube struct {
	Project string `json:"project" yaml:"project"`
}

/***********************************************************************************************************************
 * SumoLogic - https://docs.cortex.io/docs/reference/integrations/sumologic
 **********************************************************************************************************************/

type CatalogEntitySLOSumoLogic struct {
	ID string `json:"id" yaml:"id"`
}

/***********************************************************************************************************************
 * OpsGenie -  https://docs.cortex.io/docs/reference/integrations/opsgenie
 **********************************************************************************************************************/

type CatalogEntityOnCallOpsGenie struct {
	ID   string `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}

/***********************************************************************************************************************
 * Veracode -  https://docs.cortex.io/docs/reference/integrations/veracode
 **********************************************************************************************************************/

type CatalogEntityStaticAnalysisVeracode struct {
	ApplicationNames []string                                     `json:"applicationNames,omitempty" yaml:"applicationNames,omitempty"`
	Sandboxes        []CatalogEntityStaticAnalysisVeracodeSandbox `json:"sandboxes,omitempty" yaml:"sandboxes,omitempty"`
}

type CatalogEntityStaticAnalysisVeracodeSandbox struct {
	ApplicationName string `json:"applicationName" yaml:"applicationName"`
	SandboxName     string `json:"sandboxName" yaml:"sandboxName"`
}

/***********************************************************************************************************************
 * VictorOps -  https://docs.cortex.io/docs/reference/integrations/victorops
 **********************************************************************************************************************/

type CatalogEntityOnCallVictorOps struct {
	ID   string `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}
