package eventstore

import (
	"sync"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// InMemoryEventStore stores and loads events from memory.
type InMemoryEventStore struct {
	eventStreams   map[domain.Identifier][]domain.DomainEvent
	eventStreamsMu sync.RWMutex
}

// NewInInMemoryEventStore creates a new instance of InMemoryEventStore.
func NewInInMemoryEventStore() *InMemoryEventStore {
	return &InMemoryEventStore{
		eventStreams: make(map[domain.Identifier][]domain.DomainEvent),
	}
}

// LoadEventsFor loads events for the given aggregate.
func (s *InMemoryEventStore) LoadEventsFor(aggregateID domain.Identifier) ([]domain.DomainEvent, error) {
	s.eventStreamsMu.RLock()
	defer s.eventStreamsMu.RUnlock()

	return s.eventStreams[aggregateID], nil
}

//StoreEventsFor saves evens of the given aggregate.
func (s *InMemoryEventStore) StoreEventsFor(
	aggregateID domain.Identifier, version int, events []domain.DomainEvent) error {
	s.eventStreamsMu.Lock()
	defer s.eventStreamsMu.Unlock()

	s.eventStreams[aggregateID] = events
	return nil
}
