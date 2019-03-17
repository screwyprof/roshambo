package game

import (
	"testing"

	"github.com/segmentio/ksuid"

	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/dispatcher"
	. "github.com/screwyprof/roshambo/internal/pkg/cqrs/dispatcher/testdata/fixture"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventstore"

	"github.com/screwyprof/roshambo/pkg/command"
	"github.com/screwyprof/roshambo/pkg/domain"
	"github.com/screwyprof/roshambo/pkg/domain/game"
	"github.com/screwyprof/roshambo/pkg/event"
)

func TestVictory(t *testing.T) {
	ID := ksuid.New()
	player1 := "tom@game.net"
	player2 := "jerry@game.net"

	Test(t)(
		Given(createDispatcher()),
		When(
			command.CreateNewGame{GameID: ID},
			command.MakeMove{GameID: ID, PlayerEmail: player1, Move: int(game.Rock)},
			command.MakeMove{GameID: ID, PlayerEmail: player2, Move: int(game.Paper)},
		),
		Then(
			event.GameCreated{GameID: ID.String()},
			event.MoveDecided{GameID: ID.String(), PlayerEmail: player1, Move: int(game.Rock)},
			event.MoveDecided{GameID: ID.String(), PlayerEmail: player2, Move: int(game.Paper)},
			event.GameWon{GameID: ID.String(), Winner: player2, Loser: player1},
		),
	)
}

func TestTie(t *testing.T) {
	ID := ksuid.New()
	player1 := "tom@game.net"
	player2 := "jerry@game.net"

	Test(t)(
		Given(createDispatcher()),
		When(
			command.CreateNewGame{GameID: ID},
			command.MakeMove{GameID: ID, PlayerEmail: player1, Move: int(game.Scissors)},
			command.MakeMove{GameID: ID, PlayerEmail: player2, Move: int(game.Scissors)},
		),
		Then(
			event.GameCreated{GameID: ID.String()},
			event.MoveDecided{GameID: ID.String(), PlayerEmail: player1, Move: int(game.Scissors)},
			event.MoveDecided{GameID: ID.String(), PlayerEmail: player2, Move: int(game.Scissors)},
			event.GameTied{GameID: ID.String()},
		),
	)
}

func createDispatcher() *dispatcher.Dispatcher {
	eventStore := eventstore.NewInInMemoryEventStore()
	f := aggregate.NewFactory()
	f.RegisterAggregate(func(ID domain.Identifier) domain.AdvancedAggregate {
		return aggregate.NewBase(game.NewAggregate(ID), nil, nil)
	})

	return dispatcher.NewDispatcher(eventStore, f)
}
