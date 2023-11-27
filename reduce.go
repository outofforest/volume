package volume

import (
	"errors"
	"fmt"
)

// Hop represents part of the journey
type Hop struct {
	Start string
	End   string
}

func (h Hop) String() string {
	return fmt.Sprintf("(%s, %s)", h.Start, h.End)
}

// Reduce merges list of hops into a single hop containing places where journey was started and finished
func Reduce(hops []Hop) (Hop, error) {
	counters := map[string]int{}
	for _, h := range hops {
		if h.Start == "" {
			return Hop{}, fmt.Errorf("no starting point defined for hop %s", h)
		}
		if h.End == "" {
			return Hop{}, fmt.Errorf("no finishing point defined for hop %s", h)
		}
		counters[h.Start]--
		counters[h.End]++
	}
	if len(counters) == 1 {
		// AAA -> AAA
		// AAA -> AAA -> AAA
		for symbol := range counters {
			return Hop{Start: symbol, End: symbol}, nil
		}
	}

	var res Hop
	for symbol, c := range counters {
		switch c {
		case 0:
			// intermediary vertex
			continue
		case -1:
			// starting point
			if res.Start != "" {
				return Hop{}, fmt.Errorf("more than one starting point detected: %s", symbol)
			}
			res.Start = symbol
		case 1:
			// finishing point

			// Checking if we have already found an end point before is not needed here because it's not possible
			// to create many of them without creating many starting points too.

			res.End = symbol
		default:
			return Hop{}, fmt.Errorf("invalid point detected: %s", symbol)
		}
	}

	// Checking if end point is empty is not needed here because in such case starting point would be empty too.
	if res.Start == "" {
		return Hop{}, errors.New("invalid list of hops 2")
	}
	return res, nil
}
