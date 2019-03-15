package eventstore

import "github.com/screwyprof/roshambo/pkg/domain"

// InMemoryEventStore stores and loads events from memory.
type InMemoryEventStore struct{}

// NewInInMemoryEventStore creates a new instance of InMemoryEventStore.
func NewInInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{}
}

// LoadEventsFor loads events for the given aggregate.
func (s *InMemoryEventStore) LoadEventsFor(aggregateID domain.Identifier) ([]domain.DomainEvent, error) {
	return nil, nil
}

//StoreEventsFor saves evens of the given aggregate.
func (s *InMemoryEventStore) StoreEventsFor(
	aggregateID domain.Identifier, version int, events []domain.DomainEvent) error {
	return nil
}
