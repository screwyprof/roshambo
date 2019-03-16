package dispatcher

import "github.com/screwyprof/roshambo/pkg/domain"

// Dispatcher is a basic message dispatcher.
//
// It drives the overall command handling and event application/distribution process.
// It is suitable for a simple, single node application that can safely build its subscriber list
// at startup and keep it in memory.
// Depends on some kind of event storage mechanism.
type Dispatcher struct {
	eventStore       domain.EventStore
	aggregateFactory domain.AggregateFactory
}

func NewDispatcher(eventStore domain.EventStore, aggregateFactory domain.AggregateFactory) *Dispatcher {
	if eventStore == nil {
		panic("eventStore is required")
	}

	if aggregateFactory == nil {
		panic("aggregateFactory is required")
	}

	return &Dispatcher{
		eventStore:       eventStore,
		aggregateFactory: aggregateFactory,
	}
}

func (d *Dispatcher) Handle(c domain.Command) ([]domain.DomainEvent, error) {
	agg := d.aggregateFactory.CreateAggregate(c.AggregateType(), c.AggregateID())

	loadedEvents, err := d.eventStore.LoadEventsFor(c.AggregateID())
	if err != nil {
		return nil, err
	}

	_ = agg.Apply(loadedEvents...)

	return nil, nil
}
