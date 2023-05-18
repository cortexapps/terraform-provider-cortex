package cortex

import (
	"context"
	"errors"
	"fmt"
	"github.com/dghubble/sling"
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
	Tag         string                 `json:"tag" yaml:"x-cortex-tag"`
	Description string                 `json:"description" yaml:"description"`
	Groups      []string               `json:"groups" yaml:"x-cortex-groups"`
	Links       []CatalogEntityLink    `json:"links" yaml:"x-cortex-link"`
	Metadata    map[string]interface{} `json:"metadata" yaml:"x-cortex-custom-metadata"`
	Ownership   CatalogEntityOwnership `json:"ownership" yaml:"x-cortex-owners"`
}

type CatalogEntityOwnership struct {
	Emails        []CatalogEntityEmail      `json:"emails"`
	Groups        []CatalogEntityGroup      `json:"groups"`
	SlackChannels CatalogEntitySlackChannel `json:"slackChannels"`
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
