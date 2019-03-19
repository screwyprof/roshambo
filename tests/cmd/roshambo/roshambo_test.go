package roshambo_test

import (
	"fmt"

	"os"

	"github.com/screwyprof/roshambo/internal/pkg/cqrs/aggregate"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/dispatcher"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventbus"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventhandler"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/eventstore"
	"github.com/screwyprof/roshambo/internal/pkg/cqrs/testdata/mock"

	"github.com/screwyprof/roshambo/pkg/command"
	"github.com/screwyprof/roshambo/pkg/domain"
	"github.com/screwyprof/roshambo/pkg/domain/game"
	gameEventHandler "github.com/screwyprof/roshambo/pkg/eventhandler"
	"github.com/screwyprof/roshambo/pkg/report"
)

func Example() {
	ID := mock.StringIdentifier("TestGame")
	gameInfo := report.GameShortInfo{}
	d := createDispatcher(&gameInfo)

	failOnError(d.Handle(command.CreateNewGame{GameID: ID}))
	failOnError(d.Handle(command.MakeMove{GameID: ID, PlayerEmail: "gopher@happy", Move: int(game.Rock)}))
	failOnError(d.Handle(command.MakeMove{GameID: ID, PlayerEmail: "tiger@happy", Move: int(game.Scissors)}))

	printGameInfo(gameInfo)
	// Output:
	// Game Info
	// Status: game won
	// Creator: tiger@happy
	// Winner: gopher@happy
}

func printGameInfo(gameInfo report.GameShortInfo) {
	fmt.Println("Game Info")
	fmt.Println("Status:", gameInfo.State)
	fmt.Println("Creator:", gameInfo.Creator)
	fmt.Println("Winner:", gameInfo.Winner)
}

func createDispatcher(gameInfo *report.GameShortInfo) *dispatcher.Dispatcher {
	f := aggregate.NewFactory()
	f.RegisterAggregate(func(ID domain.Identifier) domain.AdvancedAggregate {
		return aggregate.NewBase(game.NewAggregate(ID), nil, nil)
	})

	gameInfoProjector := eventhandler.NewDynamic()
	gameInfoProjector.RegisterHandlers(&gameEventHandler.GameShortInfoProjector{Projection: gameInfo})

	eventBus := eventbus.NewInMemoryEventBus()
	eventBus.Register(gameInfoProjector)

	return dispatcher.NewDispatcher(eventstore.NewInInMemoryEventStore(), f, eventBus)
}

func failOnError(_ []domain.DomainEvent, err error) {
	if err != nil {
		fmt.Printf("an error occured: %v", err)
		os.Exit(1)
	}
}
