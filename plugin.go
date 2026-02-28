package logmsglint

import (
	"encoding/json"
	"fmt"
	"os"

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
	b, _ := json.Marshal(conf)
	fmt.Fprintln(os.Stderr, "LOGMSGLINT conf =", string(b))

	cfg, err := loganalysis.ConfigFromSettings(conf)
	if err != nil {
		return nil, err
	}

	b2, _ := json.Marshal(cfg)
	fmt.Fprintln(os.Stderr, "LOGMSGLINT cfg  =", string(b2))

	return plugin{cfg: cfg}, nil
}

func (p plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		loganalysis.NewAnalyzer(p.cfg),
	}, nil
}

func (p plugin) GetLoadMode() string {
	return register.LoadModeSyntax
}
