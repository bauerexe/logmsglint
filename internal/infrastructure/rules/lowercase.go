package rules

import (
	"unicode"

	"github.com/bauerexe/logmsglint/internal/domain"
)

type LowercaseRule struct{}

func (r LowercaseRule) Check(call domain.LogCall) *domain.Violation {
	if call.Message == "" {
		return nil
	}

	for _, ch := range call.Message {
		if unicode.IsUpper(ch) {
			return &domain.Violation{
				Code:    domain.ViolationLowercase,
				Message: "log message must start with a lowercase letter",
			}
		}
		break
	}

	return nil
}
