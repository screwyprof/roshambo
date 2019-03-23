package eventhandler

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// EventHandler handles events.
type EventHandler struct {
	handlers   map[string]domain.EventHandlerFunc
	handlersMu sync.RWMutex
}

// New creates new instance of New.
func New() *EventHandler {
	return &EventHandler{
		handlers: make(map[string]domain.EventHandlerFunc),
	}
}

// RegisterHandler registers an event handler for the given method.
func (s *EventHandler) RegisterHandler(method string, handler domain.EventHandlerFunc) {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()
	s.handlers[method] = handler
}

// SubscribedTo implements domain.EventHandler interface.
func (s *EventHandler) SubscribedTo() domain.EventMatcher {
	var subscribedTo []string
	for m := range s.handlers {
		subscribedTo = append(subscribedTo, strings.TrimPrefix(m, "On"))
	}
	return domain.MatchAnyEventOf(subscribedTo...)
}

// Handle implements domain.EventHandler interface.
func (s *EventHandler) Handle(e domain.DomainEvent) error {
	s.handlersMu.RLock()
	defer s.handlersMu.RUnlock()

	handlerID := "On" + e.EventType()
	handler, ok := s.handlers[handlerID]
	if !ok {
		return fmt.Errorf("event handler for %s event is not found", handlerID)
	}

	return handler(e)
}

// RegisterHandlers registers all the event handlers found in .
func (h *EventHandler) RegisterHandlers(entity interface{}) {
	entityType := reflect.TypeOf(entity)
	for i := 0; i < entityType.NumMethod(); i++ {
		method := entityType.Method(i)
		h.registerHandlerDynamically(method, entity)
	}
}

func (h *EventHandler) registerHandlerDynamically(method reflect.Method, entity interface{}) {
	if !strings.HasPrefix(method.Name, "On") {
		return
	}

	h.RegisterHandler(method.Name, func(e domain.DomainEvent) error {
		return h.invokeEventHandler(method, entity, e)
	})
}

func (h *EventHandler) invokeEventHandler(method reflect.Method, entity interface{}, e domain.DomainEvent) error {
	result := method.Func.Call([]reflect.Value{reflect.ValueOf(entity), reflect.ValueOf(e)})
	resErr := result[0].Interface()
	if resErr != nil {
		return resErr.(error)
	}
	return nil
}
