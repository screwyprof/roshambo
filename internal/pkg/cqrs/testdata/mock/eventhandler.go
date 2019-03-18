package mock

import (
	"errors"

	"github.com/screwyprof/roshambo/pkg/domain"
)

var (
	ErrCannotHandleEvent = errors.New("cannot handle event")
)

type EventHandlerMock struct {
	Err error
	Matcher domain.EventMatcher
	Happened []domain.DomainEvent
}

func (h *EventHandlerMock) SubscribedTo() domain.EventMatcher {
	if h.Matcher != nil {
		return h.Matcher
	}
	return domain.MatchAnyEventOf("SomethingHappened", "SomethingElseHappened")
}

func (h *EventHandlerMock) Handle(events ...domain.DomainEvent) error {
	if h.Err != nil {
		return h.Err
	}

	for _, e := range events {
		h.handle(e)
	}

	return nil
}

func (h *EventHandlerMock) handle(event domain.DomainEvent) {
	switch e := event.(type) {
	case SomethingHappened:
		h.OnSomethingHappened(e)
	case SomethingElseHappened:
		h.OnSomethingElseHappened(e)
	}
}

func (h *EventHandlerMock) OnSomethingHappened(e SomethingHappened) {
	h.Happened = append(h.Happened, e)
}

func (h *EventHandlerMock) OnSomethingElseHappened(e SomethingElseHappened) {
	h.Happened = append(h.Happened, e)
}