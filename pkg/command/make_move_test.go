package command_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/screwyprof/roshambo/internal/pkg/assert"
	"github.com/screwyprof/roshambo/pkg/command"
)

func TestMakeMoveAggregateID(t *testing.T) {
	ID := uuid.New()
	assert.Equals(t, ID, command.MakeMove{GameID: ID}.AggregateID())
}

func TestMakeMoveAggregateType(t *testing.T) {
	assert.Equals(t, "game.Aggregate", command.MakeMove{}.AggregateType())
}

func TestMakeMoveCommandType(t *testing.T) {
	assert.Equals(t, "MakeMove", command.MakeMove{}.CommandType())
}
