package aggregate_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"

	"github.com/screwyprof/roshambo/pkg/domain"
)

func TestNewBase(t *testing.T) {
	t.Run("ItPanicsIfThePureAggregateIsNotGiven", func(t *testing.T) {
		factory := func() {
			aggregate.NewBase(nil, nil)
		}
		assert.Panic(t, factory)
	})

	t.Run("ItReturnsAggregateID", func(t *testing.T) {
		ID := testdata.StringIdentifier("TestAgg1")
		pureAgg := testdata.NewTestAggregate(ID)
		agg := aggregate.NewBase(pureAgg, nil)

		assert.Equals(t, ID, agg.AggregateID())
	})
}

func TestBaseHandle(t *testing.T) {
	t.Run("ItHandlesTheGivenCommand", func(t *testing.T) {
		// arrange
		ID := testdata.StringIdentifier("TestAgg1")
		pureAgg := testdata.NewTestAggregate(ID)

		commandHandler := aggregate.NewStaticCommandHandler()
		commandHandler.RegisterHandler(
			"MakeSomethingHappen",
			func(c domain.Command) ([]domain.DomainEvent, error) {
				return pureAgg.MakeSomethingHappen(c.(testdata.MakeSomethingHappen))
			},
		)

		agg := aggregate.NewBase(pureAgg, commandHandler)
		expectedEvents := []domain.DomainEvent{testdata.SomethingHappened{}}

		// act
		events, err := agg.Handle(testdata.MakeSomethingHappen{})

		// assert
		assert.Ok(t, err)
		assert.Equals(t, expectedEvents, events)
	})
}
