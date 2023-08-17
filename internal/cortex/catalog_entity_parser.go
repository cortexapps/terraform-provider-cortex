package cortex

import (
	"fmt"
	"reflect"
)

type CatalogEntityParser struct{}

// YamlToEntity converts YAML into a CatalogEntity, from the specification.
func (c *CatalogEntityParser) YamlToEntity(yamlEntity map[string]interface{}) (CatalogEntityData, error) {
	entity := CatalogEntityData{}
	info := yamlEntity["info"].(map[string]interface{})

	entity.Title = MapFetchToString(info, "title")
	entity.Description = MapFetchToString(info, "description")
	entity.Tag = MapFetchToString(info, "x-cortex-tag")
	entity.Type = MapFetch(info, "x-cortex-type", "service").(string)

	entity.Definition = map[string]interface{}{}
	if info["x-cortex-definition"] != nil {
		entity.Definition = info["x-cortex-definition"].(map[string]interface{})
	}

	entity.Links = []CatalogEntityLink{}
	if info["x-cortex-link"] != nil {
		c.interpolateLinks(&entity, info["x-cortex-link"].([]interface{}))
	}

	entity.Groups = []string{}
	if info["x-cortex-groups"] != nil {
		for _, group := range info["x-cortex-groups"].([]interface{}) {
			entity.Groups = append(entity.Groups, group.(string))
		}
	}

	entity.Owners = []CatalogEntityOwner{}
	if info["x-cortex-owners"] != nil {
		c.interpolateOwners(&entity, info["x-cortex-owners"].([]interface{}))
	}

	entity.Children = []CatalogEntityChild{}
	if info["x-cortex-children"] != nil {
		c.interpolateChildren(&entity, info["x-cortex-children"].([]interface{}))
	}

	entity.DomainParents = []CatalogEntityDomainParent{}
	if info["x-cortex-domain-parents"] != nil {
		c.interpolateDomainParents(&entity, info["x-cortex-domain-parents"].([]interface{}))
	}

	entity.Metadata = map[string]interface{}{}
	if info["x-cortex-custom-metadata"] != nil {
		entity.Metadata = info["x-cortex-custom-metadata"].(map[string]interface{})
	}

	entity.Dependencies = []CatalogEntityDependency{}
	if info["x-cortex-dependency"] != nil {
		c.interpolateDependencies(&entity, info["x-cortex-dependency"].([]interface{}))
	}

	if info["x-cortex-git"] != nil {
		entity.Git = CatalogEntityGit{}
		c.interpolateGit(&entity, info["x-cortex-git"].(map[string]interface{}))
	}

	if info["x-cortex-dashboards"] != nil {
		c.interpolateDashboards(&entity, info["x-cortex-dashboards"].(map[string]interface{}))
	}

	if info["x-cortex-issues"] != nil {
		c.interpolateIssues(&entity, info["x-cortex-issues"].(map[string]interface{}))
	}

	if info["x-cortex-slos"] != nil {
		slosMap := info["x-cortex-slos"].(map[string]interface{})
		err := c.interpolateSLOs(&entity, slosMap)
		if err != nil {
			return entity, err
		}
	}

	if info["x-cortex-apm"] != nil {
		c.interpolateApm(&entity, info["x-cortex-apm"].(map[string]interface{}))
	}

	if info["x-cortex-static-analysis"] != nil {
		c.interpolateStaticAnalysis(&entity, info["x-cortex-static-analysis"].(map[string]interface{}))
	}

	if info["x-cortex-ci-cd"] != nil {
		c.interpolateCiCd(&entity, info["x-cortex-ci-cd"].(map[string]interface{}))
	}

	if info["x-cortex-oncall"] != nil {
		onCallMap := info["x-cortex-oncall"].(map[string]interface{})
		c.interpolateOnCall(&entity, onCallMap)
	}

	if info["x-cortex-alerts"] != nil {
		c.interpolateAlerts(&entity, info["x-cortex-alerts"].([]interface{}))
	}

	if info["x-cortex-bugsnag"] != nil {
		c.interpolateBugSnag(&entity, info["x-cortex-bugsnag"].(map[string]interface{}))
	}

	if info["x-cortex-checkmarx"] != nil {
		c.interpolateCheckmarx(&entity, info["x-cortex-checkmarx"].(map[string]interface{}))
	}

	if info["x-cortex-firehydrant"] != nil {
		c.interpolateFirehydrant(&entity, info["x-cortex-firehydrant"].(map[string]interface{}))
	}

	if info["x-cortex-rollbar"] != nil {
		c.interpolateRollbar(&entity, info["x-cortex-rollbar"].(map[string]interface{}))
	}

	if info["x-cortex-sentry"] != nil {
		c.interpolateSentry(&entity, info["x-cortex-sentry"].(map[string]interface{}))
	}

	if info["x-cortex-slack"] != nil {
		c.interpolateSlack(&entity, info["x-cortex-slack"].(map[string]interface{}))
	}

	if info["x-cortex-snyk"] != nil {
		c.interpolateSnyk(&entity, info["x-cortex-snyk"].(map[string]interface{}))
	}

	if info["x-cortex-wiz"] != nil {
		c.interpolateWiz(&entity, info["x-cortex-wiz"].(map[string]interface{}))
	}

	// team-specific entity attributes
	if info["x-cortex-team"] != nil {
		c.interpolateTeam(&entity, info["x-cortex-team"].(map[string]interface{}))
	}

	return entity, nil
}

func (c *CatalogEntityParser) interpolateLinks(entity *CatalogEntityData, links []interface{}) {
	for _, link := range links {
		linkMap := link.(map[string]interface{})
		entity.Links = append(entity.Links, CatalogEntityLink{
			Name: MapFetchToString(linkMap, "name"),
			Type: MapFetchToString(linkMap, "type"),
			Url:  MapFetchToString(linkMap, "url"),
		})
	}
}

func (c *CatalogEntityParser) interpolateOwners(entity *CatalogEntityData, owners []interface{}) {
	for _, owner := range owners {
		ownerMap := owner.(map[string]interface{})
		entity.Owners = append(entity.Owners, CatalogEntityOwner{
			Type:                 MapFetchToString(ownerMap, "type"),
			Name:                 MapFetchToString(ownerMap, "name"),
			Email:                MapFetchToString(ownerMap, "email"),
			Description:          MapFetchToString(ownerMap, "description"),
			Provider:             MapFetchToString(ownerMap, "provider"),
			Channel:              MapFetchToString(ownerMap, "channel"),
			NotificationsEnabled: MapFetch(ownerMap, "notificationsEnabled", false).(bool),
		})
	}
}

func (c *CatalogEntityParser) interpolateChildren(entity *CatalogEntityData, children []interface{}) {
	for _, child := range children {
		childMap := child.(map[string]interface{})
		entity.Children = append(entity.Children, CatalogEntityChild{
			Tag: MapFetchToString(childMap, "tag"),
		})
	}
}

func (c *CatalogEntityParser) interpolateDomainParents(entity *CatalogEntityData, parents []interface{}) {
	for _, parent := range parents {
		parentMap := parent.(map[string]interface{})
		entity.DomainParents = append(entity.DomainParents, CatalogEntityDomainParent{
			Tag: MapFetchToString(parentMap, "tag"),
		})
	}
}

func (c *CatalogEntityParser) interpolateDependencies(entity *CatalogEntityData, dependencies []interface{}) {
	for _, dependency := range dependencies {
		dependencyMap := dependency.(map[string]interface{})
		entity.Dependencies = append(entity.Dependencies, CatalogEntityDependency{
			Tag:         MapFetchToString(dependencyMap, "tag"),
			Method:      MapFetchToString(dependencyMap, "method"),
			Path:        MapFetchToString(dependencyMap, "path"),
			Description: MapFetchToString(dependencyMap, "description"),
			Metadata:    MapFetch(dependencyMap, "metadata", map[string]interface{}{}).(map[string]interface{}),
		})
	}
}

func (c *CatalogEntityParser) interpolateDashboards(entity *CatalogEntityData, dashboardsMap map[string]interface{}) {
	entity.Dashboards = CatalogEntityDashboards{
		Embeds: []CatalogEntityDashboardsEmbed{},
	}
	if dashboardsMap["embeds"] != nil {
		embeds := dashboardsMap["embeds"].([]interface{})
		for _, embed := range embeds {
			embedMap := embed.(map[string]interface{})
			entity.Dashboards.Embeds = append(entity.Dashboards.Embeds, CatalogEntityDashboardsEmbed{
				Type: MapFetchToString(embedMap, "type"),
				URL:  MapFetchToString(embedMap, "url"),
			})
		}
	}
}

// OnCall

func (c *CatalogEntityParser) interpolateOnCall(entity *CatalogEntityData, onCallMap map[string]interface{}) {
	entity.OnCall = CatalogEntityOnCall{}
	if onCallMap["pagerduty"] != nil {
		pdMap := onCallMap["pagerduty"].(map[string]interface{})
		entity.OnCall.PagerDuty = CatalogEntityOnCallPagerDuty{
			ID:   MapFetchToString(pdMap, "id"),
			Type: MapFetchToString(pdMap, "type"),
		}
	}
	if onCallMap["opsgenie"] != nil {
		ogMap := onCallMap["opsgenie"].(map[string]interface{})
		entity.OnCall.OpsGenie = CatalogEntityOnCallOpsGenie{
			ID:   MapFetchToString(ogMap, "id"),
			Type: MapFetchToString(ogMap, "type"),
		}
	}
	if onCallMap["victorops"] != nil {
		voMap := onCallMap["victorops"].(map[string]interface{})
		entity.OnCall.VictorOps = CatalogEntityOnCallVictorOps{
			ID:   MapFetchToString(voMap, "id"),
			Type: MapFetchToString(voMap, "type"),
		}
	}
	if onCallMap["xmatters"] != nil {
		voMap := onCallMap["xmatters"].(map[string]interface{})
		entity.OnCall.XMatters = CatalogEntityOnCallXMatters{
			ID:   MapFetchToString(voMap, "id"),
			Type: MapFetchToString(voMap, "type"),
		}
	}
}

// Git

func (c *CatalogEntityParser) interpolateGit(entity *CatalogEntityData, gitMap map[string]interface{}) {
	if gitMap["github"] != nil {
		githubMap := gitMap["github"].(map[string]interface{})
		entity.Git.Github = CatalogEntityGitGithub{
			Repository: MapFetchToString(githubMap, "repository"),
			BasePath:   MapFetchToString(githubMap, "basePath"),
		}
	} else {
		entity.Git.Github = CatalogEntityGitGithub{}
	}
	if gitMap["gitlab"] != nil {
		gitlabMap := gitMap["gitlab"].(map[string]interface{})
		entity.Git.Gitlab = CatalogEntityGitGitlab{
			Repository: MapFetchToString(gitlabMap, "repository"),
			BasePath:   MapFetchToString(gitlabMap, "basePath"),
		}
	} else {
		entity.Git.Gitlab = CatalogEntityGitGitlab{}
	}
	if gitMap["azure"] != nil {
		azureMap := gitMap["azure"].(map[string]interface{})
		entity.Git.Azure = CatalogEntityGitAzureDevOps{
			Project:    MapFetchToString(azureMap, "project"),
			Repository: MapFetchToString(azureMap, "repository"),
			BasePath:   MapFetchToString(azureMap, "basePath"),
		}
	} else {
		entity.Git.Azure = CatalogEntityGitAzureDevOps{}
	}
	if gitMap["bitbucket"] != nil {
		bitbucketMap := gitMap["bitbucket"].(map[string]interface{})
		entity.Git.BitBucket = CatalogEntityGitBitBucket{
			Repository: MapFetchToString(bitbucketMap, "repository"),
		}
	} else {
		entity.Git.BitBucket = CatalogEntityGitBitBucket{}
	}
}

// Issues

func (c *CatalogEntityParser) interpolateIssues(entity *CatalogEntityData, issuesMap map[string]interface{}) {
	entity.Issues = CatalogEntityIssues{}
	if issuesMap["jira"] != nil {
		c.interpolateJira(entity, issuesMap["jira"].(map[string]interface{}))
	}
}

func (c *CatalogEntityParser) interpolateJira(entity *CatalogEntityData, jiraMap map[string]interface{}) {
	if jiraMap["defaultJql"] != nil {
		jql := MapFetchToString(jiraMap, "defaultJql")
		if jql != "" {
			entity.Issues.Jira.DefaultJQL = jql
		}
	}
	if jiraMap["projects"] != nil {
		projects := jiraMap["projects"].([]interface{})
		for _, project := range projects {
			entity.Issues.Jira.Projects = append(entity.Issues.Jira.Projects, project.(string))
		}
	}
	if jiraMap["labels"] != nil {
		labels := jiraMap["labels"].([]interface{})
		for _, label := range labels {
			entity.Issues.Jira.Labels = append(entity.Issues.Jira.Labels, label.(string))
		}
	}
	if jiraMap["components"] != nil {
		components := jiraMap["components"].([]interface{})
		for _, component := range components {
			entity.Issues.Jira.Components = append(entity.Issues.Jira.Components, component.(string))
		}
	}
}

func (c *CatalogEntityParser) interpolateSLOs(entity *CatalogEntityData, slosMap map[string]interface{}) error {
	entity.SLOs = CatalogEntitySLOs{}
	if slosMap["datadog"] != nil {
		c.interpolateDataDogSLOs(entity, slosMap["datadog"].([]interface{}))
	}
	if slosMap["dynatrace"] != nil {
		c.interpolateDynatraceSLOs(entity, slosMap["dynatrace"].([]interface{}))
	}
	if slosMap["lightstep"] != nil {
		c.interpolateLightstepSLOs(entity, slosMap["lightstep"].([]interface{}))
	}
	if slosMap["prometheus"] != nil {
		err := c.interpolatePrometheusSLOs(entity, slosMap["prometheus"].([]interface{}))
		if err != nil {
			return err
		}
	}
	if slosMap["signalfx"] != nil {
		c.interpolateSignalFXSLOs(entity, slosMap["signalfx"].([]interface{}))
	}
	if slosMap["sumologic"] != nil {
		c.interpolateSumoLogicSLOs(entity, slosMap["sumologic"].([]interface{}))
	}
	return nil
}

func (c *CatalogEntityParser) interpolateLightstepSLOs(entity *CatalogEntityData, streams []interface{}) {
	if len(streams) == 0 {
		return
	}

	entity.SLOs.Lightstep = make([]CatalogEntitySLOLightstepStream, len(streams))
	for i, stream := range streams {
		streamMap := stream.(map[string]interface{})
		streamSLO := CatalogEntitySLOLightstepStream{
			StreamID: MapFetchToString(streamMap, "streamId"),
			Targets:  CatalogEntitySLOLightstepTargets{},
		}
		if streamMap["targets"] != nil {
			streamTargetMap := streamMap["targets"].(map[string]interface{})
			if streamTargetMap["latency"] != nil {
				latencies := streamTargetMap["latency"].([]interface{})
				for _, latency := range latencies {
					latencyMap := latency.(map[string]interface{})
					streamSLO.Targets.Latencies = append(streamSLO.Targets.Latencies, CatalogEntitySLOLightstepTargetLatency{
						Percentile: MapFetch(latencyMap, "percentile", 0.0).(float64),
						Target:     int64(MapFetch(latencyMap, "target", 0).(int)),
						SLO:        MapFetch(latencyMap, "slo", 0.0).(float64),
					})
				}
			}
		}
		entity.SLOs.Lightstep[i] = streamSLO
	}
}

func (c *CatalogEntityParser) interpolateDataDogSLOs(entity *CatalogEntityData, slos []interface{}) {
	entity.SLOs.DataDog = []CatalogEntitySLODataDog{}
	for _, slo := range slos {
		sloMap := slo.(map[string]interface{})
		entity.SLOs.DataDog = append(entity.SLOs.DataDog, CatalogEntitySLODataDog{
			ID: MapFetchToString(sloMap, "id"),
		})
	}
}

func (c *CatalogEntityParser) interpolatePrometheusSLOs(entity *CatalogEntityData, prometheusQueries []interface{}) error {
	entity.SLOs.Prometheus = []CatalogEntitySLOPrometheusQuery{}
	for _, query := range prometheusQueries {
		queryMap := query.(map[string]interface{})
		sloVal := MapFetch(queryMap, "slo", 0.0)
		sloValFloat64, err := AnyToFloat64(sloVal)
		if err != nil {
			return fmt.Errorf("error converting SLO value to float64: %+v - %+v - %+v", err, sloVal, reflect.TypeOf(sloVal))
		}
		entity.SLOs.Prometheus = append(entity.SLOs.Prometheus, CatalogEntitySLOPrometheusQuery{
			ErrorQuery: MapFetchToString(queryMap, "errorQuery"),
			TotalQuery: MapFetchToString(queryMap, "totalQuery"),
			Name:       MapFetchToString(queryMap, "name"),
			Alias:      MapFetchToString(queryMap, "alias"),
			SLO:        sloValFloat64,
		})
	}
	return nil
}

func (c *CatalogEntityParser) interpolateSignalFXSLOs(entity *CatalogEntityData, signalFxSLOs []interface{}) {
	entity.SLOs.SignalFX = []CatalogEntitySLOSignalFX{}
	for _, slo := range signalFxSLOs {
		sloMap := slo.(map[string]interface{})
		entity.SLOs.SignalFX = append(entity.SLOs.SignalFX, CatalogEntitySLOSignalFX{
			Query:     MapFetchToString(sloMap, "query"),
			Rollup:    MapFetchToString(sloMap, "rollup"),
			Target:    int64(MapFetch(sloMap, "target", 0).(int)),
			Lookback:  MapFetchToString(sloMap, "lookback"),
			Operation: MapFetchToString(sloMap, "operation"),
		})
	}
}

func (c *CatalogEntityParser) interpolateDynatraceSLOs(entity *CatalogEntityData, slos []interface{}) {
	entity.SLOs.Dynatrace = []CatalogEntitySLODynatrace{}
	for _, slo := range slos {
		sloMap := slo.(map[string]interface{})
		entity.SLOs.Dynatrace = append(entity.SLOs.Dynatrace, CatalogEntitySLODynatrace{
			ID: MapFetchToString(sloMap, "id"),
		})
	}
}

func (c *CatalogEntityParser) interpolateSumoLogicSLOs(entity *CatalogEntityData, slos []interface{}) {
	entity.SLOs.SumoLogic = []CatalogEntitySLOSumoLogic{}
	for _, slo := range slos {
		sloMap := slo.(map[string]interface{})
		entity.SLOs.SumoLogic = append(entity.SLOs.SumoLogic, CatalogEntitySLOSumoLogic{
			ID: MapFetchToString(sloMap, "id"),
		})
	}
}

// APM

func (c *CatalogEntityParser) interpolateApm(entity *CatalogEntityData, apm map[string]interface{}) {
	entity.Apm = CatalogEntityApm{}

	if apm["datadog"] != nil {
		c.interpolateDataDogApm(entity, apm["datadog"].(map[string]interface{}))
	}
	if apm["dynatrace"] != nil {
		c.interpolateDynatraceApm(entity, apm["dynatrace"].(map[string]interface{}))
	}
	if apm["newrelic"] != nil {
		c.interpolateNewRelicApm(entity, apm["newrelic"].(map[string]interface{}))
	}
}

func (c *CatalogEntityParser) interpolateDataDogApm(entity *CatalogEntityData, apm map[string]interface{}) {
	entity.Apm.DataDog = CatalogEntityApmDataDog{}
	if apm["monitors"] != nil {
		entity.Apm.DataDog.Monitors = make([]int64, len(apm["monitors"].([]interface{})))
		for i, monitor := range apm["monitors"].([]interface{}) {
			entity.Apm.DataDog.Monitors[i] = int64(monitor.(int))
		}
	}
}

func (c *CatalogEntityParser) interpolateDynatraceApm(entity *CatalogEntityData, apm map[string]interface{}) {
	entity.Apm.Dynatrace = CatalogEntityApmDynatrace{}
	if apm["entityIds"] != nil {
		entity.Apm.Dynatrace.EntityIDs = make([]string, len(apm["entityIds"].([]interface{})))
		for i, group := range apm["entityIds"].([]interface{}) {
			entity.Apm.Dynatrace.EntityIDs[i] = group.(string)
		}
	}
	if apm["entityNameMatchers"] != nil {
		entity.Apm.Dynatrace.EntityNameMatchers = make([]string, len(apm["entityNameMatchers"].([]interface{})))
		for i, group := range apm["entityNameMatchers"].([]interface{}) {
			entity.Apm.Dynatrace.EntityNameMatchers[i] = group.(string)
		}
	}
}

func (c *CatalogEntityParser) interpolateNewRelicApm(entity *CatalogEntityData, apm map[string]interface{}) {
	entity.Apm.NewRelic = CatalogEntityApmNewRelic{}
	if apm["applicationId"] != nil {
		entity.Apm.NewRelic.ApplicationID = int64(apm["applicationId"].(int))
	}
	if apm["alias"] != nil {
		entity.Apm.NewRelic.Alias = apm["alias"].(string)
	}
}

// Integrations

func (c *CatalogEntityParser) interpolateSentry(entity *CatalogEntityData, sentryMap map[string]interface{}) {
	entity.Sentry.Project = MapFetchToString(sentryMap, "project")
}

func (c *CatalogEntityParser) interpolateRollbar(entity *CatalogEntityData, rollbarMap map[string]interface{}) {
	entity.Rollbar.Project = MapFetchToString(rollbarMap, "project")
}

func (c *CatalogEntityParser) interpolateBugSnag(entity *CatalogEntityData, bugSnagMap map[string]interface{}) {
	entity.BugSnag.Project = MapFetchToString(bugSnagMap, "project")
}

func (c *CatalogEntityParser) interpolateCheckmarx(entity *CatalogEntityData, checkmarxMap map[string]interface{}) {
	entity.Checkmarx = CatalogEntityCheckmarx{
		Projects: []CatalogEntityCheckmarxProject{},
	}
	if checkmarxMap["projects"] != nil {
		for _, project := range checkmarxMap["projects"].([]interface{}) {
			projectMap := project.(map[string]interface{})
			pe := CatalogEntityCheckmarxProject{}
			if projectMap["projectId"] != nil {
				pe.ID = MapFetch(projectMap, "projectId", 0).(int64)
			} else if projectMap["projectName"] != nil {
				pe.Name = MapFetchToString(projectMap, "projectName")
			}
			if pe.ID > 0 || pe.Name != "" {
				entity.Checkmarx.Projects = append(entity.Checkmarx.Projects, pe)
			}
		}
	}
}

func (c *CatalogEntityParser) interpolateFirehydrant(entity *CatalogEntityData, firehydrantMap map[string]interface{}) {
	entity.FireHydrant = CatalogEntityFireHydrant{
		Services: []CatalogEntityFireHydrantService{},
	}
	if firehydrantMap["services"] != nil {
		for _, service := range firehydrantMap["services"].([]interface{}) {
			serviceMap := service.(map[string]interface{})
			se := CatalogEntityFireHydrantService{}
			if serviceMap["identifier"] != nil {
				se.ID = MapFetch(serviceMap, "identifier", "").(string)
			}
			if serviceMap["identifierType"] != nil {
				se.Type = MapFetch(serviceMap, "identifierType", "").(string)
			}
			if se.Enabled() {
				entity.FireHydrant.Services = append(entity.FireHydrant.Services, se)
			}
		}
	}
}

func (c *CatalogEntityParser) interpolateSlack(entity *CatalogEntityData, slackMap map[string]interface{}) {
	if slackMap["channels"] != nil {
		channels := slackMap["channels"].([]interface{})
		for _, channel := range channels {
			channelMap := channel.(map[string]interface{})
			entity.Slack.Channels = append(entity.Slack.Channels, CatalogEntitySlackChannel{
				Name:                 MapFetchToString(channelMap, "name"),
				NotificationsEnabled: MapFetch(channelMap, "notificationsEnabled", false).(bool),
			})
		}
	}
}

func (c *CatalogEntityParser) interpolateSnyk(entity *CatalogEntityData, snykMap map[string]interface{}) {
	if snykMap["projects"] != nil {
		projects := snykMap["projects"].([]interface{})
		for _, project := range projects {
			projectMap := project.(map[string]interface{})
			entity.Snyk.Projects = append(entity.Snyk.Projects, CatalogEntitySnykProject{
				ProjectID:    MapFetchToString(projectMap, "projectId"),
				Organization: MapFetchToString(projectMap, "organizationId"),
				Source:       MapFetchToString(projectMap, "source"),
			})
		}
	}
}

func (c *CatalogEntityParser) interpolateWiz(entity *CatalogEntityData, wizMap map[string]interface{}) {
	if wizMap["projects"] != nil {
		projects := wizMap["projects"].([]interface{})
		for _, project := range projects {
			projectMap := project.(map[string]interface{})
			entity.Wiz.Projects = append(entity.Wiz.Projects, CatalogEntityWizProject{
				ProjectID: MapFetchToString(projectMap, "projectId"),
			})
		}
	}
}

func (c *CatalogEntityParser) interpolateAlerts(entity *CatalogEntityData, alerts []interface{}) {
	for _, alert := range alerts {
		alertMap := alert.(map[string]interface{})
		entity.Alerts = append(entity.Alerts, CatalogEntityAlert{
			Type:  MapFetchToString(alertMap, "type"),
			Tag:   MapFetchToString(alertMap, "tag"),
			Value: MapFetchToString(alertMap, "value"),
		})
	}
}

func (c *CatalogEntityParser) interpolateStaticAnalysis(entity *CatalogEntityData, saMap map[string]interface{}) {
	if saMap["codecov"] != nil {
		c.interpolateStaticAnalysisCodeCov(entity, saMap["codecov"].(map[string]interface{}))
	}
	if saMap["mend"] != nil {
		c.interpolateStaticAnalysisMend(entity, saMap["mend"].(map[string]interface{}))
	}
	if saMap["sonarqube"] != nil {
		c.interpolateStaticAnalysisSonarQube(entity, saMap["sonarqube"].(map[string]interface{}))
	}
	if saMap["veracode"] != nil {
		c.interpolateStaticAnalysisVeracode(entity, saMap["veracode"].(map[string]interface{}))
	}
}

func (c *CatalogEntityParser) interpolateStaticAnalysisCodeCov(entity *CatalogEntityData, ccMap map[string]interface{}) {
	entity.StaticAnalysis.CodeCov = CatalogEntityStaticAnalysisCodeCov{
		Repository: MapFetchToString(ccMap, "repo"),
		Provider:   MapFetchToString(ccMap, "provider"),
	}
}

func (c *CatalogEntityParser) interpolateStaticAnalysisMend(entity *CatalogEntityData, data map[string]interface{}) {
	entity.StaticAnalysis.Mend = CatalogEntityStaticAnalysisMend{}
	applicationIds := data["applicationIds"].([]interface{})
	for _, applicationId := range applicationIds {
		if applicationId.(string) != "" {
			entity.StaticAnalysis.Mend.ApplicationIDs = append(entity.StaticAnalysis.Mend.ApplicationIDs, applicationId.(string))
		}
	}
	projectIds := data["projectIds"].([]interface{})
	for _, projectId := range projectIds {
		if projectId.(string) != "" {
			entity.StaticAnalysis.Mend.ProjectIDs = append(entity.StaticAnalysis.Mend.ProjectIDs, projectId.(string))
		}
	}
}

func (c *CatalogEntityParser) interpolateStaticAnalysisSonarQube(entity *CatalogEntityData, data map[string]interface{}) {
	entity.StaticAnalysis.SonarQube = CatalogEntityStaticAnalysisSonarQube{
		Project: data["project"].(string),
	}
}

func (c *CatalogEntityParser) interpolateStaticAnalysisVeracode(entity *CatalogEntityData, mendMap map[string]interface{}) {
	applicationNames := mendMap["applicationNames"].([]interface{})
	if len(applicationNames) == 0 && mendMap["sandboxes"] == nil {
		return
	}

	entity.StaticAnalysis.Veracode = CatalogEntityStaticAnalysisVeracode{}
	for _, applicationName := range applicationNames {
		entity.StaticAnalysis.Veracode.ApplicationNames = append(entity.StaticAnalysis.Veracode.ApplicationNames, applicationName.(string))
	}
	if mendMap["sandboxes"] != nil {
		entity.StaticAnalysis.Veracode.Sandboxes = []CatalogEntityStaticAnalysisVeracodeSandbox{}
		sandboxes := mendMap["sandboxes"].([]interface{})
		for _, sandbox := range sandboxes {
			sandboxMap := sandbox.(map[string]interface{})
			if sandboxMap["applicationName"] != nil || sandboxMap["sandboxName"] != nil {
				sandboxEntity := CatalogEntityStaticAnalysisVeracodeSandbox{}
				if sandboxMap["applicationName"] != nil {
					sandboxEntity.ApplicationName = sandboxMap["applicationName"].(string)
				}
				if sandboxMap["sandboxName"] != nil {
					sandboxEntity.SandboxName = sandboxMap["sandboxName"].(string)
				}
				entity.StaticAnalysis.Veracode.Sandboxes = append(entity.StaticAnalysis.Veracode.Sandboxes, sandboxEntity)
			}
		}
	}
}

func (c *CatalogEntityParser) interpolateCiCd(entity *CatalogEntityData, ciCdMap map[string]interface{}) {
	if ciCdMap["buildkite"] != nil {
		c.interpolateCiCdBuildkite(entity, ciCdMap["buildkite"].(map[string]interface{}))
	}
}

func (c *CatalogEntityParser) interpolateCiCdBuildkite(entity *CatalogEntityData, bkMap map[string]interface{}) {
	entity.CiCd.Buildkite = CatalogEntityCiCdBuildkite{}
	if bkMap["pipelines"] != nil {
		pipelines := bkMap["pipelines"].([]interface{})
		for _, pipeline := range pipelines {
			pipelineMap := pipeline.(map[string]interface{})
			entity.CiCd.Buildkite.Pipelines = append(entity.CiCd.Buildkite.Pipelines, CatalogEntityCiCdBuildkitePipeline{
				Slug: MapFetchToString(pipelineMap, "slug"),
			})
		}
	}
	if bkMap["tags"] != nil {
		tags := bkMap["tags"].([]interface{})
		for _, tag := range tags {
			tagMap := tag.(map[string]interface{})
			entity.CiCd.Buildkite.Tags = append(entity.CiCd.Buildkite.Tags, CatalogEntityCiCdBuildkiteTag{
				Tag: MapFetchToString(tagMap, "tag"),
			})
		}
	}
}

/***********************************************************************************************************************
 * Team attributes
 **********************************************************************************************************************/

func (c *CatalogEntityParser) interpolateTeam(entity *CatalogEntityData, teamMap map[string]interface{}) {
	if teamMap["members"] != nil {
		members := teamMap["members"].([]interface{})
		for _, member := range members {
			memberMap := member.(map[string]interface{})
			entity.Team.Members = append(entity.Team.Members, CatalogEntityTeamMember{
				Name:                 MapFetchToString(memberMap, "name"),
				Email:                MapFetchToString(memberMap, "email"),
				NotificationsEnabled: MapFetch(memberMap, "notificationsEnabled", false).(bool),
			})
		}
	}
	if teamMap["groups"] != nil {
		groups := teamMap["groups"].([]interface{})
		for _, group := range groups {
			groupMap := group.(map[string]interface{})
			entity.Team.Groups = append(entity.Team.Groups, CatalogEntityGroupMember{
				Name:     MapFetchToString(groupMap, "name"),
				Provider: MapFetchToString(groupMap, "provider"),
			})
		}
	}
}
