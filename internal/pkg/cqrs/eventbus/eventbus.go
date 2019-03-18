package eventbus

import (
	"github.com/screwyprof/roshambo/pkg/domain"
)

type InMemoryEventBus struct{}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{}
}

func (b *InMemoryEventBus) Publish(events ...domain.DomainEvent) error {
	return nil
}
