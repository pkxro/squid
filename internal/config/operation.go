package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pkxro/squid/internal/model"
	"golang.org/x/time/rate"
)

const ConfigPath = "./operation.json"

type OperationConfig struct {
	Environment       model.ApplicationEnvironment `json:"environment"`
	Port              string                       `json:"port"`
	RateLimit         rate.Limit                   `json:"rateLimit"`
	RateLimitInterval uint                         `json:"rateLimitInterval"`
	CacheLimitMinutes int                          `json:"cacheLimitMinutes"`
	AllowedTokens     []model.Token                `json:"alowedTokens"`
}

// ParseJSONFile reads a JSON file and unmarshals it into the provided struct
func ParseJSONFile(filePath string, out interface{}) error {
	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %w", err)
	}

	// Unmarshal the JSON into the provided struct
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return nil
}

func NewOperationConfig() (*OperationConfig, error) {
	config := &OperationConfig{}
	err := ParseJSONFile(ConfigPath, config)
	if err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}
	return config, nil
}
