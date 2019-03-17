package command

import "github.com/screwyprof/roshambo/pkg/domain"

type CreateNewGame struct {
	GameID domain.Identifier
}

func (c CreateNewGame) AggregateID() domain.Identifier {
	return c.GameID
}

func (c CreateNewGame) AggregateType() string {
	return "game.Aggregate"
}

func (c CreateNewGame) CommandType() string {
	return "CreateNewGame"
}
