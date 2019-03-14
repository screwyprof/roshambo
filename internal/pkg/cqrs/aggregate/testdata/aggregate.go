package testdata

import "github.com/screwyprof/roshambo/pkg/domain"

type StringIdentifier string
func (i StringIdentifier) String() string {
	return string(i)
}

// TestAggregate a pure aggregate (has no external dependencies or dark magic method) used for testing.
type TestAggregate struct {
	id domain.Identifier
}

// NewTestAggregate creates a new instance of TestAggregate.
func NewTestAggregate(ID domain.Identifier) *TestAggregate {
	return &TestAggregate{id:ID}
}

// AggregateID implements domain.Aggregate interface.
func (a *TestAggregate) AggregateID() domain.Identifier {
	return a.id
}
