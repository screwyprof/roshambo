package event

type GameCreated struct {
	GameID  string
	Creator string
}

func (c GameCreated) EventType() string {
	return "GameCreated"
}
