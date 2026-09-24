package integrations_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/cortexapps/terraform-provider-cortex/internal/cortex/integrations"
	"github.com/stretchr/testify/assert"
)

func TestJiraRequestWireFormat(t *testing.T) {
	cloud, err := json.Marshal(integrations.JiraConfigurationRequest{
		Type: integrations.JiraTypeCloudBasic, Alias: "jira", Subdomain: "acme", Email: "bot@acme.com", ApiToken: "t", BaseUrl: "atlassian.net",
	})
	assert.Nil(t, err)
	assert.JSONEq(t, `{"type":"CLOUD_BASIC","alias":"jira","isDefault":false,"subdomain":"acme","email":"bot@acme.com","apiToken":"t","baseUrl":"atlassian.net"}`, string(cloud))

	onPrem, err := json.Marshal(integrations.JiraConfigurationRequest{
		Type: integrations.JiraTypeOnPremBasic, Alias: "jira", IsDefault: true, Host: "https://jira.invalid", Username: "bot", Password: "p",
	})
	assert.Nil(t, err)
	assert.JSONEq(t, `{"type":"ON_PREM_BASIC","alias":"jira","isDefault":true,"host":"https://jira.invalid","username":"bot","password":"p"}`, string(onPrem))
}

func TestJiraResponseParsing(t *testing.T) {
	var cfg integrations.JiraConfiguration
	err := json.Unmarshal([]byte(`{"alias":"jira","isDefault":true,"type":"CLOUD_SCOPED","subdomain":"acme","baseUrl":"api.atlassian.com/ex/jira","lastFour":"a1b2","email":"bot@acme.com","username":null}`), &cfg)
	assert.Nil(t, err)
	assert.Equal(t, integrations.JiraConfiguration{Alias: "jira", IsDefault: true, Type: integrations.JiraTypeCloudScoped, Subdomain: "acme", BaseUrl: "api.atlassian.com/ex/jira", LastFour: "a1b2", Email: "bot@acme.com"}, cfg)
}

func TestGitlabRequestWireFormat(t *testing.T) {
	create, err := json.Marshal(integrations.CreateGitlabConfigurationRequest{Alias: "gl", PersonalAccessToken: "t", GroupNames: []string{}})
	assert.Nil(t, err)
	assert.JSONEq(t, `{"alias":"gl","isDefault":false,"personalAccessToken":"t","hidePersonalProjects":false,"groupNames":[]}`, string(create))

	update, err := json.Marshal(integrations.UpdateGitlabConfigurationRequest{Alias: "gl", GroupNames: []string{"platform"}, PersonalAccessToken: "t"})
	assert.Nil(t, err)
	assert.JSONEq(t, `{"alias":"gl","isDefault":false,"hidePersonalProjects":false,"groupNames":["platform"],"personalAccessToken":"t"}`, string(update))
}

func TestIncidentIoRequestWireFormat(t *testing.T) {
	create, err := json.Marshal(integrations.CreateIncidentIoConfigurationRequest{Alias: "inc", ApiKey: "k"})
	assert.Nil(t, err)
	assert.JSONEq(t, `{"alias":"inc","isDefault":false,"apiKey":"k"}`, string(create))
}

func TestJiraClientRoutes(t *testing.T) {
	c, teardown, err := setupClient(
		"/api/v1/jira/configurations",
		map[string]any{"configurations": []integrations.JiraConfiguration{{Alias: "jira", Type: integrations.JiraTypeCloudBasic}}},
		AssertRequestMethod(t, "GET"),
	)
	assert.Nil(t, err)
	defer teardown()

	res, err := integrations.Jira(c).Get(context.Background(), "jira")
	assert.Nil(t, err)
	assert.Equal(t, integrations.JiraTypeCloudBasic, res.Type)
}

func TestIncidentIoClientRoute(t *testing.T) {
	c, teardown, err := setupClient(
		"/api/v1/incidentio/configurations",
		map[string]any{"configurations": []integrations.IncidentIoConfiguration{{Alias: "inc", LastFour: "a1b2"}}},
		AssertRequestMethod(t, "GET"),
	)
	assert.Nil(t, err)
	defer teardown()

	res, err := integrations.IncidentIo(c).List(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, "a1b2", res[0].LastFour)
}

func TestGitlabClientRoute(t *testing.T) {
	c, teardown, err := setupClient(
		"/api/v1/gitlab/configuration/gl",
		map[string]any{"configurations": []integrations.GitlabConfiguration{{Alias: "gl", LastFour: "a1b2"}}},
		AssertRequestMethod(t, "PUT"),
	)
	assert.Nil(t, err)
	defer teardown()

	res, err := integrations.Gitlab(c).Update(context.Background(), "gl", "gl", integrations.UpdateGitlabConfigurationRequest{Alias: "gl", GroupNames: []string{}})
	assert.Nil(t, err)
	assert.Equal(t, "a1b2", res.LastFour)
}

// The backend rejects null for its required lists, so a nil slice must go out as [].
func TestRequestsSendEmptyListsForNil(t *testing.T) {
	for name, req := range map[string]any{
		"datadog create": integrations.CreateDatadogConfigurationRequest{},
		"datadog update": integrations.UpdateDatadogConfigurationRequest{},
		"gitlab create":  integrations.CreateGitlabConfigurationRequest{},
		"gitlab update":  integrations.UpdateGitlabConfigurationRequest{},
	} {
		t.Run(name, func(t *testing.T) {
			body, err := json.Marshal(req)
			assert.Nil(t, err)
			assert.NotContains(t, string(body), "null")
		})
	}
}
