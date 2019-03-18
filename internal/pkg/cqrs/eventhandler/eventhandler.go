package eventhandler

import "github.com/screwyprof/roshambo/pkg/domain"

// EventHandler handles events.
type EventHandler struct{}

// NewEventHandler creates new instance of NewEventHandler.
func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

// SubscribedTo implements domain.EventHandler interface.
func (*EventHandler) SubscribedTo() domain.EventMatcher {
	panic("implement me")
}

// Handle implements domain.EventHandler interface.
func (*EventHandler) Handle(domain.DomainEvent) error {
	panic("implement me")
}
