package cortex

import (
	"context"
	"testing"
)

func TestGetTeam(t *testing.T) {
	resp := &TeamResponse{}
	testTeamTag := "test-team-id"
	c, teardown, err := setupClient(BaseUris["teams"]+testTeamTag, resp, AssertRequestMethod(t, "GET"))
	if err != nil {
		t.Fatalf("could not setup client: %v", err)
	}
	defer teardown()

	res, err := c.Teams().Get(context.Background(), testTeamTag)
	if err != nil {
		t.Fatalf("error retrieving a team: %v", err)
	}
	if resp.TeamTag != res.TeamTag {
		t.Fatalf("returned TeamTag did not match")
	}
}

func TestListTeams(t *testing.T) {
	resp := &TeamsResponse{
		Teams: []TeamResponse{
			{
				TeamTag: "test-team",
			},
		},
	}
	c, teardown, err := setupClient(BaseUris["teams"], resp, AssertRequestMethod(t, "GET"))
	if err != nil {
		t.Fatalf("could not setup client: %v", err)
	}
	defer teardown()

	var queryParams TeamListParams
	res, err := c.Teams().List(context.Background(), &queryParams)
	if err != nil {
		t.Fatalf("error retrieving a team: %v", err)
	}
	if len(res.Teams) <= 0 {
		t.Fatalf("returned no teams: %s", res)
	}
	if res.Teams[0].TeamTag != "test-team" {
		t.Fatalf("returned TeamTag did not match: %s", res.Teams[0].TeamTag)
	}
}

func TestCreateTeam(t *testing.T) {
	req := CreateTeamRequest{
		TeamTag: "fake-team",
	}
	c, teardown, err := setupClient(
		BaseUris["teams"],
		TeamResponse{
			TeamTag: "fake-team",
		},
		AssertRequestMethod(t, "POST"),
		AssertRequestBody(t, req),
	)
	if err != nil {
		t.Fatalf("could not setup client: %v", err)
	}
	defer teardown()

	res, err := c.Teams().Create(context.Background(), req)
	if err != nil {
		t.Fatalf("error creating a team: %v", err)
	}

	if res.TeamTag != req.TeamTag {
		t.Fatalf("returned TeamTag did not match: %s != %s", res, req.TeamTag)
	}
}

func TestUpdateTeam(t *testing.T) {
	req := UpdateTeamRequest{}
	teamTag := "test-team"

	c, teardown, err := setupClient(
		BaseUris["teams"],
		TeamResponse{
			TeamTag: teamTag,
		},
		AssertRequestMethod(t, "PUT"),
		AssertRequestBody(t, req),
	)
	if err != nil {
		t.Fatalf("could not setup client: %v", err)
	}
	defer teardown()

	res, err := c.Teams().Update(context.Background(), teamTag, req)
	if err != nil {
		t.Fatalf("error updating a team: %v", err)
	}

	if res.TeamTag != teamTag {
		t.Fatalf("returned TeamTag did not match: %s != %s", res.TeamTag, teamTag)
	}
}

func TestDeleteTeam(t *testing.T) {
	teamTag := "test-team"

	c, teardown, err := setupClient(
		BaseUris["teams"],
		DeleteTeamResponse{},
		AssertRequestMethod(t, "DELETE"),
	)
	if err != nil {
		t.Fatalf("could not setup client: %v", err)
	}
	defer teardown()

	err = c.Teams().Delete(context.Background(), teamTag)
	if err != nil {
		t.Fatalf("error deleting a team: %v", err)
	}
}
