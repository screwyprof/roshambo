package command

import "github.com/screwyprof/roshambo/pkg/domain"

type MakeMove struct {
	GameID      domain.Identifier
	PlayerEmail string
	Move        int
}

func (c MakeMove) AggregateID() domain.Identifier {
	return c.GameID
}

func (c MakeMove) AggregateType() string {
	return "game.Aggregate"
}

func (c MakeMove) CommandType() string {
	return "MakeMove"
}
