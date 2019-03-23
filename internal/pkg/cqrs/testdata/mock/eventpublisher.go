package mock

import (
	"errors"

	"github.com/screwyprof/roshambo/pkg/domain"
)

var (
	// ErrCannotPublishEvents happens when event publisher cannot publish the given events.
	ErrCannotPublishEvents = errors.New("cannot load aggregate")
)

// EventPublisherMock mocks event store.
type EventPublisherMock struct {
	Publisher func(e ...domain.DomainEvent) error
}

// Publish implements domain.EventPublisher interface.
func (m *EventPublisherMock) Publish(e ...domain.DomainEvent) error {
	return m.Publisher(e...)
}

