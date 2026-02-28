package rules

import (
	"testing"

	"github.com/bauerexe/logmsglint/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestEnglishRuleCheck(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
		want    domain.ViolationCode
	}{
		{name: "ok english", message: "starting server", want: ""},
		{name: "violation cyrillic", message: "запуск сервера", want: domain.ViolationEnglish},
	}

	rule := EnglishRule{}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			v := rule.Check(domain.LogCall{Message: tc.message})
			if tc.want == "" {
				assert.Nil(t, v)
				return
			}

			if assert.NotNil(t, v) {
				assert.Equal(t, tc.want, v.Code)
			}
		})
	}
}
