package types

import "context"

type Gender string

const (
	Male    Gender = "male"
	Female  Gender = "female"
	Unknown Gender = "unknown"
)

type Result struct {
	Gender     Gender
	Confidence float64
	Strategy   string
	Metadata   map[string]interface{}
}

type Strategy interface {
	Detect(ctx context.Context, firstName, lastName string) (*Result, error)
	Name() string
	MinConfidence() float64
}

type Config struct {
	MinConfidence     float64
	EnableDictionary  bool
	EnableNaiveBayes  bool
	EnableRules       bool
	DictionaryPath    string
	ModelPath         string
	PatternsPath      string
	FallbackToUnknown bool
}
