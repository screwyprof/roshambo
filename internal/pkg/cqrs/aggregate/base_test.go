package aggregate_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"
)

func TestNewBase(t *testing.T) {
	t.Run("ItPanicsIfThePureAggregateIsNotGiven", func(t *testing.T) {
		factory := func() {
			aggregate.NewBase(nil)
		}
		assert.Panic(t, factory)
	})

	t.Run("ItReturnsAggregateID", func(t *testing.T) {
		ID := testdata.StringIdentifier("TestAgg1")
		pureAgg := testdata.NewTestAggregate(ID)
		agg := aggregate.NewBase(pureAgg)

		assert.Equals(t, ID, agg.AggregateID())
	})
}
