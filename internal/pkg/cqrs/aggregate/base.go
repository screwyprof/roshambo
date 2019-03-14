package aggregate

import (
	"github.com/screwyprof/roshambo/pkg/domain"
)

type commandHandler interface {
	domain.CommandHandler
	RegisterHandlers(aggregate domain.Aggregate)
	RegisterHandler(method string, handlerFunc domain.CommandHandlerFunc)
}

// Base implements a basic aggregate root.
type Base struct {
	domain.Aggregate
	commandHandler commandHandler
}

// NewBase creates a new instance of Base.
func NewBase(pureAgg domain.Aggregate, handler commandHandler) *Base {
	if pureAgg == nil {
		panic("pureAgg is required")
	}

	if handler == nil {
		handler = NewStaticCommandHandler()
	}

	handler.RegisterHandlers(pureAgg)

	return &Base{
		Aggregate:      pureAgg,
		commandHandler: handler,
	}
}

// Handle implements domain.CommandHandler.
func (b *Base) Handle(c domain.Command) ([]domain.DomainEvent, error) {
	return b.commandHandler.Handle(c)
}
