package volume

// Hop represents part of the journey
type Hop struct {
	Start string
	End   string
}

// Reduce merges list of hops into a single hop containing places where journey was started and finished
func Reduce(hops []Hop) (Hop, error) {
	panic("not implemented")
}
