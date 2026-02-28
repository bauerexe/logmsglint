package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	loganalysis "github.com/bauerexe/logmsglint/internal/infrastructure/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	var (
		configPath     string
		settingsJSON   string
		settingsJSONFn string
	)

	flag.StringVar(&configPath, "config", "", "path to .logmsglint.yml (optional)")

	flag.StringVar(&settingsJSON, "settings-json", "", "settings JSON (like golangci-lint linter settings)")

	flag.StringVar(&settingsJSONFn, "settings-json-file", "", "path to settings JSON file")

	flag.Parse()

	cfg, err := resolveConfig(configPath, settingsJSON, settingsJSONFn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "config error:", err)
		os.Exit(2)
	}

	singlechecker.Main(loganalysis.NewAnalyzer(cfg))
}

func resolveConfig(configPath, settingsJSON, settingsJSONFn string) (loganalysis.Config, error) {

	cfg := loganalysis.DefaultConfig()
	if settingsJSONFn != "" {
		b, err := os.ReadFile(settingsJSONFn)
		if err != nil {
			return loganalysis.Config{}, err
		}
		settingsJSON = string(b)
	}
	if settingsJSON != "" {
		var anySettings any
		if err := json.Unmarshal([]byte(settingsJSON), &anySettings); err != nil {
			return loganalysis.Config{}, err
		}
		return loganalysis.ConfigFromSettings(anySettings)
	}

	if configPath != "" {
		return loganalysis.LoadConfig(configPath)
	}

	return cfg, nil
}
