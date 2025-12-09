package main_test

import (
	"strings"
	"testing"

	main "github.com/lcox74/aoc25/day04"
	"github.com/stretchr/testify/require"
)

// exampleInput is the sample input from the problem description.
const exampleInput = `..@@.@@@@.
@@@.@.@.@@
@@@@@.@.@@
@.@@@@..@.
@@.@@@@.@@
.@@@@@@@.@
.@.@.@.@@@
@.@@@.@@@@
.@@@@@@@@.
@.@.@@@.@.`

func TestExample(t *testing.T) {
	dept := main.NewPrintDept()
	dept.Parse(strings.NewReader(exampleInput))

	// Part 1: 13 rolls accessible (fewer than 4 adjacent rolls)
	require.Equal(t, 13, dept.AccessibleRolls)

	// Part 2: 43 total rolls removed after iterative removal
	require.Equal(t, 43, dept.TotalRemoved)
}
