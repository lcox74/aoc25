package main_test

import (
	"strings"
	"testing"

	main "github.com/lcox74/aoc25/day06"
	"github.com/stretchr/testify/require"
)

// exampleInput is the sample input from the problem description.
const exampleInput = `123 328  51 64
 45 64  387 23
  6 98  215 314
*   +   *   +  `

func TestExample(t *testing.T) {
	solver := main.NewMathWorksheet()
	solver.Parse(strings.NewReader(exampleInput))

	// Part 1: 123*45*6=33210, 328+64+98=490, 51*387*215=4243455, 64+23+314=401
	// Grand total: 33210 + 490 + 4243455 + 401 = 4277556
	require.Equal(t, 4277556, solver.ResultPart1)

	// Part 2: Reading columns vertically, right-to-left
	// 4+431+623=1058, 175*581*32=3253600, 8+248+369=625, 356*24*1=8544
	// Grand total: 1058 + 3253600 + 625 + 8544 = 3263827
	require.Equal(t, 3263827, solver.ResultPart2)
}
