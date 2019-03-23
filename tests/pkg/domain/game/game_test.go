package game

import (
	"testing"

	"github.com/segmentio/ksuid"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/dispatcher"
	. "github.com/screwyprof/roshambo/internal/pkg/cqrs/dispatcher/testdata/fixture"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventbus"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventhandler"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventstore"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/store"

	"github.com/screwyprof/roshambo/pkg/command"
	"github.com/screwyprof/roshambo/pkg/domain"
	"github.com/screwyprof/roshambo/pkg/domain/game"
	"github.com/screwyprof/roshambo/pkg/event"
	gameEventHandler "github.com/screwyprof/roshambo/pkg/eventhandler"
	"github.com/screwyprof/roshambo/pkg/report"
)

func TestVictory(t *testing.T) {
	ID := ksuid.New()
	player1 := "tom@game.net"
	player2 := "jerry@game.net"

	got := report.GameShortInfo{}
	want := report.GameShortInfo{
		GameID:  ID.String(),
		Creator: player1,
		State:   "game won",
		Winner:  player2,
		Loser:   player1,
	}

	Test(t)(
		Given(createDispatcher(&got)),
		When(
			command.CreateNewGame{GameID: ID, Creator: player1},
			command.MakeMove{GameID: ID, PlayerEmail: player1, Move: int(game.Rock)},
			command.MakeMove{GameID: ID, PlayerEmail: player2, Move: int(game.Paper)},
		),
		Then(
			event.GameCreated{GameID: ID.String(), Creator: player1},
			event.MoveDecided{GameID: ID.String(), PlayerEmail: player1, Move: int(game.Rock)},
			event.MoveDecided{GameID: ID.String(), PlayerEmail: player2, Move: int(game.Paper)},
			event.GameWon{GameID: ID.String(), Winner: player2, Loser: player1},
		),
	)

	assert.Equals(t, want, got)
}

func TestTie(t *testing.T) {
	ID := ksuid.New()
	player1 := "tom@game.net"
	player2 := "jerry@game.net"

	got := report.GameShortInfo{}
	want := report.GameShortInfo{
		GameID:  ID.String(),
		Creator: player2,
		State:   "game tied",
	}

	Test(t)(
		Given(createDispatcher(&got)),
		When(
			command.CreateNewGame{GameID: ID, Creator: player2},
			command.MakeMove{GameID: ID, PlayerEmail: player1, Move: int(game.Scissors)},
			command.MakeMove{GameID: ID, PlayerEmail: player2, Move: int(game.Scissors)},
		),
		Then(
			event.GameCreated{GameID: ID.String(), Creator: player2},
			event.MoveDecided{GameID: ID.String(), PlayerEmail: player1, Move: int(game.Scissors)},
			event.MoveDecided{GameID: ID.String(), PlayerEmail: player2, Move: int(game.Scissors)},
			event.GameTied{GameID: ID.String()},
		),
	)

	assert.Equals(t, want, got)
}

func createDispatcher(gameInfo *report.GameShortInfo) *dispatcher.Dispatcher {
	gameInfoProjector := eventhandler.NewDynamic()
	gameInfoProjector.RegisterHandlers(&gameEventHandler.GameShortInfoProjector{Projection: gameInfo})

	f := aggregate.NewFactory()
	f.RegisterAggregate(func(ID domain.Identifier) domain.AdvancedAggregate {
		return aggregate.NewBase(game.NewAggregate(ID), nil, nil)
	})

	aggregateStore := store.NewStore(eventstore.NewInInMemoryEventStore(), f)
	eventBus := eventbus.NewInMemoryEventBus()
	eventBus.Register(gameInfoProjector)

	return dispatcher.NewDispatcher(aggregateStore, eventBus)
}
