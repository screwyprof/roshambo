package dispatcher

import "github.com/screwyprof/roshambo/pkg/domain"

// Dispatcher is a basic message dispatcher.
//
// It drives the overall command handling and event application/distribution process.
// It is suitable for a simple, single node application that can safely build its subscriber list
// at startup and keep it in memory.
// Depends on some kind of event storage mechanism.
type Dispatcher struct {
	store          domain.AggregateStore
	eventPublisher domain.EventPublisher
}

// NewDispatcher creates a new instance of Dispatcher.
func NewDispatcher(aggregateStore domain.AggregateStore, eventPublisher domain.EventPublisher) *Dispatcher {
	if aggregateStore == nil {
		panic("aggregateStore is required")
	}

	if eventPublisher == nil {
		panic("eventPublisher is required")
	}

	return &Dispatcher{
		store:          aggregateStore,
		eventPublisher: eventPublisher,
	}
}

// Handle implements domain.CommandHandler interface.
func (d *Dispatcher) Handle(c domain.Command) ([]domain.DomainEvent, error) {
	agg, err := d.store.Load(c.AggregateID(), c.AggregateType())
	if err != nil {
		return nil, err
	}

	events, err := agg.Handle(c)
	if err != nil {
		return nil, err
	}

	err = d.storeAndPublishEvents(agg, events...)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (d *Dispatcher) storeAndPublishEvents(aggregate domain.AdvancedAggregate, events ...domain.DomainEvent) error {
	err := d.store.Store(aggregate, events...)
	if err != nil {
		return err
	}

	err = d.eventPublisher.Publish(events...)
	if err != nil {
		return err
	}

	return nil
}
