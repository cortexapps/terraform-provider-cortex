package cortex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dghubble/sling"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"strings"
)

type CatalogEntitiesClientInterface interface {
	Get(ctx context.Context, tag string) (*CatalogEntity, error)
	GetFromDescriptor(ctx context.Context, tag string) (*CatalogEntityData, error)
	List(ctx context.Context, params *CatalogEntityListParams) (*CatalogEntitiesResponse, error)
	Upsert(ctx context.Context, req UpsertCatalogEntityRequest) (*CatalogEntityData, error)
	Delete(ctx context.Context, tag string) error
}

type CatalogEntitiesClient struct {
	client *HttpClient
	parser *CatalogEntityParser
}

var _ CatalogEntitiesClientInterface = &CatalogEntitiesClient{}

func (c *CatalogEntitiesClient) Client() *sling.Sling {
	return c.client.client
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// CatalogEntity is the nested response object that is typically returned from the catalog entities endpoints.
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

type CatalogEntityGetDescriptorParams struct {
	Yaml bool `url:"yaml"`
}

type yamlDecoder struct{}

// Decode decodes the Response Body into the value pointed to by v.
// Caller must provide a non-nil v and close the resp.Body.
func (d yamlDecoder) Decode(resp *http.Response, v interface{}) error {
	return yaml.NewDecoder(resp.Body).Decode(v)
}

// jsonDecoder decodes http response JSON into a JSON-tagged struct value.
type jsonDecoder struct{}

// Decode decodes the Response Body into the value pointed to by v.
// Caller must provide a non-nil v and close the resp.Body.
func (d jsonDecoder) Decode(resp *http.Response, v interface{}) error {
	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *CatalogEntitiesClient) GetFromDescriptor(ctx context.Context, tag string) (*CatalogEntityData, error) {
	entity := &CatalogEntityData{}
	entityDescriptorResponse := map[string]interface{}{}
	entityDescriptorResponse["openapi"] = "3.0.1"

	apiError := &ApiError{}
	params := &CatalogEntityGetDescriptorParams{
		Yaml: true,
	}
	uri := Route("catalog_entities", tag+"/openapi")
	cl := c.Client().Get(uri).QueryStruct(params)
	cl = cl.ResponseDecoder(yamlDecoder{})
	response, err := cl.Receive(entityDescriptorResponse, apiError)
	if err != nil {
		return entity, errors.Join(fmt.Errorf("failed getting catalog entity descriptor for %s from %s", tag, uri), err)
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return entity, errors.Join(fmt.Errorf("failed handling response status for %s from %s", tag, uri), err)
	}
	log.Printf("body: %+v", entityDescriptorResponse)
	cl.ResponseDecoder(jsonDecoder{})
	return c.parser.YamlToEntity(entity, entityDescriptorResponse)
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

func (c *CatalogEntitiesClient) Upsert(ctx context.Context, req UpsertCatalogEntityRequest) (*CatalogEntityData, error) {
	entity := &CatalogEntityData{}
	req.OpenApi = "3.0.1"
	upsertResponse := &UpsertCatalogEntityResponse{
		Ok:         false,
		Violations: []CatalogEntityViolation{},
	}
	apiError := &ApiError{}

	// The API requires submitting the request as YAML, so we need to marshal it first.
	bytes, err := yaml.Marshal(req)
	if err != nil {
		return entity, errors.New("could not marshal yaml: " + err.Error())
	}
	body := strings.NewReader(string(bytes))
	response, err := c.Client().
		Set("Content-Type", "application/openapi;charset=UTF-8").
		Set("Accept", "application/json").
		Post(Route("open_api", "")).
		Body(body).
		Receive(upsertResponse, apiError)
	if err != nil {
		return entity, errors.New("could not upsert catalog entity: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		reqYaml, _ := yaml.Marshal(req)
		log.Printf("Failed upserting catalog entity: %+v\n\nRequest:\n%+v\n%+v", err, string(reqYaml), apiError.String())
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
	return c.GetFromDescriptor(ctx, req.Info.Tag)
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
