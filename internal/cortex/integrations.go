package cortex

/***********************************************************************************************************************
 * Catalog Entity Descriptor
 **********************************************************************************************************************/

// CatalogEntityData is a struct used from YAML-based data, since its structure does not
// match the structure of the CatalogEntity struct in other responses.
// See: https://github.com/cortexapps/solutions/blob/master/examples/yaml/catalog/resource.yaml
type CatalogEntityData struct {
	Title        string                    `json:"title" yaml:"title"`
	Description  string                    `json:"description,omitempty" yaml:"description,omitempty"`
	Tag          string                    `json:"x-cortex-tag" yaml:"x-cortex-tag"`
	Type         string                    `json:"x-cortex-type,omitempty" yaml:"x-cortex-type,omitempty"`
	Definition   map[string]interface{}    `json:"x-cortex-definition,omitempty" yaml:"x-cortex-definition,omitempty"`
	Owners       []CatalogEntityOwner      `json:"x-cortex-owners,omitempty" yaml:"x-cortex-owners,omitempty"`
	Groups       []string                  `json:"x-cortex-groups,omitempty" yaml:"x-cortex-groups"` // TODO: is this -groups or -service-groups? docs unclear
	Links        []CatalogEntityLink       `json:"x-cortex-link,omitempty" yaml:"x-cortex-link,omitempty"`
	Metadata     map[string]interface{}    `json:"x-cortex-custom-metadata,omitempty" yaml:"x-cortex-custom-metadata,omitempty"`
	Dependencies []CatalogEntityDependency `json:"x-cortex-dependency,omitempty" yaml:"x-cortex-dependency,omitempty"`

	// Various generic integration attributes
	Alerts         []CatalogEntityAlert        `json:"x-cortex-alerts,omitempty" yaml:"x-cortex-alerts,omitempty"`
	Apm            CatalogEntityApm            `json:"x-cortex-apm,omitempty" yaml:"x-cortex-apm,omitempty"`
	Dashboards     CatalogEntityDashboards     `json:"x-cortex-dashboards,omitempty" yaml:"x-cortex-dashboards,omitempty"`
	Git            CatalogEntityGit            `json:"x-cortex-git,omitempty" yaml:"x-cortex-git,omitempty"`
	Issues         CatalogEntityIssues         `json:"x-cortex-issues,omitempty" yaml:"x-cortex-issues,omitempty"`
	OnCall         CatalogEntityOnCall         `json:"x-cortex-oncall,omitempty" yaml:"x-cortex-oncall,omitempty"`
	SLOs           CatalogEntitySLOs           `json:"x-cortex-slos,omitempty" yaml:"x-cortex-slos,omitempty"`
	StaticAnalysis CatalogEntityStaticAnalysis `json:"x-cortex-static-analysis,omitempty" yaml:"x-cortex-static-analysis,omitempty"`

	// Integration-specific things
	BugSnag   CatalogEntityBugSnag   `json:"x-cortex-bugsnag,omitempty" yaml:"x-cortex-bugsnag,omitempty"`
	Checkmarx CatalogEntityCheckmarx `json:"x-cortex-checkmarx,omitempty" yaml:"x-cortex-checkmarx,omitempty"`
	Rollbar   CatalogEntityRollbar   `json:"x-cortex-rollbar,omitempty" yaml:"x-cortex-rollbar,omitempty"`
	Sentry    CatalogEntitySentry    `json:"x-cortex-sentry,omitempty" yaml:"x-cortex-sentry,omitempty"`
	Snyk      CatalogEntitySnyk      `json:"x-cortex-snyk,omitempty" yaml:"x-cortex-snyk,omitempty"`
}

type CatalogEntityLink struct {
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"` // runbook, documentation, logs, dashboard, metrics, healthcheck
	Url  string `json:"url" yaml:"url"`
}

type CatalogEntityDependency struct {
	Tag         string                 `json:"tag" yaml:"tag"`
	Method      string                 `json:"method,omitempty" yaml:"method,omitempty"`
	Path        string                 `json:"path,omitempty" yaml:"path,omitempty"`
	Description string                 `json:"description,omitempty" yaml:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

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
	Type string `json:"type" yaml:"type"` // <datadog | grafana | newrelic>
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
	Type                 string `json:"type" yaml:"type"`                     // group, user, slack
	Name                 string `json:"name,omitempty" yaml:"name,omitempty"` // Must be of form <org>/<team>
	Description          string `json:"description,omitempty" yaml:"description,omitempty"`
	Provider             string `json:"provider,omitempty" yaml:"provider,omitempty"`
	Channel              string `json:"channel,omitempty" yaml:"channel,omitempty"` // for slack, do not add # to beginning
	NotificationsEnabled bool   `json:"notificationsEnabled,omitempty" yaml:"notificationsEnabled,omitempty"`
}

type CatalogEntitySLOs struct {
	DataDog    []CatalogEntitySLODataDog         `json:"datadog,omitempty" yaml:"datadog,omitempty"`
	Dynatrace  []CatalogEntitySLODynatrace       `json:"dynatrace,omitempty" yaml:"dynatrace,omitempty"`
	Lightstep  CatalogEntitySLOLightstep         `json:"lightstep,omitempty" yaml:"lightstep,omitempty"`
	Prometheus []CatalogEntitySLOPrometheusQuery `json:"prometheus,omitempty" yaml:"prometheus,omitempty"`
	SignalFX   []CatalogEntitySLOSignalFX        `json:"signalfx,omitempty" yaml:"signalfx,omitempty"`
	SumoLogic  []CatalogEntitySLOSumoLogic       `json:"sumologic,omitempty" yaml:"sumologic,omitempty"`
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
	Labels     []string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Components []string `json:"components,omitempty" yaml:"components,omitempty"`
}

/***********************************************************************************************************************
 * LightStep - https://docs.cortex.io/docs/reference/integrations/lightstep
 **********************************************************************************************************************/

/**
x-cortex-slos:
  lightstep:
    - streamId: sc4jmdXT
      targets:
		  latency:
			- percentile: 0.5
			  target: 2
			  slo: 0.9995
			- percentile: 0.7
			  target: 1
			  slo: 0.9998
*/

type CatalogEntitySLOLightstep struct {
	Streams []CatalogEntitySLOLightstepStream `json:"stream" yaml:"stream"`
}

type CatalogEntitySLOLightstepStream struct {
	StreamID string                          `json:"streamId" yaml:"streamId"`
	Targets  CatalogEntitySLOLightstepTarget `json:"targets" yaml:"targets"`
}

type CatalogEntitySLOLightstepTarget struct {
	Latencies []CatalogEntitySLOLightstepTargetLatency `json:"latency" yaml:"latency"`
}

type CatalogEntitySLOLightstepTargetLatency struct {
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

type CatalogEntitySLOPrometheusQuery struct {
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
	Lookback  string `json:"lookback" yaml:"lookback"`
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
	Organization string `json:"organizationId" yaml:"organizationId"`
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
