package cortex

import (
	"context"
	"errors"
	"fmt"

	"github.com/dghubble/sling"
)

type EntityRelationshipsClientInterface interface {
	GetDestinations(ctx context.Context, sourceTag string, relationshipTypeTag string) ([]string, error)
	SetDestinations(ctx context.Context, sourceTag string, relationshipTypeTag string, destinationTags []string) ([]string, error)
}

type EntityRelationshipsClient struct {
	client *HttpClient
}

var _ EntityRelationshipsClientInterface = &EntityRelationshipsClient{}

func (c *EntityRelationshipsClient) Client() *sling.Sling {
	return c.client.Client()
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// RelatedEntityDetails is one element in the GET/PUT response.
type RelatedEntityDetails struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Tag         string `json:"tag"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
}

type EntityDestinationsResponse struct {
	Destinations []RelatedEntityDetails `json:"destinations"`
}

type SetEntityDestinationsRequest struct {
	Destinations []string `json:"destinations"`
}

/***********************************************************************************************************************
 * GET /api/v1/catalog/:sourceTag/relationships/:relationshipTypeTag/destinations
 **********************************************************************************************************************/

func (c *EntityRelationshipsClient) GetDestinations(ctx context.Context, sourceTag string, relationshipTypeTag string) ([]string, error) {
	data := EntityDestinationsResponse{}
	apiError := ApiError{}

	url := fmt.Sprintf("/api/v1/catalog/%s/relationships/%s/destinations", sourceTag, relationshipTypeTag)
	response, err := c.Client().Get(url).Receive(&data, &apiError)
	if err != nil {
		return nil, errors.New("could not get entity relationship destinations: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return nil, errors.Join(errors.New("failed getting entity relationship destinations: "), err)
	}

	tags := make([]string, len(data.Destinations))
	for i, dest := range data.Destinations {
		tags[i] = dest.Tag
	}
	return tags, nil
}

/***********************************************************************************************************************
 * PUT /api/v1/catalog/:sourceTag/relationships/:relationshipTypeTag/destinations
 **********************************************************************************************************************/

func (c *EntityRelationshipsClient) SetDestinations(ctx context.Context, sourceTag string, relationshipTypeTag string, destinationTags []string) ([]string, error) {
	data := EntityDestinationsResponse{}
	apiError := ApiError{}

	if destinationTags == nil {
		destinationTags = []string{}
	}

	req := SetEntityDestinationsRequest{Destinations: destinationTags}
	url := fmt.Sprintf("/api/v1/catalog/%s/relationships/%s/destinations", sourceTag, relationshipTypeTag)
	response, err := c.Client().Put(url).BodyJSON(&req).Receive(&data, &apiError)
	if err != nil {
		return nil, errors.New("could not set entity relationship destinations: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return nil, errors.Join(errors.New("failed setting entity relationship destinations: "), err)
	}

	tags := make([]string, len(data.Destinations))
	for i, dest := range data.Destinations {
		tags[i] = dest.Tag
	}
	return tags, nil
}
