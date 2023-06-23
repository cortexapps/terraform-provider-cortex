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
		for key, value := range info["x-cortex-custom-metadata"].(map[string]interface{}) {
			entity.Metadata[key] = value
		}
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

	if info["x-cortex-oncall"] != nil {
		onCallMap := info["x-cortex-oncall"].(map[string]interface{})
		c.interpolateOnCall(entity, onCallMap)
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
		entity.OnCall.Pagerduty = CatalogEntityOnCallPagerduty{
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
							Target:     MapFetch(latencyMap, "target", 0).(int),
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
			SLO:        MapFetch(queryMap, "slo", 0.0).(float64),
		})
	}
}

func (c *CatalogEntityParser) interpolateSentry(entity *CatalogEntityData, sentryMap map[string]interface{}) {
	entity.Sentry.Project = MapFetchToString(sentryMap, "project")
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
