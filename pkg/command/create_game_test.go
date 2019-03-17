package command_test

import (
	"testing"

	"github.com/segmentio/ksuid"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/pkg/command"
)

func TestCreateNewGameAggregateID(t *testing.T) {
	ID := ksuid.New()
	assert.Equals(t, ID, command.CreateNewGame{GameID: ID}.AggregateID())
}

func TestCreateNewGameAggregateType(t *testing.T) {
	assert.Equals(t, "game.Aggregate", command.CreateNewGame{}.AggregateType())
}

func TestCreateNewGameCommandType(t *testing.T) {
	assert.Equals(t, "CreateNewGame", command.CreateNewGame{}.CommandType())
}
