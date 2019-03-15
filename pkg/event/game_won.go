package event

type GameWon struct {
	GameID string
	Winner string
	Loser  string
}

func (c GameWon) EventType() string {
	return "GameWon"
}
