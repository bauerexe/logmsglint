package analysis

import (
	"fmt"
	"go/ast"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/bauerexe/logmsglint/internal/app/usecase"
	"github.com/bauerexe/logmsglint/internal/domain"
	"github.com/bauerexe/logmsglint/internal/infrastructure/rules"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const AnalyzerName = "logmsglint"

func NewAnalyzer() *analysis.Analyzer {
	var configPath string
	analyzer := &analysis.Analyzer{
		Name:     AnalyzerName,
		Doc:      "checks logging message quality in slog and zap calls",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
			// Приоритет:
			// 1) -config флаг (CLI)
			// 2) inlineConfig (golangci-lint plugin)
			// 3) .logmsglint.yml / дефолты
			var (
				cfg Config
				err error
			)

			switch {
			case configPath != "":
				cfg, err = loadConfig(configPath)
			case inlineConfig != nil:
				cfg = *inlineConfig
			default:
				cfg, err = loadConfig("")
			}

			if err != nil {
				return nil, fmt.Errorf("load config: %w", err)
			}
			return run(pass, cfg)
		},
	}
	analyzer.Flags.StringVar(&configPath, "config", "", "path to logmsglint config file")

	return analyzer
}

func newValidator(cfg Config) (*usecase.Validator, error) {
	rulesList := make([]usecase.Rule, 0, 4)
	if cfg.Rules.Lowercase {
		rulesList = append(rulesList, rules.LowercaseRule{})
	}
	if cfg.Rules.English {
		rulesList = append(rulesList, rules.EnglishRule{})
	}
	if cfg.Rules.NoSpecial {
		rulesList = append(rulesList, rules.NoSpecialRule{})
	}
	if cfg.Rules.Sensitive {
		sensitiveRule, err := rules.NewSensitiveRule(cfg.Sensitive.Keywords, cfg.Sensitive.Patterns)
		if err != nil {
			return nil, err
		}
		rulesList = append(rulesList, sensitiveRule)
	}

	return usecase.NewValidator(rulesList), nil
}

func run(pass *analysis.Pass, cfg Config) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	validator, err := newValidator(cfg)
	if err != nil {
		return nil, err
	}

	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}
	insp.Preorder(nodeFilter, func(node ast.Node) {
		callExpr := node.(*ast.CallExpr)
		extracted, ok := extractLogCall(pass, callExpr)
		if !ok {
			return
		}

		violations := validator.Validate(extracted.Call)
		for _, violation := range violations {
			diagnostic := analysis.Diagnostic{
				Pos:     extracted.MessagePos,
				End:     extracted.MessageEnd,
				Message: formatDiagnostic(violation),
			}

			if fix, ok := makeSuggestedFix(extracted, violation, cfg); ok {
				diagnostic.SuggestedFixes = []analysis.SuggestedFix{fix}
			}

			pass.Report(diagnostic)
		}
	})

	return nil, nil
}

func formatDiagnostic(violation domain.Violation) string {
	return fmt.Sprintf("%s: %s", violation.Code, violation.Message)
}

func makeSuggestedFix(extracted *ExtractedLogCall, violation domain.Violation, cfg Config) (analysis.SuggestedFix, bool) {
	message := extracted.Call.Message
	if message == "" {
		return analysis.SuggestedFix{}, false
	}

	replacement := message
	switch violation.Code {
	case domain.ViolationLowercase:
		first, size := utf8.DecodeRuneInString(message)
		if first == utf8.RuneError && size == 0 {
			return analysis.SuggestedFix{}, false
		}
		replacement = string(unicode.ToLower(first)) + message[size:]
	case domain.ViolationEnglish:
		replacement = strings.Map(func(r rune) rune {
			if r == '\n' || r == '\t' || r == '\r' {
				return -1
			}
			if r >= 0x20 && r <= 0x7E {
				return r
			}
			if unicode.IsSpace(r) {
				return ' '
			}
			return -1
		}, message)
	case domain.ViolationNoSpecial:
		replacement = strings.Map(func(r rune) rune {
			if unicode.IsControl(r) || unicode.IsSymbol(r) {
				return -1
			}
			return r
		}, message)
		replacement = collapsePunctuation(replacement)
	case domain.ViolationSensitive:
		replacement = redactSensitive(message, cfg.Sensitive)
	default:
		return analysis.SuggestedFix{}, false
	}

	if replacement == "" || replacement == message {
		return analysis.SuggestedFix{}, false
	}

	return analysis.SuggestedFix{
		Message: "auto-fix log message",
		TextEdits: []analysis.TextEdit{{
			Pos:     extracted.MessagePos,
			End:     extracted.MessageEnd,
			NewText: []byte(strconv.Quote(replacement)),
		}},
	}, true
}

func redactSensitive(message string, cfg SensitiveConfig) string {
	replacement := message
	keywords := cfg.Keywords
	if len(keywords) == 0 {
		keywords = rules.DefaultSensitiveKeywords
	}

	for _, keyword := range keywords {
		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(keyword))
		replacement = re.ReplaceAllString(replacement, "[REDACTED]")
	}
	for _, pattern := range cfg.Patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			continue
		}
		replacement = re.ReplaceAllString(replacement, "[REDACTED]")
	}

	return replacement
}

func collapsePunctuation(input string) string {
	if input == "" {
		return input
	}

	var b strings.Builder
	var prev rune
	for i, r := range input {
		if i > 0 && (r == '.' || r == '!' || r == '?') && r == prev {
			continue
		}
		b.WriteRune(r)
		prev = r
	}

	return b.String()
}

var inlineConfig *Config

func SetInlineConfig(cfg Config) {
	c := cfg
	inlineConfig = &c
}
