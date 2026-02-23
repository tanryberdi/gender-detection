package training

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gender-detection/strategies"
)

func TrainNaiveBayes(datasetPath, outputPath string) error {
	file, err := os.Open(datasetPath)
	if err != nil {
		return fmt.Errorf("failed to open dataset: %w", err)
	}
	defer file.Close() //nolint:errcheck

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}

	model := &strategies.NaiveBayesModel{
		MaleFeatures:   make(map[string]int),
		FemaleFeatures: make(map[string]int),
		Smoothing:      1.0,
	}

	for i, record := range records {
		if i == 0 {
			continue
		}
		if len(record) < 2 {
			continue
		}

		name := strings.ToLower(strings.TrimSpace(record[0]))
		gender := strings.ToLower(strings.TrimSpace(record[1]))

		if gender != "male" && gender != "female" {
			continue
		}

		features := extractFeatures(name)

		if gender == "male" {
			model.MaleCount++
			for _, feature := range features {
				model.MaleFeatures[feature]++
			}
		} else {
			model.FemaleCount++
			for _, feature := range features {
				model.FemaleFeatures[feature]++
			}
		}
	}

	data, err := json.MarshalIndent(model, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal model: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write model: %w", err)
	}

	fmt.Printf("Model trained: %d male, %d female examples\n", model.MaleCount, model.FemaleCount)
	return nil
}

func extractFeatures(name string) []string {
	features := []string{}

	if len(name) >= 2 {
		features = append(features, "first2:"+name[:2])
	}
	if len(name) >= 3 {
		features = append(features, "first3:"+name[:3])
		features = append(features, "last3:"+name[len(name)-3:])
	}
	if len(name) >= 2 {
		features = append(features, "last2:"+name[len(name)-2:])
	}
	if len(name) >= 1 {
		features = append(features, "last1:"+name[len(name)-1:])
	}

	for i := 0; i < len(name)-1; i++ {
		features = append(features, "bigram:"+name[i:i+2])
	}

	features = append(features, fmt.Sprintf("length:%d", len(name)))

	return features
}
