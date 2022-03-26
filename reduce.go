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

// Reduce merges list of hops into a single hop containing places where journey was started and finished
func Reduce(hops []Hop) (Hop, error) {
	counters := map[string]edgeCounter{}
	for _, h := range hops {
		hStart := counters[h.Start]
		hStart.Outgoing++
		counters[h.Start] = hStart

		hEnd := counters[h.End]
		hEnd.Incoming++
		counters[h.End] = hEnd
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
		switch {
		case c.Outgoing == c.Incoming:
			// intermediary vertex
			continue
		case c.Outgoing == c.Incoming+1:
			// starting point
			if res.Start != "" {
				return Hop{}, fmt.Errorf("more than one starting point detected: %s", symbol)
			}
			res.Start = symbol
		case c.Incoming == c.Outgoing+1:
			// finishing point

			// Checking if we have already found an end point before is not needed here because it's not possible
			// to create many of them without creating many starting points too.

			res.End = symbol
		default:
			return Hop{}, fmt.Errorf("invalid point detected: %s", symbol)
		}
	}

	// Checking if end point s empty is not needed here because in such case starting point would be empty too.
	if res.Start == "" {
		return Hop{}, errors.New("invalid list of hops")
	}
	return res, nil
}

type edgeCounter struct {
	Incoming uint64
	Outgoing uint64
}
