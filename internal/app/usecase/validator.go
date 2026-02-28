package usecase

import "github.com/bauerex/logmsglint/internal/domain"

type Validator struct {
	rules []Rule
}

func NewValidator(rules []Rule) *Validator {
	copied := make([]Rule, 0, len(rules))
	copied = append(copied, rules...)

	return &Validator{rules: copied}
}

func (v *Validator) Validate(call domain.LogCall) []domain.Violation {
	if v == nil {
		return nil
	}

	violations := make([]domain.Violation, 0)
	for _, rule := range v.rules {
		if rule == nil {
			continue
		}

		if violation := rule.Check(call); violation != nil {
			violations = append(violations, *violation)
		}
	}

	return violations
}
