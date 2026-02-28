package analysis

import (
	"encoding/json"
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = ".logmsglint.yml"

type Config struct {
	Rules     RulesConfig     `json:"rules" yaml:"rules"`
	Sensitive SensitiveConfig `json:"sensitive" yaml:"sensitive"`
}

type RulesConfig struct {
	Lowercase bool `json:"lowercase" yaml:"lowercase"`
	English   bool `json:"english" yaml:"english"`
	NoSpecial bool `json:"nospecial" yaml:"nospecial"`
	Sensitive bool `json:"sensitive" yaml:"sensitive"`
}

type SensitiveConfig struct {
	Keywords []string `json:"keywords" yaml:"keywords"`
	Patterns []string `json:"patterns" yaml:"patterns"`
}

func defaultConfig() Config {
	return Config{
		Rules: RulesConfig{
			Lowercase: true,
			English:   true,
			NoSpecial: true,
			Sensitive: true,
		},
	}
}

func loadConfig(path string) (Config, error) {
	cfg := defaultConfig()
	if path == "" {
		path = defaultConfigPath
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return Config{}, err
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		if errJSON := json.Unmarshal(data, &cfg); errJSON != nil {
			return Config{}, err
		}
	}

	return cfg, nil
}
