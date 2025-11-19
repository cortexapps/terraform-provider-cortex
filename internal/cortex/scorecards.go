package cortex

import (
	"context"
	"errors"
	"fmt"
	"github.com/dghubble/sling"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gopkg.in/yaml.v3"
	"strings"
)

type ScorecardsClientInterface interface {
	Get(ctx context.Context, tag string) (Scorecard, error)
	Upsert(ctx context.Context, scorecard Scorecard) (Scorecard, error)
	Delete(ctx context.Context, tag string) error
}

type ScorecardsClient struct {
	client *HttpClient
	parser *ScorecardParser
}

var _ ScorecardsClientInterface = &ScorecardsClient{}

func (c *ScorecardsClient) Client() *sling.Sling {
	return c.client.Client()
}

func (c *ScorecardsClient) YamlClient() *sling.Sling {
	return c.client.YamlClient()
}

/***********************************************************************************************************************
 * Types
 **********************************************************************************************************************/

// Scorecard is the nested response object that is typically returned from the scorecards endpoints.
type Scorecard struct {
	Tag         string              `json:"tag" yaml:"tag"`
	Name        string              `json:"name" yaml:"name"`
	Description string              `json:"description,omitempty" yaml:"description,omitempty"`
	Draft       bool                `json:"draft,omitempty" yaml:"draft,omitempty"`
	Rules       []ScorecardRule     `json:"rules,omitempty" yaml:"rules,omitempty"`
	Ladder      ScorecardLadder     `json:"ladder,omitempty" yaml:"ladder,omitempty"`
	Filter      ScorecardFilter     `json:"filter,omitempty" yaml:"filter,omitempty"`
	Evaluation  ScorecardEvaluation `json:"evaluation,omitempty" yaml:"evaluation,omitempty"`
}

func (s *Scorecard) ToYaml() (string, error) {
	// The API requires submitting the request as YAML, so we need to marshal it first.
	bytes, err := yaml.Marshal(s)
	if err != nil {
		return "", errors.New("could not marshal yaml: " + err.Error())
	}
	return string(bytes), nil
}

func (s *Scorecard) ToYamlStringReader() (*strings.Reader, error) {
	yamlString, err := s.ToYaml()
	if err != nil {
		return nil, err
	}
	return strings.NewReader(yamlString), nil
}

type ScorecardLadder struct {
	Levels []ScorecardLevel `json:"levels" yaml:"levels"`
}

func (s *ScorecardLadder) Enabled() bool {
	return len(s.Levels) > 0
}

type ScorecardLevel struct {
	Name        string `json:"name" yaml:"name"`
	Rank        int64  `json:"rank" yaml:"rank"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Color       string `json:"color" yaml:"color"`
}

func (s *ScorecardLevel) Enabled() bool {
	return s.Name != "" && s.Color != "" && s.Rank > 0
}

type ScorecardRule struct {
	Title          string `json:"title" yaml:"title"`
	Expression     string `json:"expression" yaml:"expression"`
	Weight         int64  `json:"weight" yaml:"weight"`
	Level          string `json:"level" yaml:"level"`
	Description    string `json:"description,omitempty" yaml:"description,omitempty"`
	FailureMessage string `json:"failure_message,omitempty" yaml:"failureMessage,omitempty"`
}

type ScorecardFilter struct {
	Kind   string                 `json:"kind,omitempty" yaml:"kind,omitempty"`
	Types  *ScorecardFilterTypes  `json:"types,omitempty" yaml:"types,omitempty"`
	Groups *ScorecardFilterGroups `json:"groups,omitempty" yaml:"groups,omitempty"`
	Query  string                 `json:"query,omitempty" yaml:"query,omitempty"`
}

type ScorecardFilterTypes struct {
	Include []string `json:"include,omitempty" yaml:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty" yaml:"exclude,omitempty"`
}

type ScorecardFilterGroups struct {
	Include []string `json:"include,omitempty" yaml:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty" yaml:"exclude,omitempty"`
}

func (s *ScorecardFilter) Enabled() bool {
	return s.Types != nil || s.Groups != nil || s.Query != ""
}

type ScorecardEvaluation struct {
	Window int64 `json:"window,omitempty" yaml:"window,omitempty"`
}

func (s *ScorecardEvaluation) Enabled() bool {
	return s.Window > 0
}

/***********************************************************************************************************************
 * GET /api/v1/scorecards/:tag
 **********************************************************************************************************************/

// GetScorecardResponse is the generic root response object for scorecards.
type GetScorecardResponse struct {
	Scorecard Scorecard `json:"scorecard"`
}

func (c *ScorecardsClient) Get(ctx context.Context, tag string) (Scorecard, error) {
	scorecardDescriptorResponse := map[string]interface{}{}
	apiError := ApiError{}

	uri := Route("scorecards", tag+"/descriptor")
	cl := c.YamlClient().Get(uri)
	response, err := cl.Receive(scorecardDescriptorResponse, &apiError)
	if err != nil {
		return Scorecard{}, errors.Join(fmt.Errorf("failed getting scorecard descriptor for %s from %s", tag, uri), err)
	}
	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return Scorecard{}, errors.Join(fmt.Errorf("failed handling response status for %s from %s", tag, uri), err)
	}

	tflog.Debug(ctx, fmt.Sprintf("body: %+v", scorecardDescriptorResponse))

	return c.parser.YamlToEntity(scorecardDescriptorResponse)
}

/***********************************************************************************************************************
 * POST /api/v1/scorecards/descriptor
 **********************************************************************************************************************/

type UpsertScorecardResponse struct {
	Scorecard Scorecard `json:"scorecard" yaml:"scorecard"`
}

func (c *ScorecardsClient) Upsert(ctx context.Context, scorecard Scorecard) (Scorecard, error) {
	upsertResponse := UpsertScorecardResponse{
		Scorecard: Scorecard{},
	}
	apiError := ApiError{}

	// The API requires submitting the request as YAML, so we need to marshal it first.
	body, err := scorecard.ToYamlStringReader()
	if err != nil {
		return upsertResponse.Scorecard, errors.New("could not marshal yaml: " + err.Error())
	}

	tflog.Debug(ctx, fmt.Sprintf("CREATE body: %+v", body))
	response, err := c.Client().
		Set("Content-Type", "application/yaml;charset=UTF-8").
		Set("Accept", "application/json").
		Post(Route("scorecards", "descriptor")).
		Body(body).
		Receive(&upsertResponse, &apiError)
	if err != nil {
		return upsertResponse.Scorecard, errors.New("could not upsert scorecard: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return upsertResponse.Scorecard, err
	}

	// re-fetch the scorecard, since it's not fully returned here
	return c.Get(ctx, scorecard.Tag)
}

/***********************************************************************************************************************
 * DELETE /api/v1/scorecards/:tag - Delete a scorecard
 **********************************************************************************************************************/

type DeleteScorecardResponse struct{}

func (c *ScorecardsClient) Delete(ctx context.Context, tag string) error {
	scorecardResponse := DeleteScorecardResponse{}
	apiError := ApiError{}

	response, err := c.Client().Delete(Route("scorecards", tag)).Receive(&scorecardResponse, &apiError)
	if err != nil {
		return errors.New("could not delete scorecard: " + err.Error())
	}

	err = c.client.handleResponseStatus(response, &apiError)
	if err != nil {
		return err
	}

	return nil
}
