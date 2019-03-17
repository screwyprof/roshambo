package game_test

import (
	"fmt"

	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/roshambo/pkg/command"
	"github.com/screwyprof/roshambo/pkg/domain"
	"github.com/screwyprof/roshambo/pkg/domain/game"
)

func Example_Victory() {
	suite := newExampleSuite("TheGame")
	suite.run(func(ID domain.Identifier, agg domain.AdvancedAggregate, player1, player2 string) {
		agg.Handle(command.CreateNewGame{GameID: ID})
		agg.Handle(command.MakeMove{GameID: ID, PlayerEmail: player1, Move: int(game.Rock)})
		events, _ := agg.Handle(command.MakeMove{GameID: ID, PlayerEmail: player2, Move: int(game.Paper)})

		fmt.Printf("%#v", events[1])
	})

	// Output:
	// event.GameWon{GameID:"TheGame", Winner:"jerry@game.net", Loser:"tom@game.net"}
}

func Example_Tie() {
	suite := newExampleSuite("TieGame")
	suite.run(func(ID domain.Identifier, agg domain.AdvancedAggregate, player1, player2 string) {
		agg.Handle(command.CreateNewGame{GameID: ID})
		agg.Handle(command.MakeMove{GameID: ID, PlayerEmail: player1, Move: int(game.Rock)})
		events, _ := agg.Handle(command.MakeMove{GameID: ID, PlayerEmail: player2, Move: int(game.Rock)})

		fmt.Printf("%#v", events[1])
	})

	// Output:
	// event.GameTied{GameID:"TieGame"}
}

type exampleRunner func(ID domain.Identifier, gameAggregate domain.AdvancedAggregate, player1, player2 string)

type exampleSuite struct {
	ID   domain.Identifier
	game domain.AdvancedAggregate

	player1 string
	player2 string
}

func newExampleSuite(ID string) exampleSuite {
	id := mock.StringIdentifier(ID)
	return exampleSuite{
		ID:      id,
		game:    aggregate.NewBase(game.NewAggregate(id), nil, nil),
		player1: "tom@game.net",
		player2: "jerry@game.net",
	}
}

func (s exampleSuite) run(example exampleRunner) {
	example(s.ID, s.game, s.player1, s.player2)
}
