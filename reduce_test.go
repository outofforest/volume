package volume

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tCases = []tCase{
	{
		Name:  "TestNilJourney",
		Fails: true,
	},
	{
		Name:  "TestEmptyJourney",
		Fails: true,
		Hops:  []Hop{},
	},
	{
		Name: "TestSingleHopJourney",
		Hops: []Hop{
			{Start: "AAA", End: "BBB"},
		},
		Result: Hop{"AAA", "BBB"},
	},
	{
		Name: "TestTwoHops",
		Hops: []Hop{
			{Start: "BBB", End: "CCC"},
			{Start: "AAA", End: "BBB"},
		},
		Result: Hop{"AAA", "CCC"},
	},
	{
		Name: "TestLongJourney",
		Hops: []Hop{
			// AAA -> BBB -> CCC -> DDD -> EEE -> FFF -> GGG -> HHH -> III
			{Start: "FFF", End: "GGG"},
			{Start: "GGG", End: "HHH"},
			{Start: "BBB", End: "CCC"},
			{Start: "EEE", End: "FFF"},
			{Start: "HHH", End: "III"},
			{Start: "AAA", End: "BBB"},
			{Start: "DDD", End: "EEE"},
			{Start: "CCC", End: "DDD"},
		},
		Result: Hop{"AAA", "III"},
	},
	{
		Name:  "TestDisconnectedParts",
		Fails: true,
		Hops: []Hop{
			// AAA -> BBB -> CCC -> <missing CCC -> DDD> -> EEE -> FFF
			{Start: "BBB", End: "CCC"},
			{Start: "AAA", End: "BBB"},
			// CCC -> DDD is missing
			{Start: "DDD", End: "EEE"},
			{Start: "EEE", End: "FFF"},
		},
	},
	{
		Name:  "TestFork",
		Fails: true,
		Hops: []Hop{
			{Start: "AAA", End: "BBB"},
			{Start: "1CC", End: "1DD"},
			{Start: "2CC", End: "2DD"},
			{Start: "2DD", End: "2EE"},

			// At point BBB route forks and two finish points are possible
			{Start: "BBB", End: "1CC"},
			{Start: "BBB", End: "2CC"},
		},
	},
	{
		Name: "TestSingleHopLoop",
		Hops: []Hop{
			{Start: "AAA", End: "AAA"},
		},
		Result: Hop{Start: "AAA", End: "AAA"},
	},
	{
		Name:  "TestLoop",
		Fails: true,
		Hops: []Hop{
			// In this case it's not possible to tell where journey started and finished
			{Start: "CCC", End: "AAA"},
			{Start: "AAA", End: "BBB"},
			{Start: "BBB", End: "CCC"},
		},
	},
	{
		Name: "TestLoopAtTheEnd",
		Hops: []Hop{
			{Start: "DDD", End: "BBB"},
			{Start: "AAA", End: "BBB"},
			{Start: "BBB", End: "CCC"},
			{Start: "CCC", End: "DDD"},
		},
		Result: Hop{Start: "AAA", End: "BBB"},
	},
	{
		Name:  "TestForkBeforeLoop",
		Fails: true,
		Hops: []Hop{
			// Two paths are possible:
			// - AAA -> BBB -> CCC -> DDD -> EEE
			// - AAA -> BBB -> CCC -> DDD -> BBB
			{Start: "DDD", End: "BBB"},
			{Start: "AAA", End: "BBB"},
			{Start: "BBB", End: "CCC"},
			{Start: "CCC", End: "DDD"},
			{Start: "DDD", End: "EEE"},
		},
	},
	{
		Name: "TestForkAfterLoop",
		Hops: []Hop{
			// There are two possible paths to EEE:
			// - AAA -> BBB -> EEE
			// - AAA -> BBB -> CCC -> DDD -> BBB -> EEE (BBB here occurs twice)
			// but only the second one is correct because it uses all the hops.
			// Anyway, result is the same in both cases.
			{Start: "DDD", End: "BBB"},
			{Start: "AAA", End: "BBB"},
			{Start: "BBB", End: "CCC"},
			{Start: "CCC", End: "DDD"},
			{Start: "BBB", End: "EEE"},
		},
		Result: Hop{Start: "AAA", End: "EEE"},
	},
	{
		Name: "TestLongerLoop",
		Hops: []Hop{
			// There are two possible paths to EEE:
			// - AAA -> BBB -> CCC -> EEE
			// - AAA -> BBB -> CCC -> DDD -> BBB -> CCC -> EEE (BBB and CCC here occur twice)
			// but only the second one is correct because it uses all the hops.
			// Anyway, result is the same in both cases.
			{Start: "AAA", End: "BBB"},
			{Start: "BBB", End: "CCC"},
			{Start: "BBB", End: "CCC"}, // to make this journey possible BBB -> CCC pair has to occur twice
			{Start: "CCC", End: "EEE"},
			{Start: "CCC", End: "DDD"},
			{Start: "DDD", End: "BBB"},
		},
		Result: Hop{Start: "AAA", End: "EEE"},
	},
	{
		Name:  "TestDisconnectionAfterLoop",
		Fails: true,
		Hops: []Hop{
			// There are two possible paths:
			// - AAA -> BBB -> CCC -> EEE
			// - AAA -> BBB -> CCC -> DDD -> BBB
			// Solution can't be found.
			{Start: "AAA", End: "BBB"},
			{Start: "BBB", End: "CCC"}, // now there is only one BBB -> CCC pair so after the loop, continuation to EEE is not possible
			{Start: "CCC", End: "EEE"},
			{Start: "CCC", End: "DDD"},
			{Start: "DDD", End: "BBB"},
		},
	},
	{
		Name: "TestManyLoops",
		Hops: []Hop{
			// AAA -> BBB -> CCC -> DDD -> BBB -> CCC -> DDD -> EEE -> CCC -> DDD -> EEE
			// There are two loops, DDD -> BBB and EEE -> CCC
			{Start: "EEE", End: "CCC"},
			{Start: "AAA", End: "BBB"},
			{Start: "BBB", End: "CCC"},
			{Start: "DDD", End: "EEE"},
			{Start: "BBB", End: "CCC"},
			{Start: "CCC", End: "DDD"},
			{Start: "DDD", End: "BBB"},
			{Start: "DDD", End: "EEE"},
			{Start: "CCC", End: "DDD"},
			{Start: "CCC", End: "DDD"},
		},
		Result: Hop{Start: "AAA", End: "EEE"},
	},
}

func TestReduce(t *testing.T) {
	for _, tc := range tCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			res, err := Reduce(tc.Hops)
			if tc.Fails {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.Result, res)
		})
	}
}

type tCase struct {
	Name   string
	Hops   []Hop
	Result Hop
	Fails  bool
}
