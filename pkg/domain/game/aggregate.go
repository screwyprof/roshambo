package game

import (
	"errors"
	"github.com/screwyprof/roshambo/pkg/command"
	"github.com/screwyprof/roshambo/pkg/domain"
	"github.com/screwyprof/roshambo/pkg/event"
)

type state int

const (
	notCreated state = iota
	created
	waiting
	tied
	won
)

var (
	ErrGameIsAlreadyStarted = errors.New("game is already started")
	ErrPlayerIsTheSame      = errors.New("the player is already in the game")
)

type Aggregate struct {
	id domain.Identifier

	state       state
	playerEmail string
	move        Move
}

// NewAggregate creates a new instance of Aggregate.
func NewAggregate(ID domain.Identifier) *Aggregate {
	if ID == nil {
		panic("ID is required")
	}
	return &Aggregate{id: ID, state: notCreated}
}

// AggregateID implements domain.Aggregate interface.
func (a *Aggregate) AggregateID() domain.Identifier {
	return a.id
}

// CreateNewGame starts a new game.
// If the game has already started then returns an error.
func (a *Aggregate) CreateNewGame(c command.CreateNewGame) ([]domain.DomainEvent, error) {
	if a.state != notCreated {
		return nil, ErrGameIsAlreadyStarted
	}
	return []domain.DomainEvent{event.GameCreated{GameID: c.GameID}}, nil
}

// MakeMove makes a move.
// When the second player has moved, the game is finished with a tie or a win.
func (a *Aggregate) MakeMove(c command.MakeMove) ([]domain.DomainEvent, error) {
	if a.playerEmail == c.PlayerEmail {
		return nil, ErrPlayerIsTheSame
	}

	if a.state == waiting {
		return []domain.DomainEvent{
			event.MoveDecided{GameID: c.GameID, PlayerEmail: c.PlayerEmail, Move: c.Move},
			a.finish(c.GameID, c.PlayerEmail, NewMove(c.Move)),
		}, nil
	}

	return []domain.DomainEvent{event.MoveDecided{GameID: c.GameID, PlayerEmail: c.PlayerEmail, Move: c.Move}}, nil
}

func (a *Aggregate) OnGameCreated(e event.GameCreated) {
	a.state = created
}

func (a *Aggregate) OnMoveDecided(e event.MoveDecided) {
	a.playerEmail = e.PlayerEmail
	a.move = Move(e.Move)
	a.state = waiting
}

func (a *Aggregate) OnGameWon(e event.GameWon) {
	a.state = won
}

func (a *Aggregate) OnGameTied(e event.GameTied) {
	a.state = tied
}

func (a *Aggregate) finish(gameID string, opponentEmail string, opponentMove Move) domain.DomainEvent {
	switch {
	case a.move.defeats(opponentMove):
		return event.GameWon{GameID: gameID, Winner: a.playerEmail, Loser: opponentEmail}
	case opponentMove.defeats(a.move):
		return event.GameWon{GameID: gameID, Winner: opponentEmail, Loser: a.playerEmail}
	default:
		return event.GameTied{GameID: gameID}
	}
}
