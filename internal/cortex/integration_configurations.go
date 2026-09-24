package cortex

import (
	"context"
	"fmt"
	"net/url"
)

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// MultiInstanceConfiguration is a configuration of an integration that supports several configurations per tenant.
type MultiInstanceConfiguration interface {
	GetAlias() string
}

// MultiInstanceClientInterface is the API of a multi-instance integration.
type MultiInstanceClientInterface[C any, U any, R MultiInstanceConfiguration] interface {
	List(ctx context.Context) ([]R, error)
	Get(ctx context.Context, alias string) (R, error)
	Create(ctx context.Context, alias string, req C) (R, error)
	Update(ctx context.Context, currentAlias string, newAlias string, req U) (R, error)
	Delete(ctx context.Context, alias string) error
}

// SingleInstanceClientInterface is the API of a single-instance integration.
type SingleInstanceClientInterface[C any, R any] interface {
	Get(ctx context.Context) (R, error)
	Create(ctx context.Context, req C) (R, error)
	Replace(ctx context.Context, req C) (R, error)
	Delete(ctx context.Context) error
}

// configurationsResponse wraps list, create, update, and replace responses. For multi-instance integrations, create
// and update return every configuration in the tenant, not only the one that changed.
type configurationsResponse[R any] struct {
	Configurations []R `json:"configurations"`
}

/***********************************************************************************************************************
 * Multi-instance client
 **********************************************************************************************************************/

// MultiInstanceClient calls the configuration routes that every multi-instance integration shares.
type MultiInstanceClient[C any, U any, R MultiInstanceConfiguration] struct {
	client *HttpClient
	domain string // key in BaseUris
	name   string // name in error messages
}

// List calls GET <base>/configurations.
func (c *MultiInstanceClient[C, U, R]) List(ctx context.Context) ([]R, error) {
	response := configurationsResponse[R]{}
	apiError := ApiError{}

	body, err := c.client.Client().Get(Route(c.domain, "configurations")).Receive(&response, &apiError)
	if err != nil {
		return nil, fmt.Errorf("failed listing %s configurations: %+v", c.name, err)
	}
	if err := c.client.handleResponseStatus(body, &apiError); err != nil {
		return nil, fmt.Errorf("failed listing %s configurations: %w", c.name, err)
	}
	return response.Configurations, nil
}

// Get finds a configuration by alias in the list response. GET <base>/configuration/:alias returns 400, not 404,
// for an unknown alias, so the list gives a reliable not-found signal.
func (c *MultiInstanceClient[C, U, R]) Get(ctx context.Context, alias string) (R, error) {
	configurations, err := c.List(ctx)
	if err != nil {
		var zero R
		return zero, err
	}
	return findByAlias(c.name, configurations, alias)
}

// Create calls POST <base>/configuration and returns the new configuration, found by alias in the response.
func (c *MultiInstanceClient[C, U, R]) Create(ctx context.Context, alias string, req C) (R, error) {
	response := configurationsResponse[R]{}
	apiError := ApiError{}
	var zero R

	body, err := c.client.Client().Post(Route(c.domain, "configuration")).BodyJSON(&req).Receive(&response, &apiError)
	if err != nil {
		return zero, fmt.Errorf("failed creating %s configuration: %+v", c.name, err)
	}
	if err := c.client.handleResponseStatus(body, &apiError); err != nil {
		return zero, fmt.Errorf("failed creating %s configuration: %w", c.name, err)
	}
	return findByAlias(c.name, response.Configurations, alias)
}

// Update calls PUT <base>/configuration/:currentAlias. The body can rename the configuration to newAlias.
func (c *MultiInstanceClient[C, U, R]) Update(ctx context.Context, currentAlias string, newAlias string, req U) (R, error) {
	response := configurationsResponse[R]{}
	apiError := ApiError{}
	var zero R

	body, err := c.client.Client().Put(Route(c.domain, "configuration/"+url.PathEscape(currentAlias))).BodyJSON(&req).Receive(&response, &apiError)
	if err != nil {
		return zero, fmt.Errorf("failed updating %s configuration: %+v", c.name, err)
	}
	if err := c.client.handleResponseStatus(body, &apiError); err != nil {
		return zero, fmt.Errorf("failed updating %s configuration: %w", c.name, err)
	}
	return findByAlias(c.name, response.Configurations, newAlias)
}

// Delete calls DELETE <base>/configuration/:alias.
func (c *MultiInstanceClient[C, U, R]) Delete(ctx context.Context, alias string) error {
	apiError := ApiError{}

	body, err := c.client.Client().Delete(Route(c.domain, "configuration/"+url.PathEscape(alias))).Receive(nil, &apiError)
	if err != nil {
		return fmt.Errorf("failed deleting %s configuration: %+v", c.name, err)
	}
	if err := c.client.handleResponseStatus(body, &apiError); err != nil {
		return fmt.Errorf("failed deleting %s configuration: %w", c.name, err)
	}
	return nil
}

func findByAlias[R MultiInstanceConfiguration](name string, configurations []R, alias string) (R, error) {
	for _, configuration := range configurations {
		if configuration.GetAlias() == alias {
			return configuration, nil
		}
	}
	var zero R
	return zero, fmt.Errorf("%s configuration %s: %w", name, alias, ApiErrorNotFound)
}

/***********************************************************************************************************************
 * Single-instance client
 **********************************************************************************************************************/

// SingleInstanceClient calls the configuration routes that every single-instance integration shares. A tenant has at
// most one configuration, so no route takes an alias.
type SingleInstanceClient[C any, R any] struct {
	client *HttpClient
	domain string
	name   string
}

// Get calls GET <base>/default-configuration. The API returns 404 when no configuration exists.
func (c *SingleInstanceClient[C, R]) Get(ctx context.Context) (R, error) {
	var response R
	apiError := ApiError{}

	var zero R

	body, err := c.client.Client().Get(Route(c.domain, "default-configuration")).Receive(&response, &apiError)
	if err != nil {
		return zero, fmt.Errorf("failed getting %s configuration: %+v", c.name, err)
	}
	if err := c.client.handleResponseStatus(body, &apiError); err != nil {
		return zero, fmt.Errorf("failed getting %s configuration: %w", c.name, err)
	}
	return response, nil
}

// Create calls POST <base>/configuration. The API rejects it when a configuration exists.
func (c *SingleInstanceClient[C, R]) Create(ctx context.Context, req C) (R, error) {
	return c.write(false, req)
}

// Replace calls PUT <base>/configuration with the full configuration. The API creates one when none exists.
func (c *SingleInstanceClient[C, R]) Replace(ctx context.Context, req C) (R, error) {
	return c.write(true, req)
}

func (c *SingleInstanceClient[C, R]) write(replace bool, req C) (R, error) {
	response := configurationsResponse[R]{}
	apiError := ApiError{}
	var zero R

	request, verb := c.client.Client().Post(Route(c.domain, "configuration")), "creating"
	if replace {
		request, verb = c.client.Client().Put(Route(c.domain, "configuration")), "replacing"
	}
	body, err := request.BodyJSON(&req).Receive(&response, &apiError)
	if err != nil {
		return zero, fmt.Errorf("failed %s %s configuration: %+v", verb, c.name, err)
	}
	if err := c.client.handleResponseStatus(body, &apiError); err != nil {
		return zero, fmt.Errorf("failed %s %s configuration: %w", verb, c.name, err)
	}
	if len(response.Configurations) == 0 {
		return zero, fmt.Errorf("failed %s %s configuration: response contained no configuration", verb, c.name)
	}
	return response.Configurations[0], nil
}

// Delete calls DELETE <base>/configurations. The API succeeds also when no configuration exists.
func (c *SingleInstanceClient[C, R]) Delete(ctx context.Context) error {
	apiError := ApiError{}

	body, err := c.client.Client().Delete(Route(c.domain, "configurations")).Receive(nil, &apiError)
	if err != nil {
		return fmt.Errorf("failed deleting %s configuration: %+v", c.name, err)
	}
	if err := c.client.handleResponseStatus(body, &apiError); err != nil {
		return fmt.Errorf("failed deleting %s configuration: %w", c.name, err)
	}
	return nil
}
