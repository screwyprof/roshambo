package mock

import (
	"errors"

	"github.com/screwyprof/roshambo/pkg/domain"
)

var (
	// ErrEventStoreCannotLoadEvents happens when event store can't load events.
	ErrEventStoreCannotLoadEvents = errors.New("cannot load events")
	// ErrEventStoreCannotStoreEvents happens when event store can't store events.
	ErrEventStoreCannotStoreEvents = errors.New("cannot store events")
)

// EventStoreMock mocks event store.
type EventStoreMock struct {
	Loader func(aggregateID domain.Identifier) ([]domain.DomainEvent, error)
	Saver  func(aggregateID domain.Identifier, version int, events []domain.DomainEvent) error
}

// LoadEventsFor implements domain.EventStore interface.
func (m *EventStoreMock) LoadEventsFor(aggregateID domain.Identifier) ([]domain.DomainEvent, error) {
	return m.Loader(aggregateID)
}

// StoreEventsFor implements domain.EventStore interface.
func (m *EventStoreMock) StoreEventsFor(aggregateID domain.Identifier, version int, events []domain.DomainEvent) error {
	return m.Saver(aggregateID, version, events)
}
