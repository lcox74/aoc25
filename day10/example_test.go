package main_test

import (
	"strings"
	"testing"

	main "github.com/lcox74/aoc25/day10"
	"github.com/stretchr/testify/require"
)

// exampleInput is the sample input from the problem description.
const exampleInput = `[.##.] (3) (1,3) (2) (2,3) (0,2) (0,1) {3,5,4,7}
[...#.] (0,2,3,4) (2,3) (0,4) (0,1,2) (1,2,3,4) {7,5,12,7,2}
[.###.#] (0,1,2,3,4) (0,3,4) (0,1,2,4,5) (1,2) {10,11,11,5,10,5}`

func TestExample(t *testing.T) {
	factory := main.NewFactory()
	factory.Parse(strings.NewReader(exampleInput))

	// Part 1: Minimum button presses for all machines (XOR/toggle)
	// Machine 1: [.##.] -> 2 presses (buttons (0,2) and (0,1))
	// Machine 2: [...#.] -> 3 presses (buttons (0,4), (0,1,2), (1,2,3,4))
	// Machine 3: [.###.#] -> 2 presses (buttons (0,3,4) and (0,1,2,4,5))
	// Total: 2 + 3 + 2 = 7
	require.Equal(t, 7, factory.ResultPart1)

	// Part 2: Minimum button presses for joltage counters (addition)
	// Machine 1: {3,5,4,7} -> 10 presses
	// Machine 2: {7,5,12,7,2} -> 12 presses
	// Machine 3: {10,11,11,5,10,5} -> 11 presses
	// Total: 10 + 12 + 11 = 33
	require.Equal(t, 33, factory.ResultPart2)
}
