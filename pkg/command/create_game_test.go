package command_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/pkg/command"
)

func TestCreateNewGameAggregateID(t *testing.T) {
	ID := uuid.New()
	assert.Equals(t, ID, command.CreateNewGame{GameID: ID}.AggregateID())
}

func TestCreateNewGameAggregateType(t *testing.T) {
	assert.Equals(t, "game.Aggregate", command.CreateNewGame{}.AggregateType())
}

func TestCreateNewGameCommandType(t *testing.T) {
	assert.Equals(t, "CreateNewGame", command.CreateNewGame{}.CommandType())
}
