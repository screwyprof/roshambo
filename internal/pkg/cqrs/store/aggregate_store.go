package store

import "github.com/screwyprof/roshambo/pkg/domain"

// AggregateStore loads and stores aggregates.
type AggregateStore struct {
	aggregateFactory domain.AggregateFactory
	eventStore       domain.EventStore
}

// NewStore creates a new instance of AggregateStore.
func NewStore(eventStore domain.EventStore, aggregateFactory domain.AggregateFactory) *AggregateStore {
	if eventStore == nil {
		panic("eventStore is required")
	}

	if aggregateFactory == nil {
		panic("aggregateFactory is required")
	}

	return &AggregateStore{
		eventStore:       eventStore,
		aggregateFactory: aggregateFactory,
	}
}

// Load implements domain.AggregateStore interface.
func (s *AggregateStore) Load(aggregateID domain.Identifier, aggregateType string) (domain.AdvancedAggregate, error) {
	loadedEvents, err := s.eventStore.LoadEventsFor(aggregateID)
	if err != nil {
		return nil, err
	}

	agg, err := s.aggregateFactory.CreateAggregate(aggregateType, aggregateID)
	if err != nil {
		return nil, err
	}

	err = agg.Apply(loadedEvents...)
	if err != nil {
		return nil, err
	}

	return agg, nil
}

// Store implements domain.AggregateStore interface.
func (s *AggregateStore) Store(agg domain.AdvancedAggregate, events ...domain.DomainEvent) error {
	return s.eventStore.StoreEventsFor(agg.AggregateID(), agg.Version(), events)
}
