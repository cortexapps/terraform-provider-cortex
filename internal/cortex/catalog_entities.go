package cortex

import (
	"context"
	"errors"
	"fmt"
	"github.com/dghubble/sling"
	"gopkg.in/yaml.v3"
)

type CatalogEntitiesClientInterface interface {
	Get(ctx context.Context, tag string) (*CatalogEntity, error)
	List(ctx context.Context, params *CatalogEntityListParams) (*CatalogEntitiesResponse, error)
	Delete(ctx context.Context, tag string) error
}

type CatalogEntitiesClient struct {
	client *HttpClient
}

var _ CatalogEntitiesClientInterface = &CatalogEntitiesClient{}

func (c *CatalogEntitiesClient) Client() *sling.Sling {
	return c.client.client
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

type CatalogEntity struct {
	Tag          string                 `json:"tag" yaml:"x-cortex-tag"`
	Title        string                 `json:"title" yaml:"title"`
	Description  string                 `json:"description" yaml:"description"`
	Type         string                 `json:"type" yaml:"x-cortex-type"`
	Groups       []string               `json:"groups" yaml:"x-cortex-groups"`
	Links        []CatalogEntityLink    `json:"links" yaml:"x-cortex-link"`
	Metadata     map[string]interface{} `json:"metadata" yaml:"x-cortex-custom-metadata"`
	Dependencies []string               `json:"dependencies" yaml:"x-cortex-dependency"`
	Ownership    CatalogEntityOwnership `json:"ownership" yaml:"x-cortex-owners"`
}

type CatalogEntityOwnership struct {
	Emails        []CatalogEntityEmail        `json:"emails"`
	Groups        []CatalogEntityGroup        `json:"groups"`
	SlackChannels []CatalogEntitySlackChannel `json:"slackChannels"`
}

type CatalogEntityEmail struct {
	Email       string `json:"email" yaml:"email"`
	Description string `json:"description" yaml:"description"`
}

type CatalogEntityGroup struct {
	GroupName   string `json:"group" yaml:"groupName"`
	Description string `json:"description" yaml:"description"`
	Provider    string `json:"provider" yaml:"provider"`
}

type CatalogEntitySlackChannel struct {
	Channel              string `json:"channel" yaml:"channel"`
	Description          string `json:"description" yaml:"description"`
	NotificationsEnabled bool   `json:"notificationsEnabled" yaml:"notificationsEnabled"`
}

type CatalogEntityLink struct {
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
	Url  string `json:"url" yaml:"url"`
}

// CatalogEntityDefinition Required for non-service catalog entities.
type CatalogEntityDefinition struct {
	Version      string `json:"version" yaml:"version"`
	Distribution string `json:"distribution" yaml:"distribution"`
}

type CatalogEntityViolation struct {
	Description   string   `json:"description"`
	ViolationType string   `json:"violationType"`
	StartLine     int      `json:"startLine"`
	EndLine       int      `json:"endLine"`
	Paths         []string `json:"paths"`
	Pointer       string   `json:"pointer"`
	RuleLink      string   `json:"ruleLink"`
	Title         string   `json:"title"`
}

func (v *CatalogEntityViolation) String() string {
	return fmt.Sprintf("%s (%s): %s (L%d:L%d) - %s", v.Title, v.ViolationType, v.Description, v.StartLine, v.EndLine, v.Pointer)
}

/***********************************************************************************************************************
 * GET /api/v1/catalog/:tag
 **********************************************************************************************************************/

func (c *CatalogEntitiesClient) Get(ctx context.Context, tag string) (*CatalogEntity, error) {
	catalogEntityResponse := &CatalogEntity{}
	apiError := &ApiError{}
	response, err := c.Client().Get(Route("catalog_entities", tag)).Receive(catalogEntityResponse, apiError)
	if err != nil {
		return catalogEntityResponse, errors.New("could not get catalog entity: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return catalogEntityResponse, errors.Join(errors.New("Failed getting catalog entity: "), err)
	}

	return catalogEntityResponse, nil
}

func (c *CatalogEntitiesClient) GetFromDescriptor(ctx context.Context, tag string) (*CatalogEntity, error) {
	entity := &CatalogEntity{}
	entityDescriptorResponse := ""

	apiError := &ApiError{}
	response, err := c.Client().Get(Route("catalog", tag+"/openapi")).Receive(entityDescriptorResponse, apiError)

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return entity, errors.Join(errors.New("Failed getting catalog entity descriptor: "), err)
	}
	yamlEntity := map[string]interface{}{}
	err = yaml.Unmarshal([]byte(entityDescriptorResponse), yamlEntity)
	if err != nil {
		return entity, errors.Join(errors.New("Failed decoding catalog entity descriptor into YAML: "), err)
	}

	return c.YamlToEntity(ctx, entity, yamlEntity)
}

// YamlToEntity converts YAML into a CatalogEntity, from the following specification example:
/*
openapi: 3.0.0
info:
  title: Chat Service
  description: Chat service is responsible for handling chat feature.
  x-cortex-tag: chat-service
  x-cortex-type: service
  x-cortex-link:
    - name: Chat ServiceAPI Spec
      type: OPENAPI
      url: ./docs/chat-service-openapi-spec.yaml
  x-cortex-groups:
    - python-services
  x-cortex-owners:
    - type: group
      name: Delta
      provider: OKTA
      description: Delta Team
    - type: slack
      channel: delta-team
      notificationsEnabled: true
  x-cortex-custom-metadata:
    core-service: true
  x-cortex-dependency:
    tag: authentication-service
    tag: chat-database
  x-cortex-git:
    github:
      repository: org/chat-service
  x-cortex-oncall:
    pagerduty:
      id: ASDF1234
      type: SCHEDULE
  x-cortex-apm:
    datadog:
      monitors:
        - 12345
  x-cortex-issues:
    jira:
      projects:
        - CS
*/
func (c *CatalogEntitiesClient) YamlToEntity(ctx context.Context, entity *CatalogEntity, yamlEntity map[string]interface{}) (*CatalogEntity, error) {
	info := yamlEntity["info"].(map[string]interface{})

	entity.Title = info["title"].(string)
	entity.Description = info["description"].(string)
	entity.Tag = info["x-cortex-tag"].(string)
	entity.Type = info["x-cortex-type"].(string)

	entity.Links = []CatalogEntityLink{}
	if info["x-cortex-link"] != nil {
		for _, link := range info["x-cortex-link"].([]interface{}) {
			linkMap := link.(map[string]interface{})
			entity.Links = append(entity.Links, CatalogEntityLink{
				Name: linkMap["name"].(string),
				Type: linkMap["type"].(string),
				Url:  linkMap["url"].(string),
			})
		}
	}

	entity.Groups = []string{}
	if info["x-cortex-groups"] != nil {
		for _, group := range info["x-cortex-groups"].([]interface{}) {
			entity.Groups = append(entity.Groups, group.(string))
		}
	}

	entity.Ownership = CatalogEntityOwnership{
		Groups:        []CatalogEntityGroup{},
		SlackChannels: []CatalogEntitySlackChannel{},
		Emails:        []CatalogEntityEmail{},
	}
	if info["x-cortex-owners"] != nil {
		for _, owner := range info["x-cortex-owners"].([]interface{}) {
			ownerMap := owner.(map[string]interface{})
			if ownerMap["type"] == "group" {
				entity.Ownership.Groups = append(entity.Ownership.Groups, CatalogEntityGroup{
					GroupName:   ownerMap["name"].(string),
					Provider:    ownerMap["provider"].(string),
					Description: ownerMap["description"].(string),
				})
			} else if ownerMap["type"] == "slack" {
				entity.Ownership.SlackChannels = append(entity.Ownership.SlackChannels, CatalogEntitySlackChannel{
					Channel:              ownerMap["channel"].(string),
					Description:          ownerMap["description"].(string),
					NotificationsEnabled: ownerMap["notificationsEnabled"].(bool),
				})
			} else if ownerMap["type"] == "email" {
				entity.Ownership.Emails = append(entity.Ownership.Emails, CatalogEntityEmail{
					Email: ownerMap["email"].(string),
				})
			}
		}
	}

	entity.Metadata = map[string]interface{}{}
	if info["x-cortex-custom-metadata"] != nil {
		for key, value := range info["x-cortex-custom-metadata"].(map[string]interface{}) {
			entity.Metadata[key] = value
		}
	}

	entity.Dependencies = []string{}
	if info["x-cortex-dependency"] != nil {
		for _, dependency := range info["x-cortex-dependency"].([]interface{}) {
			entity.Dependencies = append(entity.Dependencies, dependency.(string))
		}
	}

	/*
		TODO: handle these
		  x-cortex-children:
			# children can be of type service, resource, or domain
			- tag: chat-service
			- tag: chat-database
		  x-cortex-domain-parents:
			# parents can be of type domain only
			- tag: payments-domain
			- tag: web-domain
		  x-cortex-git:
		    github:
		      repository: org/chat-service
		  x-cortex-oncall:
		    pagerduty:
		      id: ASDF1234
		      type: SCHEDULE
		  x-cortex-apm:
		    datadog:
		      monitors:
		        - 12345
		  x-cortex-issues:
		    jira:
		      projects:
		        - CS
	*/

	return entity, nil
}

/***********************************************************************************************************************
 * GET /api/v1/catalog
 **********************************************************************************************************************/

// CatalogEntityListParams are the query parameters for the GET /v1/catalog endpoint.
type CatalogEntityListParams struct {
	Groups          []string `url:"groups,omitempty"`
	Types           []string `url:"types,omitempty"`
	GitRepositories []string `url:"gitRepositories,omitempty"`
}

// CatalogEntitiesResponse is the response from the GET /v1/scorecards endpoint.
type CatalogEntitiesResponse struct {
	Entities []CatalogEntity `json:"entities" yaml:"entities"`
}

// List retrieves a list of scorecards based on a query.
func (c *CatalogEntitiesClient) List(ctx context.Context, params *CatalogEntityListParams) (*CatalogEntitiesResponse, error) {
	entitiesResponse := &CatalogEntitiesResponse{}
	apiError := &ApiError{}

	response, err := c.Client().Get(Route("catalog_entities", "")).QueryStruct(params).Receive(entitiesResponse, apiError)
	if err != nil {
		return nil, errors.New("could not get entities: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return nil, err
	}

	return entitiesResponse, nil
}

/***********************************************************************************************************************
 * POST /api/v1/open-api
 **********************************************************************************************************************/

type UpsertCatalogEntityRequest struct {
	Info    CatalogEntityData `json:"info"`
	OpenApi string            `json:"openapi"`
}

type UpsertCatalogEntityResponse struct {
	Ok         bool                     `json:"ok"`
	Violations []CatalogEntityViolation `json:"violations"`
}

// CatalogEntityData is a struct used in upsert requests in the `info` parameter, since its structure does not
// match the structure of the CatalogEntity struct in responses.
type CatalogEntityData struct {
	Title       string                  `json:"title"`
	Description string                  `json:"description,omitempty"`
	Tag         string                  `json:"x-cortex-tag"`
	Ownership   []CatalogEntityOwner    `json:"x-cortex-owners,omitempty"`
	Groups      []string                `json:"x-cortex-groups,omitempty"` // TODO: is this -groups or -service-groups? docs unclear
	Links       []CatalogEntityLink     `json:"x-cortex-link,omitempty"`
	Metadata    map[string]interface{}  `json:"x-cortex-custom-metadata,omitempty"`
	Type        string                  `json:"x-cortex-type,omitempty"`
	Definition  CatalogEntityDefinition `json:"x-cortex-definition,omitempty"`

	// Various generic integration attributes
	Alerts         []CatalogEntityAlert        `json:"x-cortex-alerts,omitempty"`
	Apm            CatalogEntityApm            `json:"x-cortex-apm,omitempty"`
	Dashboards     CatalogEntityDashboards     `json:"x-cortex-dashboards,omitempty"`
	Git            CatalogEntityGit            `json:"x-cortex-git,omitempty"`
	Issues         CatalogEntityIssues         `json:"x-cortex-issues,omitempty"`
	OnCall         CatalogEntityOnCall         `json:"x-cortex-oncall,omitempty"`
	SLOs           CatalogEntitySLOs           `json:"x-cortex-slos,omitempty"`
	StaticAnalysis CatalogEntityStaticAnalysis `json:"x-cortex-static-analysis,omitempty"`

	// Integration-specific things
	BugSnag   CatalogEntityBugSnag   `json:"x-cortex-bugsnag,omitempty"`
	Checkmarx CatalogEntityCheckmarx `json:"x-cortex-checkmarx,omitempty"`
	Rollbar   CatalogEntityRollbar   `json:"x-cortex-rollbar,omitempty"`
	Sentry    CatalogEntitySentry    `json:"x-cortex-sentry,omitempty"`
	Snyk      CatalogEntitySnyk      `json:"x-cortex-snyk,omitempty"`
}

func (c *CatalogEntitiesClient) Upsert(ctx context.Context, req UpsertCatalogEntityRequest) (*CatalogEntity, error) {
	entity := &CatalogEntity{}
	upsertResponse := &UpsertCatalogEntityResponse{
		Ok:         false,
		Violations: []CatalogEntityViolation{},
	}
	apiError := &ApiError{}

	response, err := c.Client().Post(Route("open_api", "")).BodyJSON(req).Receive(upsertResponse, apiError)
	if err != nil {
		return entity, errors.New("could not upsert catalog entity: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return entity, err
	}

	// coerce violations into an error
	if len(upsertResponse.Violations) > 0 {
		o := ""
		for _, v := range upsertResponse.Violations {
			o += v.String() + "\n"
		}
		return entity, errors.New(o)
	}

	// re-fetch the catalog entity, since it's not returned here
	entity, err = c.Get(ctx, req.Info.Tag)
	return entity, err
}

/***********************************************************************************************************************
 * DELETE /api/v1/catalog/:tag - Delete a catalog entity
 **********************************************************************************************************************/

type DeleteCatalogEntityResponse struct{}

func (c *CatalogEntitiesClient) Delete(ctx context.Context, tag string) error {
	entityResponse := &DeleteCatalogEntityResponse{}
	apiError := &ApiError{}

	response, err := c.Client().Delete(Route("catalog_entities", tag)).Receive(entityResponse, apiError)
	if err != nil {
		return errors.New("could not delete catalog entity: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
