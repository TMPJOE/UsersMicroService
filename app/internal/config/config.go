package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Server struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	ReadTimeout  string `yaml:"read_timeout"`
	WriteTimeout string `yaml:"write_timeout"`
	IdleTimeout  string `yaml:"idle_timeout"`
}

type JWT struct {
	Secret     string `yaml:"secret"`
	Issuer     string `yaml:"issuer"`
	Expiration string `yaml:"expiration"`
}

type Database struct {
	URL string `yaml:"url"`
}

type Health struct {
	Path      string `yaml:"path"`
	ReadyPath string `yaml:"ready_path"`
}

type Config struct {
	Server         Server                 `yaml:"server"`
	Logging        map[string]interface{} `yaml:"logging"`
	RateLimit      map[string]interface{} `yaml:"rate_limit"`
	CircuitBreaker map[string]interface{} `yaml:"circuit_breaker"`
	Health         Health                 `yaml:"health"`
}

// Load reads a YAML config file, expands environment variables, and returns the parsed config.
func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	expanded := os.ExpandEnv(string(b))

	var cfg Config
	if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
