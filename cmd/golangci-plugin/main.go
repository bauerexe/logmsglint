package main

import (
	loganalysis "github.com/bauerexe/logmsglint/internal/infrastructure/analysis"
	"golang.org/x/tools/go/analysis"
)

// New is the entrypoint expected by golangci-lint module plugin system
func New(any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{loganalysis.NewAnalyzer()}, nil
}
