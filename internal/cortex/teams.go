package cortex

import (
	"context"
	"errors"
	"github.com/dghubble/sling"
)

type TeamsClientInterface interface {
	Get(ctx context.Context, teamTag string) (*TeamResponse, error)
}

type TeamsClient struct {
	client *HttpClient
}

var _ TeamsClientInterface = &TeamsClient{}

func (c *TeamsClient) Client() *sling.Sling {
	return c.client.client
}

/***********************************************************************************************************************
 * GET /api/v1/teams/:tag
 **********************************************************************************************************************/

// TeamResponse is the response from the GET /v1/teams/{team_tag} endpoint
type TeamResponse struct {
	TeamTag string `json:"team_tag"`
}

func (c *TeamsClient) Get(ctx context.Context, teamTag string) (*TeamResponse, error) {
	teamResponse := &TeamResponse{}
	apiError := &ApiError{}
	response, err := c.Client().Get("v1/teams/"+teamTag).Receive(teamResponse, apiError)
	if err != nil {
		return teamResponse, errors.New("could not get team: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return teamResponse, err
	}

	return teamResponse, nil
}

/***********************************************************************************************************************
 * GET /api/v1/teams
 **********************************************************************************************************************/

// TeamListParams are the query parameters for the GET /v1/teams endpoint
type TeamListParams struct {
}

// TeamsResponse is the response from the GET /v1/teams endpoint
type TeamsResponse struct {
	Teams []TeamResponse `json:"teams"`
}

// List retrieves a list of teams based on a team query
func (c *TeamsClient) List(ctx context.Context, req *TeamListParams) (*TeamsResponse, error) {
	teamsResponse := &TeamsResponse{}
	apiError := &ApiError{}

	response, err := c.Client().Get("/v1/teams").QueryStruct(req).Receive(teamsResponse, apiError)
	if err != nil {
		return nil, errors.New("could not get teams: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return nil, err
	}

	return teamsResponse, nil
}
