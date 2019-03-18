package eventhandler

import (
	"fmt"
	"strings"
	"sync"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// Static handles events.
type Static struct {
	handlers   map[string]domain.EventHandlerFunc
	handlersMu sync.RWMutex
}

// NewStatic creates new instance of NewStatic.
func NewStatic() *Static {
	return &Static{
		handlers: make(map[string]domain.EventHandlerFunc),
	}
}

// RegisterHandler registers an event handler for the given method.
func (s *Static) RegisterHandler(method string, handler domain.EventHandlerFunc) {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.handlers[method] = handler
}

// SubscribedTo implements domain.EventHandler interface.
func (s *Static) SubscribedTo() domain.EventMatcher {
	var subscribedTo []string
	for m := range s.handlers {
		subscribedTo = append(subscribedTo, strings.TrimPrefix(m, "On"))
	}
	return domain.MatchAnyEventOf(subscribedTo...)
}

// Handle implements domain.EventHandler interface.
func (s *Static) Handle(e domain.DomainEvent) error {
	s.handlersMu.RLock()
	defer s.handlersMu.RUnlock()

	handlerID := "On" + e.EventType()
	handler, ok := s.handlers[handlerID]
	if !ok {
		return fmt.Errorf("event handler for %s event is not found", handlerID)
	}

	return handler(e)
}
