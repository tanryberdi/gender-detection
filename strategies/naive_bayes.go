package strategies

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"

	"gender-detection/types"
)

type NaiveBayesModel struct {
	MaleCount      int            `json:"male_count"`
	FemaleCount    int            `json:"female_count"`
	MaleFeatures   map[string]int `json:"male_features"`
	FemaleFeatures map[string]int `json:"female_features"`
	Smoothing      float64        `json:"smoothing"`
}

type NaiveBayesStrategy struct {
	model *NaiveBayesModel
}

func NewNaiveBayesStrategy(path string) (*NaiveBayesStrategy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read model: %w", err)
	}

	var model NaiveBayesModel
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, fmt.Errorf("failed to parse model: %w", err)
	}

	if model.Smoothing == 0 {
		model.Smoothing = 1.0
	}

	return &NaiveBayesStrategy{model: &model}, nil
}

func (n *NaiveBayesStrategy) extractFeatures(name string) []string {
	name = strings.ToLower(strings.TrimSpace(name))
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

func (n *NaiveBayesStrategy) Detect(ctx context.Context, firstName, lastName string) (*types.Result, error) {
	firstName = strings.ToLower(strings.TrimSpace(firstName))
	if firstName == "" {
		return nil, nil
	}

	features := n.extractFeatures(firstName)

	totalCount := n.model.MaleCount + n.model.FemaleCount
	if totalCount == 0 {
		return nil, fmt.Errorf("model not trained")
	}

	pMale := float64(n.model.MaleCount) / float64(totalCount)
	pFemale := float64(n.model.FemaleCount) / float64(totalCount)

	logProbMale := math.Log(pMale)
	logProbFemale := math.Log(pFemale)

	vocabSize := len(n.model.MaleFeatures) + len(n.model.FemaleFeatures)

	for _, feature := range features {
		maleCount := float64(n.model.MaleFeatures[feature])
		femaleCount := float64(n.model.FemaleFeatures[feature])

		maleDenom := float64(n.model.MaleCount) + n.model.Smoothing*float64(vocabSize)
		femaleDenom := float64(n.model.FemaleCount) + n.model.Smoothing*float64(vocabSize)

		logProbMale += math.Log((maleCount + n.model.Smoothing) / maleDenom)
		logProbFemale += math.Log((femaleCount + n.model.Smoothing) / femaleDenom)
	}

	var gender types.Gender
	var confidence float64

	if logProbMale > logProbFemale {
		gender = types.Male
		total := math.Exp(logProbMale) + math.Exp(logProbFemale)
		confidence = math.Exp(logProbMale) / total
	} else {
		gender = types.Female
		total := math.Exp(logProbMale) + math.Exp(logProbFemale)
		confidence = math.Exp(logProbFemale) / total
	}

	return &types.Result{
		Gender:     gender,
		Confidence: confidence,
		Metadata: map[string]interface{}{
			"features_count":  len(features),
			"log_prob_male":   logProbMale,
			"log_prob_female": logProbFemale,
		},
	}, nil
}

func (n *NaiveBayesStrategy) Name() string {
	return "naive_bayes"
}

func (n *NaiveBayesStrategy) MinConfidence() float64 {
	return 0.70
}
