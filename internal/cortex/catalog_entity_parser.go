package cortex

type CatalogEntityParser struct{}

// YamlToEntity converts YAML into a CatalogEntity, from the specification.
func (c *CatalogEntityParser) YamlToEntity(entity *CatalogEntityData, yamlEntity map[string]interface{}) (*CatalogEntityData, error) {
	info := yamlEntity["info"].(map[string]interface{})

	entity.Title = MapFetchToString(info, "title")
	entity.Description = MapFetchToString(info, "description")
	entity.Tag = MapFetchToString(info, "x-cortex-tag")
	entity.Type = MapFetch(info, "x-cortex-type", "service").(string)

	entity.Links = []CatalogEntityLink{}
	if info["x-cortex-link"] != nil {
		c.interpolateLinks(entity, info["x-cortex-link"].([]interface{}))
	}

	entity.Groups = []string{}
	if info["x-cortex-groups"] != nil {
		for _, group := range info["x-cortex-groups"].([]interface{}) {
			entity.Groups = append(entity.Groups, group.(string))
		}
	}

	entity.Owners = []CatalogEntityOwner{}
	if info["x-cortex-owners"] != nil {
		c.interpolateOwners(entity, info["x-cortex-owners"].([]interface{}))
	}

	entity.Metadata = map[string]interface{}{}
	if info["x-cortex-custom-metadata"] != nil {
		entity.Metadata = info["x-cortex-custom-metadata"].(map[string]interface{})
	}

	entity.Dependencies = []CatalogEntityDependency{}
	if info["x-cortex-dependency"] != nil {
		c.interpolateDependencies(entity, info["x-cortex-dependency"].([]interface{}))
	}

	if info["x-cortex-git"] != nil {
		entity.Git = CatalogEntityGit{}
		c.interpolateGit(entity, info["x-cortex-git"].(map[string]interface{}))
	}

	if info["x-cortex-dashboards"] != nil {
		c.interpolateDashboards(entity, info["x-cortex-dashboards"].(map[string]interface{}))
	}

	if info["x-cortex-issues"] != nil {
		entity.Issues = CatalogEntityIssues{}
		issuesMap := info["x-cortex-issues"].(map[string]interface{})
		if issuesMap["jira"] != nil {
			c.interpolateJira(entity, issuesMap["jira"].(map[string]interface{}))
		}
	}

	if info["x-cortex-slos"] != nil {
		slosMap := info["x-cortex-slos"].(map[string]interface{})
		c.interpolateSLOs(entity, slosMap)
	}

	if info["x-cortex-static-analysis"] != nil {
		c.interpolateStaticAnalysis(entity, info["x-cortex-static-analysis"].(map[string]interface{}))
	}

	if info["x-cortex-oncall"] != nil {
		onCallMap := info["x-cortex-oncall"].(map[string]interface{})
		c.interpolateOnCall(entity, onCallMap)
	}

	if info["x-cortex-alerts"] != nil {
		c.interpolateAlerts(entity, info["x-cortex-alerts"].([]interface{}))
	}

	if info["x-cortex-bugsnag"] != nil {
		c.interpolateBugSnag(entity, info["x-cortex-bugsnag"].(map[string]interface{}))
	}

	if info["x-cortex-checkmarx"] != nil {
		c.interpolateCheckmarx(entity, info["x-cortex-checkmarx"].(map[string]interface{}))
	}

	if info["x-cortex-sentry"] != nil {
		c.interpolateSentry(entity, info["x-cortex-sentry"].(map[string]interface{}))
	}

	if info["x-cortex-snyk"] != nil {
		c.interpolateSnyk(entity, info["x-cortex-snyk"].(map[string]interface{}))
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
			NotificationsEnabled: MapFetch(ownerMap, "notifications_enabled", false).(bool),
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
}

func (c *CatalogEntityParser) interpolateGit(entity *CatalogEntityData, gitMap map[string]interface{}) {
	if gitMap["github"] != nil {
		githubMap := gitMap["github"].(map[string]interface{})
		entity.Git.Github = CatalogEntityGitGithub{
			Repository: MapFetchToString(githubMap, "repository"),
			BasePath:   MapFetchToString(githubMap, "basePath"),
		}
	}
	if gitMap["gitlab"] != nil {
		gitlabMap := gitMap["gitlab"].(map[string]interface{})
		entity.Git.Gitlab = CatalogEntityGitGitlab{
			Repository: MapFetchToString(gitlabMap, "repository"),
			BasePath:   MapFetchToString(gitlabMap, "basePath"),
		}
	}
	if gitMap["azure"] != nil {
		azureMap := gitMap["azure"].(map[string]interface{})
		entity.Git.Azure = CatalogEntityGitAzureDevOps{
			Project:    MapFetchToString(azureMap, "project"),
			Repository: MapFetchToString(azureMap, "repository"),
			BasePath:   MapFetchToString(azureMap, "basePath"),
		}
	}
	if gitMap["bitbucket"] != nil {
		bitbucketMap := gitMap["bitbucket"].(map[string]interface{})
		entity.Git.BitBucket = CatalogEntityGitBitBucket{
			Repository: MapFetchToString(bitbucketMap, "repository"),
		}
	}
}

func (c *CatalogEntityParser) interpolateJira(entity *CatalogEntityData, jiraMap map[string]interface{}) {
	if jiraMap["defaultJql"] != nil {
		entity.Issues.Jira.DefaultJQL = MapFetchToString(jiraMap, "defaultJQL")
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
			entity.Issues.Jira.Projects = append(entity.Issues.Jira.Labels, label.(string))
		}
	}
	if jiraMap["components"] != nil {
		components := jiraMap["components"].([]interface{})
		for _, component := range components {
			entity.Issues.Jira.Projects = append(entity.Issues.Jira.Components, component.(string))
		}
	}
}

func (c *CatalogEntityParser) interpolateSLOs(entity *CatalogEntityData, slosMap map[string]interface{}) {
	entity.SLOs = CatalogEntitySLOs{}
	if slosMap["lightstep"] != nil {
		c.interpolateLightstep(entity, slosMap["lightstep"].(map[string]interface{}))
	}
	if slosMap["prometheus"] != nil {
		c.interpolatePrometheus(entity, slosMap["prometheus"].([]interface{}))
	}
	// TODO: SignalFX, DataDog, DynaTrace, SumoLogic
}

func (c *CatalogEntityParser) interpolateLightstep(entity *CatalogEntityData, lightstepMap map[string]interface{}) {
	entity.SLOs.Lightstep = CatalogEntitySLOLightstep{
		Streams: []CatalogEntitySLOLightstepStream{},
	}
	if lightstepMap["streams"] != nil {
		streams := lightstepMap["streams"].([]interface{})
		for _, stream := range streams {
			streamMap := stream.(map[string]interface{})
			streamSLO := CatalogEntitySLOLightstepStream{
				StreamID: MapFetchToString(streamMap, "streamId"),
				Targets:  CatalogEntitySLOLightstepTarget{},
			}
			if streamMap["targets"] != nil {
				streamTargetMap := streamMap["targets"].(map[string]interface{})
				if streamTargetMap["latency"] != nil {
					latencies := streamTargetMap["latency"].([]interface{})
					for _, latency := range latencies {
						latencyMap := latency.(map[string]interface{})
						streamSLO.Targets.Latencies = append(streamSLO.Targets.Latencies, CatalogEntitySLOLightstepTargetLatency{
							Percentile: MapFetch(latencyMap, "percentile", 0.0).(float64),
							Target:     MapFetch(latencyMap, "target", 0).(int64),
							SLO:        MapFetch(latencyMap, "slo", 0.0).(float64),
						})
					}
				}
			}
			entity.SLOs.Lightstep.Streams = append(entity.SLOs.Lightstep.Streams, streamSLO)
		}
	}
}

func (c *CatalogEntityParser) interpolatePrometheus(entity *CatalogEntityData, prometheusQueries []interface{}) {
	entity.SLOs.Prometheus = []CatalogEntitySLOPrometheusQuery{}
	for _, query := range prometheusQueries {
		queryMap := query.(map[string]interface{})
		entity.SLOs.Prometheus = append(entity.SLOs.Prometheus, CatalogEntitySLOPrometheusQuery{
			ErrorQuery: MapFetchToString(queryMap, "errorQuery"),
			TotalQuery: MapFetchToString(queryMap, "totalQuery"),
			SLO:        MapFetch(queryMap, "slo", "0.0").(float64),
		})
	}
}

func (c *CatalogEntityParser) interpolateSentry(entity *CatalogEntityData, sentryMap map[string]interface{}) {
	entity.Sentry.Project = MapFetchToString(sentryMap, "project")
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

func (c *CatalogEntityParser) interpolateSnyk(entity *CatalogEntityData, snykMap map[string]interface{}) {
	if snykMap["projects"] != nil {
		projects := snykMap["projects"].([]interface{})
		for _, project := range projects {
			projectMap := project.(map[string]interface{})
			entity.Snyk.Projects = append(entity.Snyk.Projects, CatalogEntitySnykProject{
				ProjectID:    MapFetchToString(projectMap, "projectId"),
				Organization: MapFetchToString(projectMap, "organizationId"),
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
	if saMap["code_cov"] != nil {
		c.interpolateStaticAnalysisCodeCov(entity, saMap["code_cov"].(map[string]interface{}))
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
		Repository: ccMap["repository"].(string),
		Provider:   ccMap["provider"].(string),
	}
}

func (c *CatalogEntityParser) interpolateStaticAnalysisMend(entity *CatalogEntityData, mendMap map[string]interface{}) {
	entity.StaticAnalysis.Mend = CatalogEntityStaticAnalysisMend{}
	applicationIds := mendMap["applicationIds"].([]interface{})
	for _, applicationId := range applicationIds {
		entity.StaticAnalysis.Mend.ApplicationIDs = append(entity.StaticAnalysis.Mend.ApplicationIDs, applicationId.(string))
	}
	projectIds := mendMap["projectIds"].([]interface{})
	for _, projectId := range projectIds {
		entity.StaticAnalysis.Mend.ProjectIDs = append(entity.StaticAnalysis.Mend.ProjectIDs, projectId.(string))
	}
}

func (c *CatalogEntityParser) interpolateStaticAnalysisSonarQube(entity *CatalogEntityData, mendMap map[string]interface{}) {
	entity.StaticAnalysis.SonarQube = CatalogEntityStaticAnalysisSonarQube{
		Project: mendMap["project"].(string),
	}
}

func (c *CatalogEntityParser) interpolateStaticAnalysisVeracode(entity *CatalogEntityData, mendMap map[string]interface{}) {
	entity.StaticAnalysis.Veracode = CatalogEntityStaticAnalysisVeracode{}
	applicationNames := mendMap["applicationNames"].([]interface{})
	for _, applicationName := range applicationNames {
		entity.StaticAnalysis.Veracode.ApplicationNames = append(entity.StaticAnalysis.Veracode.ApplicationNames, applicationName.(string))
	}
	if mendMap["sandboxes"] != nil {
		entity.StaticAnalysis.Veracode.Sandboxes = []CatalogEntityStaticAnalysisVeracodeSandbox{}
		sandboxes := mendMap["sandboxes"].([]interface{})
		for _, sandbox := range sandboxes {
			sandboxMap := sandbox.(map[string]interface{})
			if sandboxMap["applicationName"] != nil || sandboxMap["sandboxName"] != nil {
				sandboxEntity := CatalogEntityStaticAnalysisVeracodeSandbox{
					ApplicationName: sandboxMap["applicationName"].(string),
					SandboxName:     sandboxMap["sandboxName"].(string),
				}
				entity.StaticAnalysis.Veracode.Sandboxes = append(entity.StaticAnalysis.Veracode.Sandboxes, sandboxEntity)
			}
		}
	}
}
