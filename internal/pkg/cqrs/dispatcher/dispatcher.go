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

// NewDispatcher creates a new instance of Dispatcher.
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

// Handle implements domain.CommandHandler interface.
func (d *Dispatcher) Handle(c domain.Command) ([]domain.DomainEvent, error) {
	agg, err := d.loadAggregate(c)
	if err != nil {
		return nil, err
	}

	events, err := agg.Handle(c)
	if err != nil {
		return nil, err
	}

	err = d.eventStore.StoreEventsFor(agg.AggregateID(), agg.Version(), events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (d *Dispatcher) loadAggregate(c domain.Command) (domain.AdvancedAggregate, error) {
	loadedEvents, err := d.eventStore.LoadEventsFor(c.AggregateID())
	if err != nil {
		return nil, err
	}

	agg := d.aggregateFactory.CreateAggregate(c.AggregateType(), c.AggregateID())
	err = agg.Apply(loadedEvents...)
	if err != nil {
		return nil, err
	}

	return agg, nil
}
