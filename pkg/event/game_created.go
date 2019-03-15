package event

type GameCreated struct {
	GameID string
}

func (c GameCreated) EventType() string {
	return "GameCreated"
}
