package event

type MoveDecided struct {
	GameID      string
	PlayerEmail string
	Move        int
}

func (c MoveDecided) EventType() string {
	return "MoveDecided"
}
