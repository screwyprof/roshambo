package aggregate

import (
	"github.com/screwyprof/roshambo/pkg/domain"
)

// Advanced implements an advanced aggregate root.
type Advanced struct {
	domain.Aggregate
	version int

	commandHandler domain.CommandHandler
	eventApplier   domain.EventApplier
}

// NewAdvanced creates a new instance of Advanced.
func NewAdvanced(pureAgg domain.Aggregate, commandHandler domain.CommandHandler, eventApplier domain.EventApplier) *Advanced {
	if pureAgg == nil {
		panic("pureAgg is required")
	}

	if commandHandler == nil {
		panic("commandHandler is required")
	}

	if eventApplier == nil {
		panic("eventApplier is required")
	}

	return &Advanced{
		Aggregate:      pureAgg,
		commandHandler: commandHandler,
		eventApplier:   eventApplier,
	}
}

// Version implements domain.Versionable interface.
func (b *Advanced) Version() int {
	return b.version
}

// Handle implements domain.CommandHandler.
func (b *Advanced) Handle(c domain.Command) ([]domain.DomainEvent, error) {
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
func (b *Advanced) Apply(e ...domain.DomainEvent) error {
	if err := b.eventApplier.Apply(e...); err != nil {
		return err
	}
	b.version += len(e)
	return nil
}
