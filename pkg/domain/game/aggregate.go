package game

import "github.com/screwyprof/roshambo/pkg/domain"

type Aggregate struct {
	id domain.Identifier
}

// NewAggregate creates a new instance of Aggregate.
func NewAggregate(ID domain.Identifier) *Aggregate {
	if ID == nil {
		panic("ID is required")
	}
	return &Aggregate{id: ID}
}

// AggregateID implements domain.Aggregate interface.
func (a *Aggregate) AggregateID() domain.Identifier {
	return a.id
}
