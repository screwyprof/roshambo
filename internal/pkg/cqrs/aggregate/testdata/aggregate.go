package testdata

import (
	"errors"

	"github.com/screwyprof/roshambo/pkg/domain"
)

var (
	ErrMakeSomethingHandlerNotFound  = errors.New("handler for MakeSomethingHappen command is not found")
	ErrOnSomethingHappenedApplierNotFound  = errors.New("event applier for OnSomethingHappened event is not found")
)

type StringIdentifier string
func (i StringIdentifier) String() string {
	return string(i)
}

// TestAggregate a pure aggregate (has no external dependencies or dark magic method) used for testing.
type TestAggregate struct {
	id domain.Identifier
	happened bool
}

// NewTestAggregate creates a new instance of TestAggregate.
func NewTestAggregate(ID domain.Identifier) *TestAggregate {
	return &TestAggregate{id:ID}
}

// AggregateID implements domain.Aggregate interface.
func (a *TestAggregate) AggregateID() domain.Identifier {
	return a.id
}

func (a *TestAggregate) MakeSomethingHappen(c MakeSomethingHappen) ([]domain.DomainEvent, error) {
	return []domain.DomainEvent{SomethingHappened{}}, nil
}

func (a *TestAggregate) OnSomethingHappened(e SomethingHappened) {
	a.happened = true
}
