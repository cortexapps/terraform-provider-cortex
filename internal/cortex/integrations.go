package cortex

/***********************************************************************************************************************
 * Catalog Entity Descriptor
 **********************************************************************************************************************/

// CatalogEntityData is a struct used from YAML-based data, since its structure does not
// match the structure of the CatalogEntity struct in other responses.
// See: https://github.com/cortexapps/solutions/blob/master/examples/yaml/catalog/resource.yaml
type CatalogEntityData struct {
	Title          string                      `json:"title" yaml:"title"`
	Description    string                      `json:"description,omitempty" yaml:"description,omitempty"`
	Tag            string                      `json:"x-cortex-tag" yaml:"x-cortex-tag"`
	Type           string                      `json:"x-cortex-type,omitempty" yaml:"x-cortex-type,omitempty"`
	Definition     map[string]interface{}      `json:"x-cortex-definition,omitempty" yaml:"x-cortex-definition,omitempty"`
	Owners         []CatalogEntityOwner        `json:"x-cortex-owners,omitempty" yaml:"x-cortex-owners,omitempty"`
	Children       []CatalogEntityChild        `json:"x-cortex-children,omitempty" yaml:"x-cortex-children,omitempty"`
	DomainParents  []CatalogEntityDomainParent `json:"x-cortex-domain-parents,omitempty" yaml:"x-cortex-domain-parents,omitempty"`
	Groups         []string                    `json:"x-cortex-groups,omitempty" yaml:"x-cortex-groups,omitempty"` // TODO: is this -groups or -service-groups? docs unclear
	Links          []CatalogEntityLink         `json:"x-cortex-link,omitempty" yaml:"x-cortex-link,omitempty"`
	IgnoreMetadata bool                        `json:"-" yaml:"-"`
	Metadata       map[string]interface{}      `json:"x-cortex-custom-metadata,omitempty" yaml:"x-cortex-custom-metadata,omitempty"`
	Dependencies   []CatalogEntityDependency   `json:"x-cortex-dependency,omitempty" yaml:"x-cortex-dependency,omitempty"`

	// Various generic integration attributes
	Alerts         []CatalogEntityAlert        `json:"x-cortex-alerts,omitempty" yaml:"x-cortex-alerts,omitempty"`
	Apm            CatalogEntityApm            `json:"x-cortex-apm,omitempty" yaml:"x-cortex-apm,omitempty"`
	Dashboards     CatalogEntityDashboards     `json:"x-cortex-dashboards,omitempty" yaml:"x-cortex-dashboards,omitempty"`
	Git            CatalogEntityGit            `json:"x-cortex-git,omitempty" yaml:"x-cortex-git,omitempty"`
	Issues         CatalogEntityIssues         `json:"x-cortex-issues,omitempty" yaml:"x-cortex-issues,omitempty"`
	OnCall         CatalogEntityOnCall         `json:"x-cortex-oncall,omitempty" yaml:"x-cortex-oncall,omitempty"`
	SLOs           CatalogEntitySLOs           `json:"x-cortex-slos,omitempty" yaml:"x-cortex-slos,omitempty"`
	StaticAnalysis CatalogEntityStaticAnalysis `json:"x-cortex-static-analysis,omitempty" yaml:"x-cortex-static-analysis,omitempty"`
	CiCd           CatalogEntityCiCd           `json:"x-cortex-ci-cd,omitempty" yaml:"x-cortex-ci-cd,omitempty"`

	// Integration-specific things
	BugSnag        CatalogEntityBugSnag         `json:"x-cortex-bugsnag,omitempty" yaml:"x-cortex-bugsnag,omitempty"`
	Checkmarx      CatalogEntityCheckmarx       `json:"x-cortex-checkmarx,omitempty" yaml:"x-cortex-checkmarx,omitempty"`
	Coralogix      CatalogEntityCoralogix       `json:"x-cortex-coralogix,omitempty" yaml:"x-cortex-coralogix,omitempty"`
	FireHydrant    CatalogEntityFireHydrant     `json:"x-cortex-firehydrant,omitempty" yaml:"x-cortex-firehydrant,omitempty"`
	LaunchDarkly   CatalogEntityLaunchDarkly    `json:"x-cortex-launch-darkly,omitempty" yaml:"x-cortex-launch-darkly,omitempty"`
	MicrosoftTeams []CatalogEntityMicrosoftTeam `json:"x-cortex-microsoft-teams,omitempty" yaml:"x-cortex-microsoft-teams,omitempty"`
	Rollbar        CatalogEntityRollbar         `json:"x-cortex-rollbar,omitempty" yaml:"x-cortex-rollbar,omitempty"`
	Sentry         CatalogEntitySentry          `json:"x-cortex-sentry,omitempty" yaml:"x-cortex-sentry,omitempty"`
	ServiceNow     CatalogEntityServiceNow      `json:"x-cortex-servicenow,omitempty" yaml:"x-cortex-servicenow,omitempty"`
	Slack          CatalogEntitySlack           `json:"x-cortex-slack,omitempty" yaml:"x-cortex-slack,omitempty"`
	Snyk           CatalogEntitySnyk            `json:"x-cortex-snyk,omitempty" yaml:"x-cortex-snyk,omitempty"`
	Wiz            CatalogEntityWiz             `json:"x-cortex-wiz,omitempty" yaml:"x-cortex-wiz,omitempty"`

	// Infrastructure, Resources, and Deployments attributes
	K8s CatalogEntityK8s `json:"x-cortex-k8s,omitempty" yaml:"x-cortex-k8s,omitempty"`

	// Team-specific attributes
	Team CatalogEntityTeam `json:"team" yaml:"x-cortex-team,omitempty"`
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
 * Teams - Generic structs for teams
 **********************************************************************************************************************/

type CatalogEntityTeam struct {
	Members []CatalogEntityTeamMember  `json:"members" yaml:"members,omitempty"`
	Groups  []CatalogEntityGroupMember `json:"groups" yaml:"groups,omitempty"`
}

func (o *CatalogEntityTeam) Enabled() bool {
	return len(o.Members) > 0 || len(o.Groups) > 0
}

type CatalogEntityTeamMember struct {
	Name                 string `json:"name" yaml:"name"`
	Email                string `json:"email" yaml:"email"`
	Role                 string `json:"role,omitempty" yaml:"role,omitempty"`
	NotificationsEnabled bool   `json:"notificationsEnabled" yaml:"notificationsEnabled"`
}

func (o *CatalogEntityTeamMember) Enabled() bool {
	return o.Name != "" && o.Email != ""
}

type CatalogEntityGroupMember struct {
	Name     string `json:"name" yaml:"name"`
	Provider string `json:"provider" yaml:"provider"`
}

func (o *CatalogEntityGroupMember) Enabled() bool {
	return o.Name != "" && o.Provider != ""
}

/***********************************************************************************************************************
 * Integrations - Generic structs for integrations
 **********************************************************************************************************************/

type CatalogEntityAlert struct {
	Type  string `json:"type" yaml:"type"`
	Tag   string `json:"tag" yaml:"tag"`
	Value string `json:"value" yaml:"value"`
}

func (o *CatalogEntityAlert) Enabled() bool {
	return o.Tag != ""
}

type CatalogEntityApm struct {
	DataDog   CatalogEntityApmDataDog    `json:"datadog,omitempty" yaml:"datadog,omitempty"`
	Dynatrace CatalogEntityApmDynatrace  `json:"dynatrace,omitempty" yaml:"dynatrace,omitempty"`
	NewRelic  []CatalogEntityApmNewRelic `json:"newrelic,omitempty" yaml:"newrelic,omitempty"`
}

func (c *CatalogEntityApm) Enabled() bool {
	return c.DataDog.Enabled() || c.Dynatrace.Enabled() || len(c.NewRelic) > 0
}

type CatalogEntityDashboards struct {
	Embeds []CatalogEntityDashboardsEmbed `json:"embeds,omitempty" yaml:"embeds,omitempty"`
}

func (c *CatalogEntityDashboards) Enabled() bool {
	return len(c.Embeds) > 0
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

func (o *CatalogEntityGit) Enabled() bool {
	return o.Github.Enabled() || o.Gitlab.Enabled() || o.Azure.Enabled() || o.BitBucket.Enabled()
}

type CatalogEntityIssues struct {
	Jira CatalogEntityIssuesJira `json:"jira,omitempty" yaml:"jira,omitempty"`
}

func (c *CatalogEntityIssues) Enabled() bool {
	return c.Jira.Enabled()
}

type CatalogEntityOnCall struct {
	PagerDuty CatalogEntityOnCallPagerDuty `json:"pagerduty,omitempty" yaml:"pagerduty,omitempty"`
	OpsGenie  CatalogEntityOnCallOpsGenie  `json:"opsgenie,omitempty" yaml:"opsgenie,omitempty"`
	VictorOps CatalogEntityOnCallVictorOps `json:"victorops,omitempty" yaml:"victorops,omitempty"`
	XMatters  CatalogEntityOnCallXMatters  `json:"xmatters,omitempty" yaml:"xmatters,omitempty"`
}

func (c *CatalogEntityOnCall) Enabled() bool {
	return c.PagerDuty.Enabled() || c.OpsGenie.Enabled() || c.VictorOps.Enabled()
}

type CatalogEntityOwner struct {
	Type                 string `json:"type" yaml:"type"`                       // group, user, slack
	Name                 string `json:"name,omitempty" yaml:"name,omitempty"`   // Must be of form <org>/<team>
	Email                string `json:"email,omitempty" yaml:"email,omitempty"` // user only
	Description          string `json:"description,omitempty" yaml:"description,omitempty"`
	Provider             string `json:"provider,omitempty" yaml:"provider,omitempty"`
	Channel              string `json:"channel,omitempty" yaml:"channel,omitempty"` // for slack, do not add # to beginning
	NotificationsEnabled bool   `json:"notificationsEnabled,omitempty" yaml:"notificationsEnabled,omitempty"`
}

type CatalogEntityChild struct {
	Tag string `json:"tag" yaml:"tag"`
}

type CatalogEntityDomainParent struct {
	Tag string `json:"tag" yaml:"tag"`
}

type CatalogEntitySLOs struct {
	DataDog    []CatalogEntitySLODataDog         `json:"datadog,omitempty" yaml:"datadog,omitempty"`
	Dynatrace  []CatalogEntitySLODynatrace       `json:"dynatrace,omitempty" yaml:"dynatrace,omitempty"`
	Lightstep  []CatalogEntitySLOLightstepStream `json:"lightstep,omitempty" yaml:"lightstep,omitempty"`
	Prometheus []CatalogEntitySLOPrometheusQuery `json:"prometheus,omitempty" yaml:"prometheus,omitempty"`
	SignalFX   []CatalogEntitySLOSignalFX        `json:"signalfx,omitempty" yaml:"signalfx,omitempty"`
	SumoLogic  []CatalogEntitySLOSumoLogic       `json:"sumologic,omitempty" yaml:"sumologic,omitempty"`
}

func (c *CatalogEntitySLOs) Enabled() bool {
	return len(c.DataDog) > 0 || len(c.Dynatrace) > 0 || len(c.Lightstep) > 0 || len(c.Prometheus) > 0 || len(c.SignalFX) > 0 || len(c.SumoLogic) > 0
}

type CatalogEntityStaticAnalysis struct {
	CodeCov   CatalogEntityStaticAnalysisCodeCov   `json:"codecov,omitempty" yaml:"codecov,omitempty"`
	Mend      CatalogEntityStaticAnalysisMend      `json:"mend,omitempty" yaml:"mend,omitempty"`
	SonarQube CatalogEntityStaticAnalysisSonarQube `json:"sonarqube,omitempty" yaml:"sonarqube,omitempty"`
	Veracode  CatalogEntityStaticAnalysisVeracode  `json:"veracode,omitempty" yaml:"veracode,omitempty"`
}

func (c *CatalogEntityStaticAnalysis) Enabled() bool {
	return c.CodeCov.Enabled() || c.Mend.Enabled() || c.SonarQube.Enabled() || c.Veracode.Enabled()
}

/***********************************************************************************************************************
 * CI/CD
 **********************************************************************************************************************/

type CatalogEntityCiCd struct {
	Buildkite CatalogEntityCiCdBuildkite `json:"buildkite,omitempty" yaml:"buildkite,omitempty"`
}

func (c *CatalogEntityCiCd) Enabled() bool {
	return c.Buildkite.Enabled()
}

/***********************************************************************************************************************
 * Azure DevOps - https://docs.cortex.io/docs/reference/integrations/azuredevops
 **********************************************************************************************************************/

type CatalogEntityGitAzureDevOps struct {
	Project    string `json:"project" yaml:"project"`
	Repository string `json:"repository" yaml:"repository"`
	BasePath   string `json:"basepath,omitempty" yaml:"basepath,omitempty"`
}

func (o *CatalogEntityGitAzureDevOps) Enabled() bool {
	return o.Repository != ""
}

/***********************************************************************************************************************
 * BitBucket - https://docs.cortex.io/docs/reference/integrations/bitbucket
 **********************************************************************************************************************/

type CatalogEntityGitBitBucket struct {
	Repository string `json:"repository" yaml:"repository"`
}

func (o *CatalogEntityGitBitBucket) Enabled() bool {
	return o.Repository != ""
}

/***********************************************************************************************************************
 * BugSnag - https://docs.cortex.io/docs/reference/integrations/bugsnag
 **********************************************************************************************************************/

type CatalogEntityBugSnag struct {
	Project string `json:"project" yaml:"project"`
}

func (o *CatalogEntityBugSnag) Enabled() bool {
	return o.Project != ""
}

/***********************************************************************************************************************
 * Buildkite - https://docs.cortex.io/docs/reference/integrations/buildkite
 **********************************************************************************************************************/

type CatalogEntityCiCdBuildkite struct {
	Pipelines []CatalogEntityCiCdBuildkitePipeline `json:"pipelines" yaml:"pipelines"`
	Tags      []CatalogEntityCiCdBuildkiteTag      `json:"tags,omitempty" yaml:"tags,omitempty"`
}

func (o *CatalogEntityCiCdBuildkite) Enabled() bool {
	return len(o.Tags) > 0 || len(o.Pipelines) > 0
}

type CatalogEntityCiCdBuildkitePipeline struct {
	Slug string `json:"slug" yaml:"slug"`
}

func (o *CatalogEntityCiCdBuildkitePipeline) Enabled() bool {
	return o.Slug != ""
}

type CatalogEntityCiCdBuildkiteTag struct {
	Tag string `json:"tag" yaml:"tag"`
}

func (o *CatalogEntityCiCdBuildkiteTag) Enabled() bool {
	return o.Tag != ""
}

/***********************************************************************************************************************
 * Coralogix - https://docs.cortex.io/docs/reference/integrations/coralogix
 **********************************************************************************************************************/

type CatalogEntityCoralogix struct {
	Applications []CatalogEntityCoralogixApplication `json:"applications" yaml:"applications"`
}

func (o *CatalogEntityCoralogix) Enabled() bool {
	return len(o.Applications) > 0
}

type CatalogEntityCoralogixApplication struct {
	Name  string `json:"applicationName,omitempty" yaml:"applicationName"`
	Alias string `json:"alias,omitempty" yaml:"alias,omitempty"`
}

func (o *CatalogEntityCoralogixApplication) Enabled() bool {
	return o.Name != ""
}

/***********************************************************************************************************************
 * Coralogix - https://docs.cortex.io/docs/reference/integrations/coralogix
 **********************************************************************************************************************/

type CatalogEntityCheckmarx struct {
	Projects []CatalogEntityCheckmarxProject `json:"projects" yaml:"projects"`
}

type CatalogEntityCheckmarxProject struct {
	ID   int64  `json:"projectId,omitempty" yaml:"projectId,omitempty"`
	Name string `json:"projectName,omitempty" yaml:"projectName,omitempty"`
}

func (o *CatalogEntityCheckmarx) Enabled() bool {
	return len(o.Projects) > 0
}

/***********************************************************************************************************************
 * CodeCov - https://docs.cortex.io/docs/reference/integrations/codecov
 **********************************************************************************************************************/

type CatalogEntityStaticAnalysisCodeCov struct {
	Repository string `json:"repo" yaml:"repo"`
	Provider   string `json:"provider" yaml:"provider"`
	Owner      string `json:"owner" yaml:"owner"`
	Flag       string `json:"flag" yaml:"flag"`
}

func (o *CatalogEntityStaticAnalysisCodeCov) Enabled() bool {
	return o.Repository != ""
}

/***********************************************************************************************************************
 * DataDog - https://docs.cortex.io/docs/reference/integrations/datadog
 **********************************************************************************************************************/

type CatalogEntityApmDataDog struct {
	Monitors []int64 `json:"monitors" yaml:"monitors"`
}

func (o *CatalogEntityApmDataDog) Enabled() bool {
	return len(o.Monitors) > 0
}

type CatalogEntitySLODataDog struct {
	ID string `json:"id" yaml:"id"`
}

func (o *CatalogEntitySLODataDog) Enabled() bool {
	return o.ID != ""
}

/***********************************************************************************************************************
 * Dynatrace - https://docs.cortex.io/docs/reference/integrations/dynatrace
 **********************************************************************************************************************/

type CatalogEntityApmDynatrace struct {
	EntityIDs          []string `json:"entityIds,omitempty" yaml:"entityIds,omitempty"`
	EntityNameMatchers []string `json:"entityNameMatchers,omitempty" yaml:"entityNameMatchers,omitempty"`
}

func (o *CatalogEntityApmDynatrace) Enabled() bool {
	return len(o.EntityIDs) > 0 || len(o.EntityNameMatchers) > 0
}

type CatalogEntitySLODynatrace struct {
	ID string `json:"id" yaml:"id"`
}

func (o *CatalogEntitySLODynatrace) Enabled() bool {
	return o.ID != ""
}

/***********************************************************************************************************************
 * FireHydrant - https://docs.cortex.io/docs/reference/integrations/firehydrant
 **********************************************************************************************************************/

type CatalogEntityFireHydrant struct {
	Services []CatalogEntityFireHydrantService `json:"services" yaml:"services"`
}

func (o *CatalogEntityFireHydrant) Enabled() bool {
	return len(o.Services) > 0
}

type CatalogEntityFireHydrantService struct {
	ID   string `json:"identifier" yaml:"identifier"`
	Type string `json:"identifierType" yaml:"identifierType"`
}

func (o *CatalogEntityFireHydrantService) Enabled() bool {
	return o.ID != "" && o.Type != ""
}

/***********************************************************************************************************************
 * GitHub - https://docs.cortex.io/docs/reference/integrations/github
 **********************************************************************************************************************/

type CatalogEntityGitGithub struct {
	Repository string `json:"repository" yaml:"repository"`
	BasePath   string `json:"basepath,omitempty" yaml:"basepath,omitempty"`
}

func (o *CatalogEntityGitGithub) Enabled() bool {
	return o.Repository != ""
}

/***********************************************************************************************************************
 * GitLab - https://docs.cortex.io/docs/reference/integrations/gitlab
 **********************************************************************************************************************/

type CatalogEntityGitGitlab struct {
	Repository string `json:"repository" yaml:"repository"`
	BasePath   string `json:"basepath,omitempty" yaml:"basepath,omitempty"`
}

func (o *CatalogEntityGitGitlab) Enabled() bool {
	return o.Repository != ""
}

/***********************************************************************************************************************
 * JIRA - https://docs.cortex.io/docs/reference/integrations/jira
 **********************************************************************************************************************/

type CatalogEntityIssuesJira struct {
	DefaultJQL string   `json:"defaultJql,omitempty" yaml:"defaultJql,omitempty"`
	Projects   []string `json:"projects,omitempty" yaml:"projects,omitempty"`
	Labels     []string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Components []string `json:"components,omitempty" yaml:"components,omitempty"`
}

func (o *CatalogEntityIssuesJira) Enabled() bool {
	return len(o.Projects) > 0 || len(o.Labels) > 0 || len(o.Components) > 0 || o.DefaultJQL != ""
}

/***********************************************************************************************************************
 * Kubernetes - https://docs.cortex.io/docs/reference/integrations/kubernetes
 **********************************************************************************************************************/

type CatalogEntityK8s struct {
	Deployments  []CatalogEntityK8sDeployment  `json:"deployment,omitempty" yaml:"deployment,omitempty"`
	ArgoRollouts []CatalogEntityK8sArgoRollout `json:"argorollout,omitempty" yaml:"argorollout,omitempty"`
	StatefulSets []CatalogEntityK8sStatefulSet `json:"statefulset,omitempty" yaml:"statefulset,omitempty"`
	CronJobs     []CatalogEntityK8sCronJob     `json:"cronjob,omitempty" yaml:"cronjob,omitempty"`
}

func (o *CatalogEntityK8s) Enabled() bool {
	return len(o.Deployments) > 0 || len(o.ArgoRollouts) > 0 || len(o.StatefulSets) > 0 || len(o.CronJobs) > 0
}

type CatalogEntityK8sDeployment struct {
	Identifier string `json:"identifier" yaml:"identifier"`
	Cluster    string `json:"cluster,omitempty" yaml:"cluster,omitempty"`
}

func (o *CatalogEntityK8sDeployment) Enabled() bool {
	return o.Identifier != ""
}

type CatalogEntityK8sArgoRollout struct {
	Identifier string `json:"identifier" yaml:"identifier"`
	Cluster    string `json:"cluster,omitempty" yaml:"cluster,omitempty"`
}

func (o *CatalogEntityK8sArgoRollout) Enabled() bool {
	return o.Identifier != ""
}

type CatalogEntityK8sStatefulSet struct {
	Identifier string `json:"identifier" yaml:"identifier"`
	Cluster    string `json:"cluster,omitempty" yaml:"cluster,omitempty"`
}

func (o *CatalogEntityK8sStatefulSet) Enabled() bool {
	return o.Identifier != ""
}

type CatalogEntityK8sCronJob struct {
	Identifier string `json:"identifier" yaml:"identifier"`
	Cluster    string `json:"cluster,omitempty" yaml:"cluster,omitempty"`
}

func (o *CatalogEntityK8sCronJob) Enabled() bool {
	return o.Identifier != ""
}

/***********************************************************************************************************************
 * LaunchDarkly - https://docs.cortex.io/docs/reference/integrations/launchdarkly
 **********************************************************************************************************************/

type CatalogEntityLaunchDarkly struct {
	Projects []CatalogEntityLaunchDarklyProject `json:"projects,omitempty" yaml:"projects,omitempty"`
}

func (o *CatalogEntityLaunchDarkly) Enabled() bool {
	return len(o.Projects) > 0
}

type CatalogEntityLaunchDarklyProject struct {
	ID           string                                        `json:"identifier" yaml:"identifier"`
	Type         string                                        `json:"identifierType" yaml:"identifierType"`
	Alias        string                                        `json:"alias,omitempty" yaml:"alias,omitempty"`
	Environments []CatalogEntityLaunchDarklyProjectEnvironment `json:"environments,omitempty" yaml:"environments,omitempty"`
}

func (o *CatalogEntityLaunchDarklyProject) Enabled() bool {
	return o.ID != "" && o.Type != ""
}

type CatalogEntityLaunchDarklyProjectEnvironment struct {
	Name string `json:"environmentName" yaml:"environmentName"`
}

func (o *CatalogEntityLaunchDarklyProjectEnvironment) Enabled() bool {
	return o.Name != ""
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

type CatalogEntitySLOLightstepStream struct {
	StreamID string                           `json:"streamId" yaml:"streamId"`
	Targets  CatalogEntitySLOLightstepTargets `json:"targets" yaml:"targets"`
}

type CatalogEntitySLOLightstepTargets struct {
	Latencies []CatalogEntitySLOLightstepTargetLatency `json:"latency" yaml:"latency"`
}

type CatalogEntitySLOLightstepTargetLatency struct {
	Percentile float64 `json:"percentile" yaml:"percentile"`
	Target     int64   `json:"target" yaml:"target"`
	SLO        float64 `json:"slo" yaml:"slo"`
}

/***********************************************************************************************************************
 * Mend - https://docs.cortex.io/docs/reference/integrations/mend
 **********************************************************************************************************************/

type CatalogEntityStaticAnalysisMend struct {
	ApplicationIDs []string `json:"applicationIds,omitempty" yaml:"applicationIds,omitempty"`
	ProjectIDs     []string `json:"projectIds,omitempty" yaml:"projectIds,omitempty"`
}

func (o *CatalogEntityStaticAnalysisMend) Enabled() bool {
	return len(o.ApplicationIDs) > 0 || len(o.ProjectIDs) > 0
}

/***********************************************************************************************************************
 * Microsoft Teams - https://docs.cortex.io/docs/reference/integrations/microsoftteams
 **********************************************************************************************************************/

/**
x-cortex-microsoft-teams:
    - name: team-engineering # exact match name of the channel
      description: This is a description for this Teams channel # optional
      notificationsEnabled: true #optional
*/

type CatalogEntityMicrosoftTeam struct {
	Name                 string `json:"name" yaml:"name"`
	Description          string `json:"description,omitempty" yaml:"description,omitempty"`
	NotificationsEnabled bool   `json:"notificationsEnabled,omitempty" yaml:"notificationsEnabled,omitempty"`
}

func (o *CatalogEntityMicrosoftTeam) Enabled() bool {
	return o.Name != ""
}

/***********************************************************************************************************************
 * New Relic - https://docs.cortex.io/docs/reference/integrations/newrelic
 **********************************************************************************************************************/

type CatalogEntityApmNewRelic struct {
	ApplicationID int64  `json:"applicationId" yaml:"applicationId"`
	Alias         string `json:"alias,omitempty" yaml:"alias,omitempty"`
}

func (o *CatalogEntityApmNewRelic) Enabled() bool {
	return o.ApplicationID > 0 || o.Alias != ""
}

/***********************************************************************************************************************
 * PagerDuty - https://docs.cortex.io/docs/reference/integrations/pagerduty
 **********************************************************************************************************************/

type CatalogEntityOnCallPagerDuty struct {
	ID   string `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}

func (o *CatalogEntityOnCallPagerDuty) Enabled() bool {
	return o.ID != ""
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

func (o *CatalogEntityRollbar) Enabled() bool {
	return o.Project != ""
}

/***********************************************************************************************************************
 * SignalFX - https://docs.cortex.io/docs/reference/integrations/signalfx
 **********************************************************************************************************************/

type CatalogEntitySLOSignalFX struct {
	Query     string `json:"query" yaml:"query"`
	Rollup    string `json:"rollup" yaml:"rollup"`
	Target    int64  `json:"target" yaml:"target"`
	Lookback  string `json:"lookback" yaml:"lookback"`
	Operation string `json:"operation" yaml:"operation"`
}

/***********************************************************************************************************************
 * Sentry - https://docs.cortex.io/docs/reference/integrations/sentry
 **********************************************************************************************************************/

type CatalogEntitySentry struct {
	Project  string                       `json:"project,omitempty" yaml:"project,omitempty"`
	Projects []CatalogEntitySentryProject `json:"projects,omitempty" yaml:"projects,omitempty"`
}

func (o *CatalogEntitySentry) Enabled() bool {
	return o.Project != "" || len(o.Projects) > 0
}

type CatalogEntitySentryProject struct {
	Name string `json:"name" yaml:"name"`
}

func (o *CatalogEntitySentryProject) Enabled() bool {
	return o.Name != ""
}

/***********************************************************************************************************************
 * ServiceNow - https://docs.cortex.io/docs/reference/integrations/servicenow
 **********************************************************************************************************************/

type CatalogEntityServiceNow struct {
	Services []CatalogEntityServiceNowService `json:"services" yaml:"services"`
}

func (o *CatalogEntityServiceNow) Enabled() bool {
	return len(o.Services) > 0
}

type CatalogEntityServiceNowService struct {
	ID        int64  `json:"id" yaml:"id"`
	TableName string `json:"tableName" yaml:"tableName"`
}

func (o *CatalogEntityServiceNowService) Enabled() bool {
	return o.ID > 0 || o.TableName != ""
}

/***********************************************************************************************************************
 * Slack - https://docs.cortex.io/docs/reference/integrations/sentry
 **********************************************************************************************************************/

type CatalogEntitySlack struct {
	Channels []CatalogEntitySlackChannel `json:"channels,omitempty" yaml:"channels,omitempty"`
}

func (o *CatalogEntitySlack) Enabled() bool {
	return len(o.Channels) > 0
}

type CatalogEntitySlackChannel struct {
	Name                 string `json:"name" yaml:"name"`
	NotificationsEnabled bool   `json:"notificationsEnabled,omitempty" yaml:"notificationsEnabled,omitempty"`
}

func (o *CatalogEntitySlackChannel) Enabled() bool {
	return o.Name != ""
}

/***********************************************************************************************************************
 * Snyk - https://docs.cortex.io/docs/reference/integrations/snyk
 **********************************************************************************************************************/

type CatalogEntitySnyk struct {
	Projects []CatalogEntitySnykProject `json:"projects,omitempty" yaml:"projects,omitempty"`
}

func (o *CatalogEntitySnyk) Enabled() bool {
	return len(o.Projects) > 0
}

type CatalogEntitySnykProject struct {
	Organization string `json:"organizationId" yaml:"organizationId"`
	ProjectID    string `json:"projectId" yaml:"projectId"`
	Source       string `json:"source,omitempty" yaml:"source,omitempty"`
}

/***********************************************************************************************************************
 * SonarQube - https://docs.cortex.io/docs/reference/integrations/sonarqube
 **********************************************************************************************************************/

type CatalogEntityStaticAnalysisSonarQube struct {
	Project string `json:"project" yaml:"project"`
	Alias   string `json:"alias,omitempty" yaml:"alias,omitempty"`
}

func (o *CatalogEntityStaticAnalysisSonarQube) Enabled() bool {
	return o.Project != ""
}

/***********************************************************************************************************************
 * SumoLogic - https://docs.cortex.io/docs/reference/integrations/sumologic
 **********************************************************************************************************************/

type CatalogEntitySLOSumoLogic struct {
	ID string `json:"id" yaml:"id"`
}

func (o *CatalogEntitySLOSumoLogic) Enabled() bool {
	return o.ID != ""
}

/***********************************************************************************************************************
 * OpsGenie -  https://docs.cortex.io/docs/reference/integrations/opsgenie
 **********************************************************************************************************************/

type CatalogEntityOnCallOpsGenie struct {
	ID   string `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}

func (o *CatalogEntityOnCallOpsGenie) Enabled() bool {
	return o.ID != ""
}

/***********************************************************************************************************************
 * Veracode -  https://docs.cortex.io/docs/reference/integrations/veracode
 **********************************************************************************************************************/

type CatalogEntityStaticAnalysisVeracode struct {
	ApplicationNames []string                                     `json:"applicationNames,omitempty" yaml:"applicationNames,omitempty"`
	Sandboxes        []CatalogEntityStaticAnalysisVeracodeSandbox `json:"sandboxes,omitempty" yaml:"sandboxes,omitempty"`
}

func (o *CatalogEntityStaticAnalysisVeracode) Enabled() bool {
	return len(o.ApplicationNames) > 0 || len(o.Sandboxes) > 0
}

type CatalogEntityStaticAnalysisVeracodeSandbox struct {
	ApplicationName string `json:"applicationName,omitempty" yaml:"applicationName,omitempty"`
	SandboxName     string `json:"sandboxName,omitempty" yaml:"sandboxName,omitempty"`
}

/***********************************************************************************************************************
 * VictorOps -  https://docs.cortex.io/docs/reference/integrations/victorops
 **********************************************************************************************************************/

type CatalogEntityOnCallVictorOps struct {
	ID   string `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}

func (o *CatalogEntityOnCallVictorOps) Enabled() bool {
	return o.ID != ""
}

/***********************************************************************************************************************
 * Wiz - https://docs.cortex.io/docs/reference/integrations/wiz
 **********************************************************************************************************************/

type CatalogEntityWiz struct {
	Projects []CatalogEntityWizProject `json:"projects,omitempty" yaml:"projects,omitempty"`
}

func (o *CatalogEntityWiz) Enabled() bool {
	return len(o.Projects) > 0
}

type CatalogEntityWizProject struct {
	ProjectID string `json:"projectId" yaml:"projectId"`
}

/***********************************************************************************************************************
 * XMatters -  https://docs.cortex.io/docs/reference/integrations/xmatters
 **********************************************************************************************************************/

type CatalogEntityOnCallXMatters struct {
	ID   string `json:"id" yaml:"id"`
	Type string `json:"type" yaml:"type"`
}

func (o *CatalogEntityOnCallXMatters) Enabled() bool {
	return o.ID != "" && o.Type != ""
}
