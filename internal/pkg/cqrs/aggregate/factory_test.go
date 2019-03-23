package aggregate_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// ensure that factory implements domain.AggregateFactory interface.
var _ domain.AggregateFactory = (*aggregate.Factory)(nil)

func TestNewFactory(t *testing.T) {
	t.Run("ItReturnsNewFactoryInstance", func(t *testing.T) {
		f := aggregate.NewFactory()
		assert.True(t, f != nil)
	})
}

func TestFactoryCreateAggregate(t *testing.T) {
	t.Run("ItPanicsIfTheAggregateIsNotRegistered", func(t *testing.T) {
		f := aggregate.NewFactory()

		_, err := f.CreateAggregate("mock.TestAggregate", mock.StringIdentifier("TestAgg"))

		assert.Equals(t, mock.ErrAggIsNotRegistered, err)
	})
}

func TestFactoryRegisterAggregate(t *testing.T) {
	t.Run("ItRegistersAnAggregateFactory", func(t *testing.T) {
		// arrange
		ID := mock.StringIdentifier("TestAgg")
		agg := mock.NewTestAggregate(ID)

		commandHandler := aggregate.NewCommandHandler()
		commandHandler.RegisterHandlers(agg)

		eventApplier := aggregate.NewEventApplier()
		eventApplier.RegisterAppliers(agg)

		expected := aggregate.NewAdvanced(
			agg,
			commandHandler,
			eventApplier,
		)

		f := aggregate.NewFactory()

		// act
		f.RegisterAggregate(func(ID domain.Identifier) domain.AdvancedAggregate {
			return expected
		})
		newAgg, err := f.CreateAggregate("mock.TestAggregate", ID)

		// assert
		assert.Ok(t, err)
		assert.Equals(t, expected, newAgg)
	})
}
