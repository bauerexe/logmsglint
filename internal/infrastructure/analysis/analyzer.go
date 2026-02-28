package analysis

import (
	"fmt"
	"go/ast"
	"os"
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

func NewAnalyzer(cfg Config) *analysis.Analyzer {
	cfgCopy := cfg
	return &analysis.Analyzer{
		Name:     "logmsglint",
		Doc:      "log message must be in English",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run: func(pass *analysis.Pass) (any, error) {
			return run(pass, cfgCopy)
		},
	}
}

func newValidator(cfg Config) (*usecase.Validator, error) {
	fmt.Fprintf(os.Stderr, "LOGMSGLINT runtime cfg = %+v\n", cfg)

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

	fmt.Fprintf(os.Stderr, "LOGMSGLINT rules count=%d sensitive=%v\n", len(rulesList), cfg.Rules.Sensitive)

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
		violations = filterViolationsByConfig(violations, cfg)
		if !cfg.Rules.Sensitive && len(violations) > 0 {
			filtered := violations[:0]
			for _, v := range violations {
				if v.Code == domain.ViolationSensitive {
					continue
				}
				filtered = append(filtered, v)
			}
			violations = filtered
		}

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
		if !cfg.Rules.Lowercase {
			return analysis.SuggestedFix{}, false
		}
		first, size := utf8.DecodeRuneInString(message)
		if first == utf8.RuneError && size == 0 {
			return analysis.SuggestedFix{}, false
		}
		replacement = string(unicode.ToLower(first)) + message[size:]

	case domain.ViolationEnglish:
		if !cfg.Rules.English {
			return analysis.SuggestedFix{}, false
		}
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
		if !cfg.Rules.NoSpecial {
			return analysis.SuggestedFix{}, false
		}
		replacement = strings.Map(func(r rune) rune {
			if unicode.IsControl(r) || unicode.IsSymbol(r) {
				return -1
			}
			return r
		}, message)
		replacement = collapsePunctuation(replacement)

	case domain.ViolationSensitive:
		if !cfg.Rules.Sensitive {
			return analysis.SuggestedFix{}, false
		}
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
func filterViolationsByConfig(violations []domain.Violation, cfg Config) []domain.Violation {
	if len(violations) == 0 {
		return violations
	}

	out := violations[:0]
	for _, v := range violations {
		switch v.Code {
		case domain.ViolationLowercase:
			if !cfg.Rules.Lowercase {
				continue
			}
		case domain.ViolationEnglish:
			if !cfg.Rules.English {
				continue
			}
		case domain.ViolationNoSpecial:
			if !cfg.Rules.NoSpecial {
				continue
			}
		case domain.ViolationSensitive:
			if !cfg.Rules.Sensitive {
				continue
			}
		}
		out = append(out, v)
	}
	return out
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
