package main

import (
	"github.com/bauerex/logmsglint/internal/infrastructure/analysis"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analysis.NewAnalyzer())
}
