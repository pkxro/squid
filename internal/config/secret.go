package config

import "github.com/kelseyhightower/envconfig"

// SecretConfig is a generic config used to fetch env attributes
type SecretConfig struct {
	RpcUrl        string `envconfig:"RPC_URL"`
	CacheUri      string `envconfig:"CACHE_URI"`
	CachePassword string `envconfig:"CACHE_PASSWORD"`
	FeePayer      string `envconfig:"FEE_PAYER"`
}

// NewSecretConfig generates a new instance of GlobalConfig
func NewSecretConfig() (*SecretConfig, error) {
	var sec SecretConfig

	err := envconfig.Process("", &sec)
	if err != nil {
		return nil, err
	}

	return &sec, nil
}
