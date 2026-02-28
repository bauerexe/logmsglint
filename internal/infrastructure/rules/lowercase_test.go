package rules

import (
	"testing"

	"github.com/bauerexe/logmsglint/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestLowercaseRuleCheck(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
		want    domain.ViolationCode
	}{
		{name: "ok lowercase", message: "starting server", want: ""},
		{name: "ok empty", message: "", want: ""},
		{name: "violation uppercase", message: "Starting server", want: domain.ViolationLowercase},
	}

	rule := LowercaseRule{}
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
