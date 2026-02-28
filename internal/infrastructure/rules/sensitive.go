package rules

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bauerexe/logmsglint/internal/domain"
)

var DefaultSensitiveKeywords = []string{
	"password",
	"passwd",
	"token",
	"api_key",
	"apikey",
	"secret",
	"bearer",
	"access_token",
	"refresh_token",
}

type SensitiveRule struct {
	Keywords []string
	Patterns []*regexp.Regexp
}

func NewSensitiveRule(keywords, patterns []string) (SensitiveRule, error) {
	if len(keywords) == 0 {
		keywords = DefaultSensitiveKeywords
	}

	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return SensitiveRule{}, fmt.Errorf("invalid sensitive pattern %q: %w", pattern, err)
		}
		compiled = append(compiled, re)
	}

	return SensitiveRule{Keywords: keywords, Patterns: compiled}, nil
}

func (r SensitiveRule) Check(call domain.LogCall) *domain.Violation {
	fmt.Fprintf(os.Stderr, "LOGMSGLINT SensitiveRule.Check called for message: %q\n", call.Message)
	message := strings.ToLower(call.Message)
	for _, keyword := range r.Keywords {
		if strings.Contains(message, strings.ToLower(keyword)) {
			return &domain.Violation{
				Code:    domain.ViolationSensitive,
				Message: "log message may contain sensitive data",
			}
		}
	}

	for _, pattern := range r.Patterns {
		if pattern.MatchString(call.Message) {
			return &domain.Violation{
				Code:    domain.ViolationSensitive,
				Message: "log message may contain sensitive data",
			}
		}
	}

	return nil
}
