package game

type Move int

const (
	Rock Move = iota
	Paper
	Scissors
)

func NewMove(m int) Move {
	return Move(m)
}

func (m Move) defeats(other Move) bool {
	switch m {
	case Rock:
		return other == Scissors
	case Paper:
		return other == Rock
	case Scissors:
		return other == Paper
	default:
		return false
	}
}
