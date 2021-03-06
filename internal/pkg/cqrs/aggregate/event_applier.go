package aggregate

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// EventApplier applies events for the registered appliers.
type EventApplier struct {
	appliers   map[string]domain.EventApplierFunc
	appliersMu sync.RWMutex
}

// NewEventApplier creates a new instance of EventApplier.
func NewEventApplier() *EventApplier {
	return &EventApplier{
		appliers: make(map[string]domain.EventApplierFunc),
	}
}

// RegisterAppliers registers all the event appliers found in the aggregate.
func (a *EventApplier) RegisterAppliers(aggregate domain.Aggregate) {
	aggregateType := reflect.TypeOf(aggregate)
	for i := 0; i < aggregateType.NumMethod(); i++ {
		method := aggregateType.Method(i)
		if !strings.HasPrefix(method.Name, "On") {
			continue
		}

		a.RegisterApplier(method.Name, func(e domain.DomainEvent) {
			method.Func.Call([]reflect.Value{reflect.ValueOf(aggregate), reflect.ValueOf(e)})
		})
	}
}

// RegisterApplier registers an event applier for the given method.
func (a *EventApplier) RegisterApplier(method string, applier domain.EventApplierFunc) {
	a.appliersMu.Lock()
	defer a.appliersMu.Unlock()
	a.appliers[method] = applier
}

// Apply implements domain.EventApplier interface.
func (a *EventApplier) Apply(events ...domain.DomainEvent) error {
	for _, e := range events {
		if err := a.apply(e); err != nil {
			return err
		}
	}
	return nil
}

func (a *EventApplier) apply(event domain.DomainEvent) error {
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
