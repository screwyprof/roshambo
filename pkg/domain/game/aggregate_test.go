package game_test

import (
	"testing"

	"github.com/segmentio/ksuid"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	. "github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata/fixture"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/testdata/mock"

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
		ID := mock.StringIdentifier("Game")
		agg := game.NewAggregate(ID)

		assert.Equals(t, ID, agg.AggregateID())
	})
}

func TestAggregateAggregateType(t *testing.T) {
	t.Run("ItReturnsAggregateType", func(t *testing.T) {
		ID := mock.StringIdentifier("Game")
		agg := game.NewAggregate(ID)

		assert.Equals(t, "game.Aggregate", agg.AggregateType())
	})
}

func TestAggregate_CreateNewGame(t *testing.T) {
	t.Run("ItCreatesNewGame", func(t *testing.T) {
		ID := mock.StringIdentifier("g777")
		Test(t)(
			Given(createTestAggregate()),
			When(command.CreateNewGame{GameID: ID, Creator: "tiger@happy.com"}),
			Then(event.GameCreated{GameID: ID.String(), Creator: "tiger@happy.com"}),
		)
	})

	t.Run("ItCannotStartANewGameIfItTheGameIsAlreadyStarted", func(t *testing.T) {
		ID := mock.StringIdentifier("g777")
		Test(t)(
			Given(createTestAggregate(), event.GameCreated{GameID: ID.String()}),
			When(command.CreateNewGame{GameID: ID}),
			ThenFailWith(game.ErrGameIsAlreadyStarted),
		)
	})
}

func TestAggregateMakeMove(t *testing.T) {
	t.Run("APlayerCanMakeAMove", func(t *testing.T) {
		ID := mock.StringIdentifier("g777")
		Test(t)(
			Given(createTestAggregate(), event.GameCreated{GameID: ID.String()}),
			When(command.MakeMove{GameID: ID, PlayerEmail: "player@game.com", Move: int(game.Rock)}),
			Then(event.MoveDecided{GameID: ID.String(), PlayerEmail: "player@game.com", Move: int(game.Rock)}),
		)
	})

	t.Run("ItFailsIfThePlayerIsTheSame", func(t *testing.T) {
		ID := mock.StringIdentifier("g777")
		Test(t)(
			Given(createTestAggregate(),
				event.GameCreated{GameID: ID.String()},
				event.MoveDecided{GameID: ID.String(), PlayerEmail: "player@game.com", Move: int(game.Rock)}),
			When(command.MakeMove{GameID: ID, PlayerEmail: "player@game.com", Move: int(game.Rock)}),
			ThenFailWith(game.ErrPlayerIsTheSame),
		)
	})

	t.Run("ItFailsIfTheGameHaveNotStarted", func(t *testing.T) {
		ID := mock.StringIdentifier("g777")
		Test(t)(
			Given(createTestAggregate()),
			When(command.MakeMove{GameID: ID, PlayerEmail: "player@game.com", Move: int(game.Rock)}),
			ThenFailWith(game.ErrTheGameHaveNotStartedOrFinished),
		)
	})

	t.Run("FirstPlayerDeclaredAWinner", func(t *testing.T) {
		ID := mock.StringIdentifier("g777")
		Test(t)(
			Given(createTestAggregate(),
				event.GameCreated{GameID: ID.String()},
				event.MoveDecided{GameID: ID.String(), PlayerEmail: "player1@game.com", Move: int(game.Scissors)}),
			When(command.MakeMove{GameID: ID, PlayerEmail: "player2@game.com", Move: int(game.Paper)}),
			Then(
				event.MoveDecided{GameID: ID.String(), PlayerEmail: "player2@game.com", Move: int(game.Paper)},
				event.GameWon{GameID: ID.String(), Winner: "player1@game.com", Loser: "player2@game.com"},
			),
		)
	})

	t.Run("ItFailsIfTheMoveIsMadeAfterTheGameIsFinished", func(t *testing.T) {
		ID := mock.StringIdentifier("g777")
		Test(t)(
			Given(createTestAggregate(),
				event.GameCreated{GameID: ID.String()},
				event.MoveDecided{GameID: ID.String(), PlayerEmail: "player1@game.com", Move: int(game.Rock)},
				event.MoveDecided{GameID: ID.String(), PlayerEmail: "player2@game.com", Move: int(game.Rock)},
				event.GameTied{GameID: ID.String()}),
			When(command.MakeMove{GameID: ID, PlayerEmail: "another@game.com", Move: int(game.Paper)}),
			ThenFailWith(game.ErrTheGameHaveNotStartedOrFinished),
		)
	})

	t.Run("SecondPlayerDeclaredAWinner", func(t *testing.T) {
		ID := mock.StringIdentifier("g777")
		Test(t)(
			Given(createTestAggregate(),
				event.GameCreated{GameID: ID.String()},
				event.MoveDecided{GameID: ID.String(), PlayerEmail: "player1@game.com", Move: int(game.Rock)}),
			When(command.MakeMove{GameID: ID, PlayerEmail: "player2@game.com", Move: int(game.Paper)}),
			Then(
				event.MoveDecided{GameID: ID.String(), PlayerEmail: "player2@game.com", Move: int(game.Paper)},
				event.GameWon{GameID: ID.String(), Winner: "player2@game.com", Loser: "player1@game.com"},
			),
		)
	})

	t.Run("GameTied", func(t *testing.T) {
		ID := mock.StringIdentifier("g777")
		Test(t)(
			Given(createTestAggregate(),
				event.GameCreated{GameID: ID.String()},
				event.MoveDecided{GameID: ID.String(), PlayerEmail: "player1@game.com", Move: int(game.Scissors)}),
			When(command.MakeMove{GameID: ID, PlayerEmail: "player2@game.com", Move: int(game.Scissors)}),
			Then(
				event.MoveDecided{GameID: ID.String(), PlayerEmail: "player2@game.com", Move: int(game.Scissors)},
				event.GameTied{GameID: ID.String()},
			),
		)
	})
}

func createTestAggregate() *aggregate.Base {
	gameAgg := game.NewAggregate(ksuid.New())

	commandHandler := aggregate.NewCommandHandler()
	commandHandler.RegisterHandlers(gameAgg)

	eventApplier := aggregate.NewEventApplier()
	eventApplier.RegisterAppliers(gameAgg)

	return aggregate.NewBase(gameAgg, commandHandler, eventApplier)
}
