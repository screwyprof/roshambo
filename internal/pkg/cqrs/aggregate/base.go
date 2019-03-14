package aggregate

import "github.com/screwyprof/roshambo/pkg/domain"

// Base implements a basic aggregate root.
type Base struct {
	domain.Aggregate
}

// NewBase creates a new instance of Base.
func NewBase(pureAgg domain.Aggregate) *Base {
	if pureAgg == nil {
		panic("pureAgg is required")
	}
	return &Base{Aggregate: pureAgg}
}
