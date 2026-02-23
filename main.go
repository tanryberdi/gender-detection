package main

import (
	"context"
	"fmt"
	"log"

	"gender-detection/detector"
	"gender-detection/training"
	"gender-detection/types"
)

func main() {
	if err := training.TrainNaiveBayes("training/dataset.csv", "data/model.json"); err != nil {
		log.Fatal(err)
	}

	config := &types.Config{
		MinConfidence:     0.70,
		EnableDictionary:  true,
		EnableNaiveBayes:  true,
		EnableRules:       true,
		DictionaryPath:    "data/names.json",
		ModelPath:         "data/model.json",
		PatternsPath:      "data/patterns.json",
		FallbackToUnknown: true,
	}

	det, err := detector.NewHybridDetector(config)
	if err != nil {
		log.Fatal(err)
	}

	testNames := []struct{ first, last string }{
		{"Ayşegül", "Yılmaz"},
		{"Mehmet", "Demir"},
		{"Gulnaz", ""},
		{"Nurбek", ""},
		{"Unknown", "Person"},
	}

	for _, name := range testNames {
		result, err := det.Detect(context.Background(), name.first, name.last)
		if err != nil {
			log.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("%s: %s (%.2f%% - %s)\n",
			name.first, result.Gender, result.Confidence*100, result.Strategy)
	}
}
