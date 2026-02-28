package rules

import (
	"regexp"
	"unicode"

	"github.com/bauerexe/logmsglint/internal/domain"
)

var repeatedPunctuationPattern = regexp.MustCompile(`(!!!+|\.\.\.+|\?\?+)`)

type NoSpecialRule struct{}

func (r NoSpecialRule) Check(call domain.LogCall) *domain.Violation {
	if repeatedPunctuationPattern.MatchString(call.Message) {
		return &domain.Violation{
			Code:    domain.ViolationNoSpecial,
			Message: "log message must not contain repeated punctuation like !!! or ...",
		}
	}

	for _, ch := range call.Message {
		if unicode.IsControl(ch) {
			return &domain.Violation{
				Code:    domain.ViolationNoSpecial,
				Message: "log message contains disallowed control/special characters",
			}
		}

		if unicode.IsSymbol(ch) {
			return &domain.Violation{
				Code:    domain.ViolationNoSpecial,
				Message: "log message contains disallowed symbols or emoji",
			}
		}
	}

	return nil
}
