package rules

import (
	"unicode"

	"github.com/bauerexe/logmsglint/internal/domain"
)

type EnglishRule struct{}

func (r EnglishRule) Check(call domain.LogCall) *domain.Violation {
	for _, ch := range call.Message {
		if unicode.IsSpace(ch) {
			continue
		}
		if ch >= 0x20 && ch <= 0x7E {
			continue
		}
		return &domain.Violation{
			Code:    domain.ViolationEnglish,
			Message: "log message must be in English (ASCII only)",
		}
	}

	return nil
}
