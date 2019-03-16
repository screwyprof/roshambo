package dispatcher_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testfixture"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/dispatcher"

	"github.com/screwyprof/roshambo/pkg/domain"
)

// ensure that Dispatcher  implements domain.CommandHandler interface.
var _ domain.CommandHandler = (*dispatcher.Dispatcher)(nil)

func TestNewDispatcher(t *testing.T) {
	t.Run("ItPanicsIfEventStoreIsNotGiven", func(t *testing.T) {
		factory := func() {
			dispatcher.NewDispatcher(nil, nil)
		}
		assert.Panic(t, factory)
	})

	t.Run("ItPanicsIfAggregateFactoryIsNotGiven", func(t *testing.T) {
		factory := func() {
			dispatcher.NewDispatcher(nil, nil)
		}
		assert.Panic(t, factory)
	})
}

func TestNewDispatcherHandle(t *testing.T) {
	t.Run("ItCreatesAggregate", func(t *testing.T) {
		// arrange
		aggID := testdata.StringIdentifier("TestAgg")
		agg := aggregate.NewBase(testdata.NewTestAggregate(aggID), nil, nil)

		eventStore := &testfixture.EventStoreMock{}
		aggFactory := &testfixture.AggregateFactoryMock{
			Creator: func(aggregateType string, ID domain.Identifier) domain.AdvancedAggregate {
				// assert
				assert.Equals(t, aggID, ID)
				assert.Equals(t, agg.AggregateType(), aggregateType)
				return agg
			},
		}

		d := dispatcher.NewDispatcher(eventStore, aggFactory)

		// act
		_, err := d.Handle(testdata.MakeSomethingHappen{AggID: aggID})

		// assert
		assert.Ok(t, err)
	})
}
