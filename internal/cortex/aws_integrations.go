package cortex

import (
	"context"
)

type AwsConfiguration struct {
	AccountID    string `json:"accountId"`
	IAMRole      string `json:"iamRole"`
	AccountAlias string `json:"accountAlias,omitempty"`
}

type AwsType struct {
	Type       string `json:"type"`
	Configured bool   `json:"configured"`
}

type AwsTypesResponse struct {
	Types []AwsType `json:"types"`
}

type AwsClientInterface interface {
	GetConfiguration(ctx context.Context, accountId string) (*AwsConfiguration, error)
	CreateConfiguration(ctx context.Context, config AwsConfiguration) (*AwsConfiguration, error)
	UpdateConfiguration(ctx context.Context, config AwsConfiguration) (*AwsConfiguration, error)
	DeleteConfiguration(ctx context.Context, accountId string) error

	GetTypes(ctx context.Context) ([]AwsType, error)
	UpdateTypes(ctx context.Context, types []AwsType) ([]AwsType, error)
}

type AwsClient struct {
	client *HttpClient
}

func (c *AwsClient) GetConfiguration(ctx context.Context, accountId string) (*AwsConfiguration, error) {
	apiError := new(ApiError)
	res := new(AwsConfiguration)

	response, err := c.client.Client().Get(Route("aws_configurations", accountId)).Receive(res, apiError)
	if err != nil {
		return nil, err
	}
	if err := c.client.handleResponseStatus(response, apiError); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *AwsClient) CreateConfiguration(ctx context.Context, config AwsConfiguration) (*AwsConfiguration, error) {
	apiError := new(ApiError)
	res := new(AwsConfiguration)

	response, err := c.client.Client().Post(Route("aws_configurations", "")).BodyJSON(&config).Receive(res, apiError)
	if err != nil {
		return nil, err
	}
	if err := c.client.handleResponseStatus(response, apiError); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *AwsClient) UpdateConfiguration(ctx context.Context, config AwsConfiguration) (*AwsConfiguration, error) {
	// POST /api/v1/aws/configurations replaces or updates the config for the account.
	return c.CreateConfiguration(ctx, config)
}

func (c *AwsClient) DeleteConfiguration(ctx context.Context, accountId string) error {
	apiError := new(ApiError)

	response, err := c.client.Client().Delete(Route("aws_configurations", accountId)).Receive(nil, apiError)
	if err != nil {
		return err
	}
	if err := c.client.handleResponseStatus(response, apiError); err != nil {
		return err
	}

	return nil
}

func (c *AwsClient) GetTypes(ctx context.Context) ([]AwsType, error) {
	apiError := new(ApiError)
	res := new(AwsTypesResponse)

	response, err := c.client.Client().Get(Route("aws_types", "")).Receive(res, apiError)
	if err != nil {
		return nil, err
	}
	if err := c.client.handleResponseStatus(response, apiError); err != nil {
		return nil, err
	}

	return res.Types, nil
}

func (c *AwsClient) UpdateTypes(ctx context.Context, types []AwsType) ([]AwsType, error) {
	apiError := new(ApiError)
	res := new(AwsTypesResponse)

	req := AwsTypesResponse{Types: types}
	response, err := c.client.Client().Put(Route("aws_types", "")).BodyJSON(&req).Receive(res, apiError)
	if err != nil {
		return nil, err
	}
	if err := c.client.handleResponseStatus(response, apiError); err != nil {
		return nil, err
	}

	return res.Types, nil
}
