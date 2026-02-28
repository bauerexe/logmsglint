package rules

import (
	"testing"

	"github.com/bauerexe/logmsglint/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestSensitiveRuleCheck(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
		want    domain.ViolationCode
	}{
		{name: "ok regular message", message: "starting server", want: ""},
		{name: "violation password", message: "password was updated", want: domain.ViolationSensitive},
		{name: "violation bearer", message: "bearer token accepted", want: domain.ViolationSensitive},
		{name: "violation custom pattern", message: "card 1234-5678-9012-3456", want: domain.ViolationSensitive},
	}

	rule, err := NewSensitiveRule(nil, []string{`\b\d{4}-\d{4}-\d{4}-\d{4}\b`})
	assert.NoError(t, err)
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
