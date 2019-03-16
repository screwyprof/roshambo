package aggregate

import (
	"github.com/screwyprof/roshambo/pkg/domain"
)

type commandHandler interface {
	domain.CommandHandler
	RegisterHandlers(aggregate domain.Aggregate)
	RegisterHandler(method string, handlerFunc domain.CommandHandlerFunc)
}

type eventApplier interface {
	domain.EventApplier
	RegisterAppliers(aggregate domain.Aggregate)
	RegisterApplier(method string, applierFunc domain.EventApplierFunc)
}

// Base implements a basic aggregate root.
type Base struct {
	domain.Aggregate
	commandHandler commandHandler
	eventApplier   eventApplier
}

// NewBase creates a new instance of Base.
func NewBase(pureAgg domain.Aggregate, handler commandHandler, applier eventApplier) *Base {
	if pureAgg == nil {
		panic("pureAgg is required")
	}

	if handler == nil {
		handler = newDynamicCommandHandler()
	}

	if applier == nil {
		applier = newDynamicEventApplier()
	}

	handler.RegisterHandlers(pureAgg)
	applier.RegisterAppliers(pureAgg)

	return &Base{
		Aggregate:      pureAgg,
		commandHandler: handler,
		eventApplier:   applier,
	}
}

// Handle implements domain.CommandHandler.
func (b *Base) Handle(c domain.Command) ([]domain.DomainEvent, error) {
	events, err := b.commandHandler.Handle(c)
	if err != nil {
		return nil, err
	}

	if applierErr := b.eventApplier.Apply(events...); applierErr != nil {
		return nil, applierErr
	}

	return events, nil
}

// Apply implements domain.EventApplier interface.
func (b *Base) Apply(e ...domain.DomainEvent) error {
	if err := b.eventApplier.Apply(e...); err != nil {
		return err
	}
	b.Aggregate.IncrementVersion(len(e))
	return nil
}
