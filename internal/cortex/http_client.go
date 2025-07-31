package cortex

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/dghubble/sling"
	"github.com/motemen/go-loghttp"
	_ "github.com/motemen/go-loghttp/global" // Just this line!
)

const (
	// UserAgentPrefix is the prefix of the User-Agent header that all terraform REST calls perform.
	UserAgentPrefix = "cortex-terraform-provider"
)

var BaseUris = map[string]string{
	"teams":                "/api/v1/teams/",
	"departments":          "/api/v1/teams/departments/",
	"scorecards":           "/api/v1/scorecards/",
	"catalog_entities":     "/api/v1/catalog/",
	"open_api":             "/api/v1/open-api",
	"resource_definitions": "/api/v1/catalog/definitions/",
}

func Route(domain string, path string) string {
	val, ok := BaseUris[domain]
	if !ok {
		return ""
	}
	return val + path
}

type HttpClient struct {
	ctx        context.Context
	client     *sling.Sling
	yamlClient *sling.Sling
	baseUrl    string
	token      string
	version    string
}

type OptionDelegator func(c *HttpClient) error

// NewClient initializes a new API client for Cortex.
func NewClient(opts ...OptionDelegator) (*HttpClient, error) {
	c := &HttpClient{}
	for _, f := range opts {
		if err := f(c); err != nil {
			return nil, err
		}
	}

	hc := &http.Client{}
	if os.Getenv("HTTP_DEBUG") == "1" {
		hc.Transport = &loghttp.Transport{}
	}
	c.client = sling.New().Doer(hc).Base(c.baseUrl).
		Set("User-Agent", fmt.Sprintf("%s (%s)", UserAgentPrefix, c.version)).
		Set("Authorization", fmt.Sprintf("Bearer %s", c.token)).
		ResponseDecoder(jsonDecoder{})
	c.yamlClient = sling.New().Doer(hc).Base(c.baseUrl).
		Set("User-Agent", fmt.Sprintf("%s (%s)", UserAgentPrefix, c.version)).
		Set("Authorization", fmt.Sprintf("Bearer %s", c.token)).
		ResponseDecoder(yamlDecoder{})

	return c, nil
}

func WithVersion(version string) func(*HttpClient) error {
	return func(c *HttpClient) error {
		if version == "" {
			return errors.New("cannot specify empty version")
		}
		c.version = version
		return nil
	}
}

// WithContext Specify the context for the cortex client to use.
func WithContext(ctx context.Context) func(*HttpClient) error {
	return func(c *HttpClient) error {
		c.ctx = ctx
		return nil
	}
}

// WithURL Specify the base URL for the cortex client to connect to.
func WithURL(baseUrl string) func(*HttpClient) error {
	return func(c *HttpClient) error {
		if baseUrl == "" {
			return errors.New("cannot specify empty API Base URL")
		}
		if _, err := url.Parse(baseUrl); err != nil {
			return err
		}
		c.baseUrl = baseUrl
		return nil
	}
}

// WithToken Specify the API token for the cortex client to use.
func WithToken(token string) func(*HttpClient) error {
	return func(c *HttpClient) error {
		if token == "" {
			return errors.New("cannot specify empty token")
		}
		c.token = token
		return nil
	}
}

func (c *HttpClient) handleResponseStatus(response *http.Response, apiError *ApiError) error {
	switch code := response.StatusCode; {
	case code >= 200 && code <= 299:
		return nil
	case code == 404:
		return ApiErrorNotFound
	case code == 401:
		return fmt.Errorf("%s\n%s", ApiErrorUnauthorized, apiError)
	default:
		return fmt.Errorf("%d request failed with error: %+v", code, apiError.String())
	}
}

func (c *HttpClient) Ping(ctx context.Context) error {
	apiError := new(ApiError)
	response, err := c.client.Get("/").Receive(nil, apiError)
	if err != nil {
		return err
	}
	return c.handleResponseStatus(response, apiError)
}

func (c *HttpClient) Client() *sling.Sling {
	return c.client.New()
}

func (c *HttpClient) YamlClient() *sling.Sling {
	return c.yamlClient.New()
}

/********** Client Interfaces **********/

func (c *HttpClient) CatalogEntities() CatalogEntitiesClientInterface {
	return &CatalogEntitiesClient{client: c, parser: &CatalogEntityParser{}}
}

func (c *HttpClient) CatalogEntityCustomData() CatalogEntityCustomDataClientInterface {
	return &CatalogEntityCustomDataClient{client: c}
}

func (c *HttpClient) CatalogEntityOpenAPI() CatalogEntityOpenAPIClientInterface {
	return &CatalogEntityOpenAPIClient{client: c}
}

func (c *HttpClient) Teams() TeamsClientInterface {
	return &TeamsClient{client: c}
}

func (c *HttpClient) Departments() DepartmentsClientInterface {
	return &DepartmentsClient{client: c}
}

func (c *HttpClient) Scorecards() ScorecardsClientInterface {
	return &ScorecardsClient{client: c}
}

func (c *HttpClient) ResourceDefinitions() ResourceDefinitionsClientInterface {
	return &ResourceDefinitionsClient{client: c}
}
