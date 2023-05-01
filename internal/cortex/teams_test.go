package cortex

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testTeamResponse = &Team{
	TeamTag:    "test-team",
	IsArchived: false,
	Metadata: TeamMetadata{
		Name:        "Test Team",
		Description: "A test team",
		Summary:     "A short summary about the team",
	},
	SlackChannels: []TeamSlackChannel{
		{
			Name:                 "#test-dev",
			NotificationsEnabled: false,
		},
		{
			Name:                 "#test-alerts",
			NotificationsEnabled: true,
		},
	},
	Links: []TeamLink{
		{
			Name: "Test Link",
			Url:  "https://cortexapp.com",
		},
	},
}

func TestGetTeam(t *testing.T) {
	testTeamTag := "test-team"
	resp := testTeamResponse
	c, teardown, err := setupClient(BaseUris["teams"]+testTeamTag, resp, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.Teams().Get(context.Background(), testTeamTag)
	assert.Nil(t, err, "error retrieving a team")
	assert.Equal(t, resp, res)
}

func TestListTeams(t *testing.T) {
	firstTeamTag := "test-team"
	resp := &TeamsResponse{
		Teams: []Team{
			*testTeamResponse,
		},
	}
	c, teardown, err := setupClient(BaseUris["teams"], resp, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	var queryParams TeamListParams
	res, err := c.Teams().List(context.Background(), &queryParams)
	assert.Nil(t, err, "error retrieving a team")
	assert.NotEmpty(t, res.Teams, "returned no teams")
	assert.Equal(t, res.Teams[0].TeamTag, firstTeamTag)
}

func TestCreateTeam(t *testing.T) {
	teamTag := "test-team"
	req := CreateTeamRequest{
		TeamTag: teamTag,
		AdditionalMembers: []TeamMember{
			{
				Name:        "Test Member",
				Description: "A test member",
				Email:       "test@cortex.io",
			},
		},
		Metadata: TeamMetadata{
			Name:        "Test Team",
			Description: "A test team",
			Summary:     "A short summary about the team",
		},
		SlackChannels: []TeamSlackChannel{
			{
				Name:                 "#test-dev",
				NotificationsEnabled: false,
			},
			{
				Name:                 "#test-alerts",
				NotificationsEnabled: true,
			},
		},
		Links: []TeamLink{
			{
				Name: "Test Link",
				Url:  "https://cortexapp.com",
			},
		},
	}
	c, teardown, err := setupClient(
		BaseUris["teams"],
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
	req := UpdateTeamRequest{
		AdditionalMembers: []TeamMember{
			{
				Name:        "Test Member",
				Description: "A test member",
				Email:       "test@cortex.io",
			},
		},
		Metadata: TeamMetadata{
			Name:        "Test Team",
			Description: "A test team",
			Summary:     "A short summary about the team",
		},
		SlackChannels: []TeamSlackChannel{
			{
				Name:                 "#test-dev",
				NotificationsEnabled: false,
			},
			{
				Name:                 "#test-alerts",
				NotificationsEnabled: true,
			},
		},
		Links: []TeamLink{
			{
				Name: "Test Link",
				Url:  "https://cortexapp.com",
			},
		},
	}
	teamTag := "test-team"

	c, teardown, err := setupClient(
		BaseUris["teams"],
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

func TestArchiveTeam(t *testing.T) {
	teamTag := "test-team"

	c, teardown, err := setupClient(
		BaseUris["teams"],
		ArchiveTeamResponse{},
		AssertRequestMethod(t, "DELETE"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.Teams().Archive(context.Background(), teamTag)
	assert.Nil(t, err, "error archiving a team")
}

func TestUnarchiveTeam(t *testing.T) {
	teamTag := "test-team"

	c, teardown, err := setupClient(
		BaseUris["teams"]+teamTag+"/unarchive",
		UnarchiveTeamResponse{},
		AssertRequestMethod(t, "PUT"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.Teams().Unarchive(context.Background(), teamTag)
	assert.Nil(t, err, "error unarchiving a team")
}
