package detector

import (
	"context"
	"fmt"
	"strings"

	"gender-detection/strategies"
	"gender-detection/types"
)

type HybridDetector struct {
	config     *types.Config
	strategies []types.Strategy
}

func NewHybridDetector(config *types.Config) (*HybridDetector, error) {
	detector := &HybridDetector{
		config:     config,
		strategies: make([]types.Strategy, 0),
	}

	if config.EnableDictionary {
		dict, err := strategies.NewDictionaryStrategy(config.DictionaryPath)
		if err != nil {
			return nil, fmt.Errorf("failed to init dictionary: %w", err)
		}
		detector.strategies = append(detector.strategies, dict)
	}

	if config.EnableNaiveBayes {
		nb, err := strategies.NewNaiveBayesStrategy(config.ModelPath)
		if err != nil {
			return nil, fmt.Errorf("failed to init naive bayes: %w", err)
		}
		detector.strategies = append(detector.strategies, nb)
	}

	if config.EnableRules {
		rules, err := strategies.NewRulesStrategy(config.PatternsPath)
		if err != nil {
			return nil, fmt.Errorf("failed to init rules: %w", err)
		}
		detector.strategies = append(detector.strategies, rules)
	}

	return detector, nil
}

func (d *HybridDetector) Detect(ctx context.Context, firstName, lastName string) (*types.Result, error) {
	firstName = strings.TrimSpace(strings.ToLower(firstName))
	lastName = strings.TrimSpace(strings.ToLower(lastName))

	if firstName == "" {
		return &types.Result{Gender: types.Unknown, Confidence: 0.0, Strategy: "none"}, nil
	}

	for _, strategy := range d.strategies {
		result, err := strategy.Detect(ctx, firstName, lastName)
		if err != nil {
			continue
		}

		if result != nil && result.Confidence >= d.config.MinConfidence {
			result.Strategy = strategy.Name()
			return result, nil
		}
	}

	if d.config.FallbackToUnknown {
		return &types.Result{Gender: types.Unknown, Confidence: 0.0, Strategy: "fallback"}, nil
	}

	return nil, fmt.Errorf("no confident detection found")
}

func (d *HybridDetector) DetectWithStrategies(ctx context.Context, firstName, lastName string) ([]*types.Result, error) {
	firstName = strings.TrimSpace(strings.ToLower(firstName))
	lastName = strings.TrimSpace(strings.ToLower(lastName))

	results := make([]*types.Result, 0, len(d.strategies))

	for _, strategy := range d.strategies {
		result, err := strategy.Detect(ctx, firstName, lastName)
		if err != nil {
			continue
		}
		if result != nil {
			result.Strategy = strategy.Name()
			results = append(results, result)
		}
	}

	return results, nil
}
