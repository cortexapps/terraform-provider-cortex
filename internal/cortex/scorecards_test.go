package cortex_test

import (
	"context"
	"github.com/bigcommerce/terraform-provider-cortex/internal/cortex"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testScorecard = &cortex.Scorecard{
	Tag:         "test-scorecard",
	Name:        "Test Scorecard",
	Description: "This is a test scorecard",
	Draft:       false,
}

func TestGetScorecard(t *testing.T) {
	tag := testScorecard.Tag
	c, teardown, err := setupYamlClient(cortex.Route("scorecards", tag+"/descriptor"), testScorecard, AssertRequestMethod(t, "GET"))
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	res, err := c.Scorecards().Get(context.Background(), tag)
	assert.Nil(t, err, "error retrieving a scorecard")
	assert.Equal(t, testScorecard.Tag, res.Tag)
	assert.Equal(t, testScorecard.Name, res.Name)
	assert.Equal(t, testScorecard.Description, res.Description)
	assert.Equal(t, testScorecard.Draft, res.Draft)
}

func TestDeleteScorecard(t *testing.T) {
	tag := testScorecard.Tag
	c, teardown, err := setupClient(
		cortex.Route("scorecards", tag),
		cortex.DeleteScorecardResponse{},
		AssertRequestMethod(t, "DELETE"),
	)
	assert.Nil(t, err, "could not setup client")
	defer teardown()

	err = c.Scorecards().Delete(context.Background(), tag)
	assert.Nil(t, err, "error deleting a scorecard")
}
