package cortex

import (
	"context"
	"errors"
	"github.com/dghubble/sling"
)

type TeamsClientInterface interface {
	Get(ctx context.Context, teamTag string) (*Team, error)
	List(ctx context.Context, params *TeamListParams) (*TeamsResponse, error)
	Create(ctx context.Context, createReq CreateTeamRequest) (*Team, error)
	Update(ctx context.Context, teamTag string, updateReq UpdateTeamRequest) (*Team, error)
	Archive(ctx context.Context, teamTag string) error
	Unarchive(ctx context.Context, teamTag string) error
}

type TeamsClient struct {
	client *HttpClient
}

var _ TeamsClientInterface = &TeamsClient{}

func (c *TeamsClient) Client() *sling.Sling {
	return c.client.client
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// Team is the response from the GET /v1/teams/{team_tag} endpoint.
type Team struct {
	AdditionalMembers []TeamMember       `json:"additionalMembers,omitempty"`
	IsArchived        bool               `json:"isArchived,omitempty"`
	Metadata          TeamMetadata       `json:"metadata"`
	SlackChannels     []TeamSlackChannel `json:"slackChannels,omitempty"`
	Links             []TeamLink         `json:"links,omitempty"`
	TeamTag           string             `json:"teamTag"`
}

type TeamMetadata struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Summary     string `json:"summary"`
}

type TeamMember struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Email       string `json:"email"`
}

type TeamLink struct {
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"`
	Url         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type TeamSlackChannel struct {
	Name                 string `json:"name"`
	NotificationsEnabled bool   `json:"notificationsEnabled"`
}

type TeamIdpGroup struct {
	Group    string               `json:"group"`
	Members  []TeamIdpGroupMember `json:"members"`
	Provider string               `json:"provider"`
}

type TeamIdpGroupMember struct {
	Email string `json:"email,omitempty"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

/***********************************************************************************************************************
 * GET /api/v1/teams/:tag
 **********************************************************************************************************************/

func (c *TeamsClient) Get(ctx context.Context, teamTag string) (*Team, error) {
	teamResponse := &Team{}
	apiError := &ApiError{}
	response, err := c.Client().Get(BaseUris["teams"]+teamTag).Receive(teamResponse, apiError)
	if err != nil {
		return teamResponse, errors.New("could not get team: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return teamResponse, errors.Join(errors.New("Failed getting team: "), err)
	}

	return teamResponse, nil
}

/***********************************************************************************************************************
 * GET /api/v1/teams
 **********************************************************************************************************************/

// TeamListParams are the query parameters for the GET /v1/teams endpoint.
type TeamListParams struct {
}

// TeamsResponse is the response from the GET /v1/teams endpoint.
type TeamsResponse struct {
	Teams []Team `json:"teams"`
}

// List retrieves a list of teams based on a team query.
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
	AdditionalMembers []TeamMember       `json:"additionalMembers,omitempty"`
	IsArchived        bool               `json:"isArchived,omitempty"`
	Metadata          TeamMetadata       `json:"metadata"`
	SlackChannels     []TeamSlackChannel `json:"slackChannels,omitempty"`
	Links             []TeamLink         `json:"links,omitempty"`
	TeamTag           string             `json:"teamTag"`
	IdpGroup          TeamIdpGroup       `json:"idpGroup"`
}

func (c *TeamsClient) Create(ctx context.Context, req CreateTeamRequest) (*Team, error) {
	teamResponse := &Team{}
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
	AdditionalMembers []TeamMember       `json:"additionalMembers,omitempty"`
	Metadata          TeamMetadata       `json:"metadata"`
	SlackChannels     []TeamSlackChannel `json:"slackChannels,omitempty"`
	Links             []TeamLink         `json:"links,omitempty"`
}

func (c *TeamsClient) Update(ctx context.Context, teamTag string, req UpdateTeamRequest) (*Team, error) {
	teamResponse := &Team{}
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
 * DELETE /api/v1/teams/:tag - Archive a team
 **********************************************************************************************************************/

type ArchiveTeamResponse struct{}

func (c *TeamsClient) Archive(ctx context.Context, teamTag string) error {
	teamResponse := &ArchiveTeamResponse{}
	apiError := &ApiError{}

	response, err := c.Client().Delete(BaseUris["teams"]+teamTag).Receive(teamResponse, apiError)
	if err != nil {
		return errors.New("could not archive team: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return err
	}

	return nil
}

/***********************************************************************************************************************
 * PUT /api/v1/teams/:tag/unarchive - Un-archive a team
 **********************************************************************************************************************/

type UnarchiveTeamResponse struct{}

func (c *TeamsClient) Unarchive(ctx context.Context, teamTag string) error {
	teamResponse := &UnarchiveTeamResponse{}
	apiError := &ApiError{}

	response, err := c.Client().Put(BaseUris["teams"]+teamTag+"/unarchive").Receive(teamResponse, apiError)
	if err != nil {
		return errors.New("could not unarchive team: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
