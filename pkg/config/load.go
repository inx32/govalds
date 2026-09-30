package config

import (
	"context"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"forge.pi.home.arpa/golang/validator"
)

type contextConfigKey struct{}

func SetContext(ctx context.Context, config *Config) context.Context {
	return context.WithValue(ctx, contextConfigKey{}, config)
}

func FromContext(ctx context.Context) *Config {
	val := ctx.Value(contextConfigKey{})
	if val == nil {
		return nil
	}
	return val.(*Config)
}

func Load(path string) (*Config, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file, yaml.DisallowUnknownField())
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("config format is invalid:\n\n%w", err)
	}

	if err := validator.Validate(&config); err != nil {
		return nil, fmt.Errorf("unable to validate config:\n\n%w", err)
	}

	return &config, nil
}
