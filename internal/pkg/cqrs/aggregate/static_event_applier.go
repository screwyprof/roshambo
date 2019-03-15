package aggregate

import (
	"fmt"
	"sync"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// StaticEventApplier applies events for the registered appliers.
type StaticEventApplier struct {
	appliers   map[string]domain.EventApplierFunc
	appliersMu sync.RWMutex
}

// NewStaticEventApplier creates a new instance of StaticEventApplier.
func NewStaticEventApplier() *StaticEventApplier {
	return &StaticEventApplier{
		appliers: make(map[string]domain.EventApplierFunc),
	}
}

// RegisterAppliers does nothing at the moment.
// This method is present to satisfy eventApplier interface of base aggregate.
func (a *StaticEventApplier) RegisterAppliers(aggregate domain.Aggregate) {}

// RegisterApplier registers an event applier for the given method.
func (a *StaticEventApplier) RegisterApplier(method string, applier domain.EventApplierFunc) {
	a.appliersMu.Lock()
	defer a.appliersMu.Unlock()
	a.appliers[method] = applier
}

// Apply implements domain.EventApplier interface.
func (a *StaticEventApplier) Apply(events ...domain.DomainEvent) error {
	for _, e := range events {
		if err := a.apply(e); err != nil {
			return err
		}
	}

	return nil
}

func (a *StaticEventApplier) apply(event domain.DomainEvent) error {
	a.appliersMu.RLock()
	defer a.appliersMu.RUnlock()

	applierID := "On" + event.EventType()
	applier, ok := a.appliers[applierID]
	if !ok {
		return fmt.Errorf("event applier for %s event is not found", applierID)
	}
	applier(event)

	return nil
}
