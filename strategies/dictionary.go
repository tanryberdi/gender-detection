package strategies

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gender-detection/types"
)

type DictionaryStrategy struct {
	names map[string]types.Gender
}

type nameEntry struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
}

func NewDictionaryStrategy(path string) (*DictionaryStrategy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dictionary: %w", err)
	}

	var entries []nameEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse dictionary: %w", err)
	}

	names := make(map[string]types.Gender)
	for _, entry := range entries {
		key := strings.ToLower(strings.TrimSpace(entry.Name))
		names[key] = types.Gender(strings.ToLower(entry.Gender))
	}

	return &DictionaryStrategy{names: names}, nil
}

func (d *DictionaryStrategy) Detect(ctx context.Context, firstName, lastName string) (*types.Result, error) {
	firstName = strings.ToLower(strings.TrimSpace(firstName))

	if gender, found := d.names[firstName]; found {
		return &types.Result{
			Gender:     gender,
			Confidence: 1.0,
			Metadata:   map[string]interface{}{"source": "exact_match"},
		}, nil
	}

	return nil, nil
}

func (d *DictionaryStrategy) Name() string {
	return "dictionary"
}

func (d *DictionaryStrategy) MinConfidence() float64 {
	return 0.95
}
