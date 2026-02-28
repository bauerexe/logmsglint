package rules

import (
	"fmt"
	"regexp"

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
	// Временно всегда возвращаем nil, если хотим отключить
	return nil
}
