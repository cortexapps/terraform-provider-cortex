package cortex

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dghubble/sling"
	"log"
)

type TeamsClientInterface interface {
	Get(ctx context.Context, tag string) (*Team, error)
	List(ctx context.Context, params *TeamListParams) (*TeamsResponse, error)
	Create(ctx context.Context, createReq CreateTeamRequest) (*Team, error)
	Update(ctx context.Context, tag string, updateReq UpdateTeamRequest) (*Team, error)
	Delete(ctx context.Context, tag string) error
	Archive(ctx context.Context, tag string) error
	Unarchive(ctx context.Context, tag string) error
}

type TeamsClient struct {
	client *HttpClient
}

var _ TeamsClientInterface = &TeamsClient{}

func (c *TeamsClient) Client() *sling.Sling {
	return c.client.Client()
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
	CortexTeam        TeamCortexManaged  `json:"cortexTeam,omitempty"`
}

type TeamMetadata struct {
	Description string `json:"description,omitempty"`
	Name        string `json:"name"`
	Summary     string `json:"summary,omitempty"`
}

type TeamMember struct {
	Description string `json:"description,omitempty" tfsdk:"description"`
	Name        string `json:"name" tfsdk:"name"`
	Email       string `json:"email" tfsdk:"email"`
}

type TeamLink struct {
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"`
	Url         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type TeamSlackChannel struct {
	Name                 string `json:"name" tfsdk:"name"`
	NotificationsEnabled bool   `json:"notificationsEnabled" tfsdk:"notifications_enabled"`
}

type TeamIdpGroup struct {
	Group    string               `json:"group"`
	Members  []TeamIdpGroupMember `json:"members"`
	Provider string               `json:"provider"`
}

type TeamCortexManaged struct {
	Members []TeamMember `json:"members"`
}

type TeamIdpGroupMember struct {
	Email string `json:"email,omitempty"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

/***********************************************************************************************************************
 * GET /api/v1/teams/:tag
 **********************************************************************************************************************/

func (c *TeamsClient) Get(ctx context.Context, tag string) (*Team, error) {
	teamResponse := &Team{}
	apiError := &ApiError{}
	response, err := c.Client().Get(Route("teams", tag)).Receive(teamResponse, apiError)
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

	response, err := c.Client().Get(Route("teams", "")).QueryStruct(&params).Receive(teamsResponse, apiError)
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
	TeamTag           string             `json:"teamTag"`
	Type              string             `json:"type"`
	Metadata          TeamMetadata       `json:"metadata"`
	IsArchived        bool               `json:"isArchived,omitempty"`
	AdditionalMembers []TeamMember       `json:"additionalMembers"`
	SlackChannels     []TeamSlackChannel `json:"slackChannels"`
	Links             []TeamLink         `json:"links"`
	CortexTeam        TeamCortexManaged  `json:"cortexTeam,omitempty"`
	IdpGroup          TeamIdpGroup       `json:"idpGroup,omitempty"`
}

func (c *TeamsClient) Create(ctx context.Context, req CreateTeamRequest) (*Team, error) {
	teamResponse := &Team{}
	apiError := &ApiError{}

	response, err := c.Client().Post(Route("teams", "")).BodyJSON(&req).Receive(teamResponse, apiError)
	if err != nil {
		return teamResponse, errors.New("could not create team: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		reqJson, _ := json.Marshal(req)
		log.Printf("Failed creating team: %+v\n\nRequest:\n%+v", err, string(reqJson))
		return teamResponse, err
	}

	return teamResponse, nil
}

/***********************************************************************************************************************
 * PUT /api/v1/teams/:tag
 **********************************************************************************************************************/

type UpdateTeamRequest struct {
	Metadata          TeamMetadata       `json:"metadata"`
	Links             []TeamLink         `json:"links"`
	SlackChannels     []TeamSlackChannel `json:"slackChannels"`
	AdditionalMembers []TeamMember       `json:"additionalMembers"`
}

func (c *TeamsClient) Update(ctx context.Context, tag string, req UpdateTeamRequest) (*Team, error) {
	teamResponse := &Team{}
	apiError := &ApiError{}

	response, err := c.Client().Put(Route("teams", tag)).BodyJSON(&req).Receive(teamResponse, apiError)
	if err != nil {
		return teamResponse, errors.New("could not update team: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		reqJson, _ := json.Marshal(req)
		log.Printf("Failed updating team: %+v\n\nRequest:\n%+v\n%+v", err, string(reqJson), apiError.String())
		return teamResponse, err
	}

	return teamResponse, nil
}

/***********************************************************************************************************************
 * DELETE /api/v1/teams - Delete a team
 **********************************************************************************************************************/

type DeleteTeamResponse struct{}
type DeleteTeamRequest struct {
	Tag string `json:"teamTag" url:"teamTag"`
}

func (c *TeamsClient) Delete(ctx context.Context, tag string) error {
	teamResponse := &ArchiveTeamResponse{}
	apiError := &ApiError{}
	req := DeleteTeamRequest{Tag: tag}

	response, err := c.Client().Delete(Route("teams", "")).QueryStruct(req).Receive(teamResponse, apiError)
	if err != nil {
		return fmt.Errorf("could not delete team %v:\n\n%+v", tag, err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		log.Printf("Could not delete team %v:\n\n%+v", tag, err.Error())
		return err
	}

	return nil
}

/***********************************************************************************************************************
 * PUT /api/v1/teams/:tag/archive - Archive a team
 **********************************************************************************************************************/

type ArchiveTeamResponse struct{}

func (c *TeamsClient) Archive(ctx context.Context, tag string) error {
	teamResponse := &ArchiveTeamResponse{}
	apiError := &ApiError{}

	response, err := c.Client().Put(Route("teams", tag+"/archive")).Receive(teamResponse, apiError)
	if err != nil {
		return fmt.Errorf("could not archive team: %v", err.Error())
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

func (c *TeamsClient) Unarchive(ctx context.Context, tag string) error {
	teamResponse := &UnarchiveTeamResponse{}
	apiError := &ApiError{}

	response, err := c.Client().Put(Route("teams", tag+"/unarchive")).Receive(teamResponse, apiError)
	if err != nil {
		return errors.New("could not unarchive team: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, apiError)
	if err != nil {
		return err
	}

	return nil
}
