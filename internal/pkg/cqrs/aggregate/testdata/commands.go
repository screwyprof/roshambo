package testdata

type MakeSomethingHappen struct{}
func (c MakeSomethingHappen) CommandType() string {
	return "MakeSomethingHappen"
}