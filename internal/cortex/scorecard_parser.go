package cortex

type ScorecardParser struct{}

// YamlToEntity converts YAML into a Scorecard, from the specification.
func (c *ScorecardParser) YamlToEntity(yamlEntity map[string]interface{}) (Scorecard, error) {
	entity := Scorecard{}

	entity.Name = MapFetchToString(yamlEntity, "name")
	entity.Tag = MapFetchToString(yamlEntity, "tag")
	entity.Description = MapFetchToString(yamlEntity, "description")
	entity.Draft = MapFetch(yamlEntity, "draft", false).(bool)

	if yamlEntity["rules"] != nil {
		c.interpolateRules(&entity, yamlEntity["rules"].([]interface{}))
	}
	if yamlEntity["ladder"] != nil {
		c.interpolateLadder(&entity, yamlEntity["ladder"].(map[string]interface{}))
	}
	if yamlEntity["filter"] != nil {
		c.interpolateFilter(&entity, yamlEntity["filter"].(map[string]interface{}))
	}
	if yamlEntity["evaluation"] != nil {
		c.interpolateEvaluation(&entity, yamlEntity["evaluation"].(map[string]interface{}))
	}

	return entity, nil
}

func (c *ScorecardParser) interpolateRules(entity *Scorecard, rules []interface{}) {
	var rs []ScorecardRule
	for _, rule := range rules {
		ruleMap := rule.(map[string]interface{})
		rs = append(rs, ScorecardRule{
			Title:          MapFetchToString(ruleMap, "title"),
			Expression:     MapFetchToString(ruleMap, "expression"),
			Weight:         int64(MapFetch(ruleMap, "weight", 1).(int)),
			Level:          MapFetchToString(ruleMap, "level"),
			Description:    MapFetchToString(ruleMap, "description"),
			FailureMessage: MapFetchToString(ruleMap, "failureMessage"),
		})
	}
	entity.Rules = rs
}

func (c *ScorecardParser) interpolateLadder(entity *Scorecard, ladder map[string]interface{}) {
	entity.Ladder = ScorecardLadder{
		Levels: []ScorecardLevel{},
	}
	if ladder["levels"] != nil {
		c.interpolateLadderLevels(entity, ladder["levels"].([]interface{}))
	}
}

func (c *ScorecardParser) interpolateLadderLevels(entity *Scorecard, levels []interface{}) {
	ls := make([]ScorecardLevel, len(levels))
	for i, level := range levels {
		levelMap := level.(map[string]interface{})
		ls[i] = ScorecardLevel{
			Name:        MapFetchToString(levelMap, "name"),
			Rank:        int64(MapFetch(levelMap, "rank", 1).(int)),
			Description: MapFetchToString(levelMap, "description"),
			Color:       MapFetchToString(levelMap, "color"),
		}
	}
	entity.Ladder.Levels = ls
}

func (c *ScorecardParser) interpolateFilter(entity *Scorecard, filter map[string]interface{}) {
	entity.Filter = ScorecardFilter{
		Category: MapFetchToString(filter, "category"),
		Query:    MapFetchToString(filter, "query"),
	}
}

func (c *ScorecardParser) interpolateEvaluation(entity *Scorecard, evaluation map[string]interface{}) {
	entity.Evaluation = ScorecardEvaluation{
		Window: int64(MapFetch(evaluation, "window", 4).(int)),
	}
}
