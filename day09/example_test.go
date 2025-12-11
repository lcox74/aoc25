package main_test

import (
	"strings"
	"testing"

	main "github.com/lcox74/aoc25/day09"
	"github.com/stretchr/testify/require"
)

// exampleInput is the sample input from the problem description.
const exampleInput = `7,1
11,1
11,7
9,7
9,5
2,5
2,3
7,3`

func TestExample(t *testing.T) {
	theater := main.NewMovieTheater()
	theater.Parse(strings.NewReader(exampleInput))

	// Part 1: Largest rectangle area is 50 (between 2,5 and 11,1)
	// Width = |11-2|+1 = 10, Height = |5-1|+1 = 5, Area = 50
	require.Equal(t, 50, theater.ResultPart1)

	// Part 2: Largest rectangle using only red/green tiles is 24 (between 9,5 and 2,3)
	// Width = |9-2|+1 = 8, Height = |5-3|+1 = 3, Area = 24
	require.Equal(t, 24, theater.ResultPart2)
}
