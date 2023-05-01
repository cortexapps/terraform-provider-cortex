package cortex

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testScorecard = &Scorecard{
	Tag: "test-scorecard",
}

func TestGetScorecard(t *testing.T) {
	testTag := "test-scorecard"
	resp := GetScorecardResponse{
		Scorecard: testScorecard,
	}
	c, teardown, err := setupClient(BaseUris["scorecards"]+testTag, resp, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.Scorecards().Get(context.Background(), testTag)
	assert.Nil(t, err, "error retrieving a scorecard")
	assert.Equal(t, testScorecard, res)
}

func TestListScorecards(t *testing.T) {
	firstTag := "test-scorecard"
	resp := &ScorecardsResponse{
		Scorecards: []Scorecard{
			*testScorecard,
		},
	}
	c, teardown, err := setupClient(BaseUris["scorecards"], resp, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	var queryParams ScorecardListParams
	res, err := c.Scorecards().List(context.Background(), &queryParams)
	assert.Nil(t, err, "error retrieving scorecards")
	assert.NotEmpty(t, res.Scorecards, "returned no scorecards")
	assert.Equal(t, res.Scorecards[0].Tag, firstTag)
}

func TestUpsertScorecard(t *testing.T) {
	tag := "test-scorecard"
	req := UpsertScorecardRequest{
		Tag: tag,
	}
	upsertScorecardResponse := UpsertScorecardResponse{
		Scorecard: testScorecard,
	}
	c, teardown, err := setupClient(
		BaseUris["scorecards"],
		upsertScorecardResponse,
		AssertRequestMethod(t, "POST"),
		AssertRequestBodyYaml(t, req),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.Scorecards().Upsert(context.Background(), req)
	assert.Nil(t, err, "error creating a scorecard")
	assert.Equal(t, tag, res.Tag)
}

func TestDeleteScorecard(t *testing.T) {
	tag := "test-scorecard"

	c, teardown, err := setupClient(
		BaseUris["scorecards"],
		ArchiveTeamResponse{},
		AssertRequestMethod(t, "DELETE"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.Scorecards().Delete(context.Background(), tag)
	assert.Nil(t, err, "error deleting a scorecard")
}
