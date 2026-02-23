package strategies

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gender-detection/types"
)

type PatternRule struct {
	Pattern    string  `json:"pattern"`
	Type       string  `json:"type"`
	Gender     string  `json:"gender"`
	Confidence float64 `json:"confidence"`
}

type PatternConfig struct {
	Suffixes []PatternRule `json:"suffixes"`
	Prefixes []PatternRule `json:"prefixes"`
	Contains []PatternRule `json:"contains"`
}

type RulesStrategy struct {
	config *PatternConfig
}

func NewRulesStrategy(path string) (*RulesStrategy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read patterns: %w", err)
	}

	var config PatternConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse patterns: %w", err)
	}

	return &RulesStrategy{config: &config}, nil
}

func (r *RulesStrategy) Detect(ctx context.Context, firstName, lastName string) (*types.Result, error) {
	firstName = strings.ToLower(strings.TrimSpace(firstName))
	if firstName == "" {
		return nil, nil
	}

	var bestMatch *PatternRule
	var matchType string

	for _, rule := range r.config.Suffixes {
		if strings.HasSuffix(firstName, strings.ToLower(rule.Pattern)) {
			if bestMatch == nil || rule.Confidence > bestMatch.Confidence {
				bestMatch = &rule
				matchType = "suffix"
			}
		}
	}

	for _, rule := range r.config.Prefixes {
		if strings.HasPrefix(firstName, strings.ToLower(rule.Pattern)) {
			if bestMatch == nil || rule.Confidence > bestMatch.Confidence {
				bestMatch = &rule
				matchType = "prefix"
			}
		}
	}

	for _, rule := range r.config.Contains {
		if strings.Contains(firstName, strings.ToLower(rule.Pattern)) {
			if bestMatch == nil || rule.Confidence > bestMatch.Confidence {
				bestMatch = &rule
				matchType = "contains"
			}
		}
	}

	if bestMatch != nil {
		return &types.Result{
			Gender:     types.Gender(strings.ToLower(bestMatch.Gender)),
			Confidence: bestMatch.Confidence,
			Metadata: map[string]interface{}{
				"pattern":    bestMatch.Pattern,
				"match_type": matchType,
			},
		}, nil
	}

	return nil, nil
}

func (r *RulesStrategy) Name() string {
	return "rules"
}

func (r *RulesStrategy) MinConfidence() float64 {
	return 0.60
}
