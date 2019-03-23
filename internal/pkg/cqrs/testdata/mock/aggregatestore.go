package mock

import (
	"errors"

	"github.com/screwyprof/roshambo/pkg/domain"
)

var (
	// ErrAggregateStoreCannotLoadAggregate happens when aggregate store can't load aggregate.
	ErrAggregateStoreCannotLoadAggregate = errors.New("cannot load aggregate")
	// ErrAggregateStoreCannotStoreAggregate happens when aggregate store can't store aggregate.
	ErrAggregateStoreCannotStoreAggregate = errors.New("cannot store aggregate")
)

// AggregateStoreMock mocks event store.
type AggregateStoreMock struct {
	Loader func(aggregateID domain.Identifier, aggregateType string) (domain.AdvancedAggregate, error)
	Saver func(aggregate domain.AdvancedAggregate, events ...domain.DomainEvent) error
}

// Load implements domain.AggregateStore interface.
func (m *AggregateStoreMock) Load(
	aggregateID domain.Identifier, aggregateType string) (domain.AdvancedAggregate, error) {
	return m.Loader(aggregateID, aggregateType)
}

// StoreEventsFor implements domain.AggregateStore interface.
func (m *AggregateStoreMock) Store(aggregate domain.AdvancedAggregate, events ...domain.DomainEvent) error {
	return m.Saver(aggregate, events...)
}
