package eventbus

import (
	"sync"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// InMemoryEventBus publishes events.
type InMemoryEventBus struct {
	eventHandlers   map[domain.EventHandler]struct{}
	eventHandlersMu sync.RWMutex
}

// NewInMemoryEventBus creates a new instance of InMemoryEventBus.
func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		eventHandlers: make(map[domain.EventHandler]struct{}),
	}
}

// Register registers event handler.
func (b *InMemoryEventBus) Register(h domain.EventHandler) {
	b.eventHandlersMu.Lock()
	defer b.eventHandlersMu.Unlock()

	b.eventHandlers[h] = struct{}{}
}

// Publish implements domain.EventPublisher interface.
func (b *InMemoryEventBus) Publish(events ...domain.DomainEvent) error {
	b.eventHandlersMu.RLock()
	defer b.eventHandlersMu.RUnlock()

	for h := range b.eventHandlers {
		if err := b.handleEvents(h, events...); err != nil {
			return err
		}
	}
	return nil
}

func (b *InMemoryEventBus) handleEvents(h domain.EventHandler, events ...domain.DomainEvent) error {
	for _, e := range events {
		err := b.handleEventIfMatches(h.SubscribedTo(), h, e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *InMemoryEventBus) handleEventIfMatches(
	m domain.EventMatcher, h domain.EventHandler, e domain.DomainEvent) error {
	if !m(e) {
		return nil
	}
	return h.Handle(e)
}
