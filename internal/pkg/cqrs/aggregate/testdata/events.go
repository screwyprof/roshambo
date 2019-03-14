package testdata

type SomethingHappened struct{}
func (c SomethingHappened) EventType() string {
	return "SomethingHappened"
}