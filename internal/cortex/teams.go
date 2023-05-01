package cortex

import (
	"context"
	"errors"
	"github.com/dghubble/sling"
)

type TeamsClientInterface interface {
	Get(ctx context.Context, teamTag string) (*TeamResponse, error)
	List(ctx context.Context, params *TeamListParams) (*TeamsResponse, error)
	Create(ctx context.Context, createReq CreateTeamRequest) (*TeamResponse, error)
	Update(ctx context.Context, teamTag string, updateReq UpdateTeamRequest) (*TeamResponse, error)
	Delete(ctx context.Context, teamTag string) error
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
	response, err := c.Client().Get(BaseUris["teams"]+teamTag).Receive(teamResponse, apiError)
	if err != nil {
		return teamResponse, errors.New("could not get team: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return teamResponse, errors.Join(errors.New("Failed creating team: "), err)
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
func (c *TeamsClient) List(ctx context.Context, params *TeamListParams) (*TeamsResponse, error) {
	teamsResponse := &TeamsResponse{}
	apiError := &ApiError{}

	response, err := c.Client().Get(BaseUris["teams"]).QueryStruct(params).Receive(teamsResponse, apiError)
	if err != nil {
		return nil, errors.New("could not get teams: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return nil, err
	}

	return teamsResponse, nil
}

/***********************************************************************************************************************
 * POST /api/v1/teams
 **********************************************************************************************************************/

type CreateTeamRequest struct {
	TeamTag string `json:"team_tag"`
}

func (c *TeamsClient) Create(ctx context.Context, req CreateTeamRequest) (*TeamResponse, error) {
	teamResponse := &TeamResponse{}
	apiError := &ApiError{}

	response, err := c.Client().Post(BaseUris["teams"]).BodyJSON(&req).Receive(teamResponse, apiError)
	if err != nil {
		return teamResponse, errors.New("could not create team: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return teamResponse, err
	}

	return teamResponse, nil
}

/***********************************************************************************************************************
 * PUT /api/v1/teams/:tag
 **********************************************************************************************************************/

type UpdateTeamRequest struct {
}

func (c *TeamsClient) Update(ctx context.Context, teamTag string, req UpdateTeamRequest) (*TeamResponse, error) {
	teamResponse := &TeamResponse{}
	apiError := &ApiError{}

	response, err := c.Client().Put(BaseUris["teams"]+teamTag).BodyJSON(&req).Receive(teamResponse, apiError)
	if err != nil {
		return teamResponse, errors.New("could not update team: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return teamResponse, err
	}

	return teamResponse, nil
}

/***********************************************************************************************************************
 * DELETE /api/v1/teams/:tag
 **********************************************************************************************************************/

type DeleteTeamResponse struct{}

func (c *TeamsClient) Delete(ctx context.Context, teamTag string) error {
	teamResponse := &DeleteTeamResponse{}
	apiError := &ApiError{}

	response, err := c.Client().Delete(BaseUris["teams"]+teamTag).Receive(teamResponse, apiError)
	if err != nil {
		return errors.New("could not delete team: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
