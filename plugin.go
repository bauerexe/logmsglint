package logmsglint

import (
	myanalysis "github.com/bauerexe/logmsglint/internal/infrastructure/analysis"
	"golang.org/x/tools/go/analysis"
)

// Analyzer — экспортируемый анализатор, который увидит golangci-lint custom.
var Analyzer *analysis.Analyzer = NewAnalyzer()

// NewAnalyzer возвращает твой анализатор.
func NewAnalyzer() *analysis.Analyzer {
	return myanalysis.NewAnalyzer()
}
