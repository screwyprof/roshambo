package game_test

import (
	"fmt"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"

	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate/testdata"

	"github.com/screwyprof/roshambo/pkg/command"
	"github.com/screwyprof/roshambo/pkg/domain"
	"github.com/screwyprof/roshambo/pkg/domain/game"
)

var (
	ID   domain.Identifier
	Game domain.AdvancedAggregate

	Player1 string
	Player2 string
)

func init() {
	ID = testdata.StringIdentifier("TheGame")

	Player1 = "tom@game.net"
	Player2 = "jerry@game.net"
}

func ExampleTie() {
	Game = aggregate.NewBase(game.NewAggregate(ID), nil, nil)

	Game.Handle(command.CreateNewGame{GameID: ID.String()})
	Game.Handle(command.MakeMove{GameID: ID.String(), PlayerEmail: Player1, Move: int(game.Rock)})
	events, _ := Game.Handle(command.MakeMove{GameID: ID.String(), PlayerEmail: Player2, Move: int(game.Rock)})

	fmt.Printf("%#v", events[1])

	// Output:
	// event.GameTied{GameID:"TheGame"}
}

func ExampleVictory() {
	Game = aggregate.NewBase(game.NewAggregate(ID), nil, nil)

	Game.Handle(command.CreateNewGame{GameID: ID.String()})
	Game.Handle(command.MakeMove{GameID: ID.String(), PlayerEmail: Player1, Move: int(game.Rock)})
	events, _ := Game.Handle(command.MakeMove{GameID: ID.String(), PlayerEmail: Player2, Move: int(game.Paper)})

	fmt.Printf("%#v", events[1])

	// Output:
	// event.GameWon{GameID:"TheGame", Winner:"jerry@game.net", Loser:"tom@game.net"}
}
