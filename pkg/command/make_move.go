package command

type MakeMove struct {
	GameID      string
	PlayerEmail string
	Move        int
}

func (c MakeMove) CommandType() string {
	return "MakeMove"
}
