package aggregate

import (
	"fmt"
	"sync"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// StaticCommandHandler handles registered commands.
type StaticCommandHandler struct {
	handlers   map[string]domain.CommandHandlerFunc
	handlersMu sync.RWMutex
}

// NewStaticCommandHandler creates a new instance of StaticCommandHandler.
func NewStaticCommandHandler() *StaticCommandHandler {
	return &StaticCommandHandler{
		handlers: make(map[string]domain.CommandHandlerFunc),
	}
}

// RegisterHandlers does nothing at the moment.
// This method is present to satisfy commandHandler interface of base aggregate.
func (h *StaticCommandHandler) RegisterHandlers(aggregate domain.Aggregate) {}

// RegisterHandler registers a command handler for the given method.
func (h *StaticCommandHandler) RegisterHandler(method string, handler domain.CommandHandlerFunc) {
	h.handlersMu.Lock()
	defer h.handlersMu.Unlock()
	h.handlers[method] = handler
}

// Handle implements domain.CommandHandler interface.
func (h *StaticCommandHandler) Handle(c domain.Command) ([]domain.DomainEvent, error) {
	h.handlersMu.RLock()
	defer h.handlersMu.RUnlock()

	handler, ok := h.handlers[c.CommandType()]
	if !ok {
		return nil, fmt.Errorf("handler for %s command is not found", c.CommandType())
	}

	return handler(c)
}
