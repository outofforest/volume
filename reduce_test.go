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
		Name: "TestSingleHopJurney",
		Hops: []Hop{
			{Start: "AAA", End: "BBB"},
		},
		Result: Hop{"AAA", "BBB"},
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
