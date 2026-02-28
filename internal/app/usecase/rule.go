package usecase

import "github.com/bauerex/logmsglint/internal/domain"

// Rule describes validation contract used by application use case.
type Rule interface {
	Check(call domain.LogCall) *domain.Violation
}
