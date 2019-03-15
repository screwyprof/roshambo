package command

type CreateNewGame struct {
	GameID string
}

func (c CreateNewGame) CommandType() string {
	return "CreateNewGame"
}
