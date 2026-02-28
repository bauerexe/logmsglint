package logmsglint

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	loganalysis "github.com/bauerexe/logmsglint/internal/infrastructure/analysis"
)

func init() {
	register.Plugin("logmsglint", New)
}

type plugin struct {
	cfg loganalysis.Config
}

func New(conf any) (register.LinterPlugin, error) {
	cfg, err := loganalysis.ConfigFromSettings(conf)
	if err != nil {
		return nil, err
	}
	return plugin{cfg: cfg}, nil
}

func (p plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	loganalysis.SetInlineConfig(p.cfg)
	return []*analysis.Analyzer{loganalysis.NewAnalyzer()}, nil
}

func (p plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
