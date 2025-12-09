package main_test

import (
	"strings"
	"testing"

	main "github.com/lcox74/aoc25/day05"
	"github.com/stretchr/testify/require"
)

// exampleInput is the sample input from the problem description.
const exampleInput = `3-5
10-14
16-20
12-18

1
5
8
11
17
32`

func TestExample(t *testing.T) {
	cafe := main.NewCafeteria()
	cafe.Parse(strings.NewReader(exampleInput))

	// Part 1: 3 fresh ingredients (5, 11, 17)
	require.Equal(t, 3, cafe.FreshCount)

	// Part 2: 14 unique fresh IDs (3-5, 10-20 merged)
	require.Equal(t, 14, cafe.TotalFresh)
}
