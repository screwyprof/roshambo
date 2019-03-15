package event

type GameTied struct {
	GameID string
}

func (c GameTied) EventType() string {
	return "GameTied"
}
