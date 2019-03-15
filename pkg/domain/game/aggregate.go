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
)

var (
	ErrGameIsAlreadyStarted = errors.New("game is already started")
)

type Aggregate struct {
	id    domain.Identifier
	state state
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

func (a *Aggregate) CreateNewGame(c command.CreateNewGame) ([]domain.DomainEvent, error) {
	if a.state != notCreated {
		return nil, ErrGameIsAlreadyStarted
	}
	return []domain.DomainEvent{event.GameCreated{GameID: c.GameID}}, nil
}

func (a *Aggregate) OnGameCreated(e event.GameCreated) {
	a.state = created
}
