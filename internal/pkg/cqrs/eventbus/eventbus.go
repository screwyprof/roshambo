package eventbus

import (
	"sync"

	"github.com/screwyprof/roshambo/pkg/domain"
)

type InMemoryEventBus struct {
	eventHandlers   map[domain.EventHandler]struct{}
	eventHandlersMu sync.RWMutex
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		eventHandlers: make(map[domain.EventHandler]struct{}),
	}
}

func (b *InMemoryEventBus) Register(h domain.EventHandler) {
	b.eventHandlersMu.Lock()
	defer b.eventHandlersMu.Unlock()

	b.eventHandlers[h] = struct{}{}
}

func (b *InMemoryEventBus) Publish(events ...domain.DomainEvent) error {
	b.eventHandlersMu.RLock()
	defer b.eventHandlersMu.RUnlock()

	for h := range b.eventHandlers {
		if err := h.Handle(events...); err != nil {
			return err
		}
	}
	return nil
}
