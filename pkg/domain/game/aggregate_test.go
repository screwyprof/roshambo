package game_test

import (
	"testing"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"
	. "github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testfixture"

	"github.com/screwyprof/roshambo/pkg/command"
	"github.com/screwyprof/roshambo/pkg/domain"
	"github.com/screwyprof/roshambo/pkg/domain/game"
	"github.com/screwyprof/roshambo/pkg/event"
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

func TestAggregate_CreateNewGame(t *testing.T) {
	t.Run("ItCreatesNewGame", func(t *testing.T) {
		ID := testdata.StringIdentifier("g777")
		Test(t)(
			Given(createTestAggregate()),
			When(command.CreateNewGame{GameID: ID.String()}),
			Then(event.GameCreated{GameID: ID.String()}),
		)
	})

	t.Run("ItCannotStartANewGameIfItTheGameIsAlreadyStarted", func(t *testing.T) {
		ID := testdata.StringIdentifier("g777")
		Test(t)(
			Given(createTestAggregate(), event.GameCreated{GameID: ID.String()}),
			When(command.CreateNewGame{GameID: ID.String()}),
			ThenFailWith(game.ErrGameIsAlreadyStarted),
		)
	})
}

func createTestAggregate() *aggregate.Base {
	ID := testdata.StringIdentifier("GameAgg")
	gameAgg := game.NewAggregate(ID)

	return aggregate.NewBase(gameAgg, nil, nil)
}
