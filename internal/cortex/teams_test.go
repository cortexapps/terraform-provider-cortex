package cortex_test

import (
	"context"
	"github.com/cortexapps/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testTeamResponse = &cortex.Team{
	TeamTag:    "test-team",
	IsArchived: false,
	Metadata: cortex.TeamMetadata{
		Name:        "Test Team",
		Description: "A test team",
		Summary:     "A short summary about the team",
	},
	SlackChannels: []cortex.TeamSlackChannel{
		{
			Name:                 "#test-dev",
			NotificationsEnabled: false,
		},
		{
			Name:                 "#test-alerts",
			NotificationsEnabled: true,
		},
	},
	Links: []cortex.TeamLink{
		{
			Name: "Test Link",
			Url:  "https://cortexapp.com",
		},
	},
}

func TestGetTeam(t *testing.T) {
	testTeamTag := "test-team"
	resp := testTeamResponse
	c, teardown, err := setupClient(cortex.Route("teams", testTeamTag), resp, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.Teams().Get(context.Background(), testTeamTag)
	assert.Nil(t, err, "error retrieving a team")
	assert.Equal(t, resp, res)
}

func TestListTeams(t *testing.T) {
	firstTeamTag := "test-team"
	resp := &cortex.TeamsResponse{
		Teams: []cortex.Team{
			*testTeamResponse,
		},
	}
	c, teardown, err := setupClient(cortex.Route("teams", ""), resp, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	var queryParams cortex.TeamListParams
	res, err := c.Teams().List(context.Background(), &queryParams)
	assert.Nil(t, err, "error retrieving a team")
	assert.NotEmpty(t, res.Teams, "returned no teams")
	assert.Equal(t, res.Teams[0].TeamTag, firstTeamTag)
}

func TestCreateTeam(t *testing.T) {
	teamTag := "test-team"
	req := cortex.CreateTeamRequest{
		TeamTag: teamTag,
		AdditionalMembers: []cortex.TeamMember{
			{
				Name:        "Test Member",
				Description: "A test member",
				Email:       "test@cortex.io",
			},
		},
		Metadata: cortex.TeamMetadata{
			Name:        "Test Team",
			Description: "A test team",
			Summary:     "A short summary about the team",
		},
		SlackChannels: []cortex.TeamSlackChannel{
			{
				Name:                 "#test-dev",
				NotificationsEnabled: false,
			},
			{
				Name:                 "#test-alerts",
				NotificationsEnabled: true,
			},
		},
		Links: []cortex.TeamLink{
			{
				Name: "Test Link",
				Url:  "https://cortexapp.com",
			},
		},
	}
	c, teardown, err := setupClient(
		cortex.Route("teams", ""),
		testTeamResponse,
		AssertRequestMethod(t, "POST"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.Teams().Create(context.Background(), req)
	assert.Nil(t, err, "error creating a team")
	assert.Equal(t, res.TeamTag, teamTag)
}

func TestUpdateTeam(t *testing.T) {
	req := cortex.UpdateTeamRequest{
		AdditionalMembers: []cortex.TeamMember{
			{
				Name:        "Test Member",
				Description: "A test member",
				Email:       "test@cortex.io",
			},
		},
		Metadata: cortex.TeamMetadata{
			Name:        "Test Team",
			Description: "A test team",
			Summary:     "A short summary about the team",
		},
		SlackChannels: []cortex.TeamSlackChannel{
			{
				Name:                 "#test-dev",
				NotificationsEnabled: false,
			},
			{
				Name:                 "#test-alerts",
				NotificationsEnabled: true,
			},
		},
		Links: []cortex.TeamLink{
			{
				Name: "Test Link",
				Url:  "https://cortexapp.com",
			},
		},
	}
	teamTag := "test-team"

	c, teardown, err := setupClient(
		cortex.Route("teams", teamTag),
		testTeamResponse,
		AssertRequestMethod(t, "PUT"),
		AssertRequestBody(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.Teams().Update(context.Background(), teamTag, req)
	assert.Nil(t, err, "error updating a team")
	assert.Equal(t, res.TeamTag, teamTag)
}

func TestDeleteTeam(t *testing.T) {
	tag := "test-team"

	c, teardown, err := setupClient(
		cortex.Route("teams", ""),
		cortex.DeleteTeamResponse{},
		AssertRequestMethod(t, "DELETE"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.Teams().Delete(context.Background(), tag)
	assert.Nil(t, err, "error deleting a team")
}

func TestArchiveTeam(t *testing.T) {
	teamTag := "test-team"

	c, teardown, err := setupClient(
		cortex.Route("teams", teamTag+"/archive"),
		cortex.ArchiveTeamResponse{},
		AssertRequestMethod(t, "PUT"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.Teams().Archive(context.Background(), teamTag)
	assert.Nil(t, err, "error archiving a team")
}

func TestUnarchiveTeam(t *testing.T) {
	teamTag := "test-team"

	c, teardown, err := setupClient(
		cortex.Route("teams", teamTag+"/unarchive"),
		cortex.UnarchiveTeamResponse{},
		AssertRequestMethod(t, "PUT"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.Teams().Unarchive(context.Background(), teamTag)
	assert.Nil(t, err, "error unarchiving a team")
}
