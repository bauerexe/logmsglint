package rules

import (
	"testing"

	"github.com/bauerexe/logmsglint/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNoSpecialRuleCheck(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
		want    domain.ViolationCode
	}{
		{name: "ok normal punctuation", message: "starting server.", want: ""},
		{name: "violation repeated punctuation", message: "starting server!!!", want: domain.ViolationNoSpecial},
		{name: "violation ellipsis", message: "starting...", want: domain.ViolationNoSpecial},
		{name: "violation emoji", message: "starting 🚀", want: domain.ViolationNoSpecial},
	}

	rule := NoSpecialRule{}
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
