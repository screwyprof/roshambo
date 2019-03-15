package game_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"

	"github.com/screwyprof/roshambo/pkg/domain"
	"github.com/screwyprof/roshambo/pkg/domain/game"
)

// ensure that game aggregate implements domain.Aggregate interface.
var _ domain.Aggregate = (*game.Aggregate)(nil)

func TestNewAggregate(t *testing.T) {
	t.Run("ItPanicsIfIDIsNotGiven", func(t *testing.T) {
		factory := func() {
			game.NewAggregate(nil)
		}
		assert.Panic(t, factory)
	})
}

func TestAggregateAggregateID(t *testing.T) {
	t.Run("ItReturnsAggregateID", func(t *testing.T) {
		ID := testdata.StringIdentifier("Game")
		agg := game.NewAggregate(ID)

		assert.Equals(t, ID, agg.AggregateID())
	})
}
