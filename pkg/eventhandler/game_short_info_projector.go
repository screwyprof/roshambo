package eventhandler

import (
	"github.com/screwyprof/roshambo/pkg/event"
	"github.com/screwyprof/roshambo/pkg/report"
)

type GameShortInfoProjector struct {
	Projection *report.GameShortInfo
}

func (p *GameShortInfoProjector) OnGameCreated(e event.GameCreated) error {
	p.Projection.GameID = e.GameID
	p.Projection.Creator = e.Creator
	p.Projection.State = "created"
	return nil
}

func (p *GameShortInfoProjector) OnGameWon(e event.GameWon) error {
	p.Projection.State = "game won"
	p.Projection.Winner = e.Winner
	p.Projection.Loser = e.Loser

	return nil
}

func (p *GameShortInfoProjector) OnGameTied(e event.GameTied) error {
	p.Projection.State = "game tied"
	return nil
}
