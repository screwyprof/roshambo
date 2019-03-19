package eventhandler_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventhandler"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/roshambo/pkg/domain"
)

func TestDynamicRegisterHandlers(t *testing.T) {
	t.Run("ItReturnersTheEventsItSubscribedTo", func(t *testing.T) {
		// arrange
		eh := &mock.TestEventHandler{}

		s := eventhandler.NewDynamic()
		s.RegisterHandlers(eh)

		want := domain.MatchAnyEventOf("SomethingHappened", "SomethingElseHappened")

		// act
		got := s.SubscribedTo()

		// assert
		assert.Equals(t, want, got)
	})
}

func TestDynamicHandle(t *testing.T) {
	t.Run("ItHandlesTheGivenEvent", func(t *testing.T) {
		// arrange
		eh := &mock.TestEventHandler{}

		s := eventhandler.NewDynamic()
		s.RegisterHandlers(eh)

		// act
		err := s.Handle(mock.SomethingHappened{})

		// assert
		assert.Ok(t, err)
		assert.Equals(t, "test", eh.SomethingHappened)
	})

	t.Run("ItFailsIfEventHandlerReturnsAnError", func(t *testing.T) {
		// arrange
		eh := &mock.TestEventHandler{}

		s := eventhandler.NewDynamic()
		s.RegisterHandlers(eh)

		// act
		err := s.Handle(mock.SomethingElseHappened{})

		// assert
		assert.Equals(t, mock.ErrCannotHandleEvent, err)
	})
}
