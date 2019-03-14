package aggregate_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"
)

func TestBaseAggregateID(t *testing.T) {
	t.Run("ItReturnsAggregateID", func(t *testing.T) {
		ID := testdata.StringIdentifier("TestAgg1")
		pureAgg := testdata.NewTestAggregate(ID)
		agg := aggregate.NewBase(pureAgg)

		assert.Equals(t, ID, agg.AggregateID())
	})
}
