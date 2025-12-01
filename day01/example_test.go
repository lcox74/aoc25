package main_test

import (
	"strings"
	"testing"

	main "github.com/lcox74/aoc25/day01"
	"github.com/stretchr/testify/require"
)

// exampleInput is the sample input from the problem description.
const exampleInput = `L68
L30
R48
L5
R60
L55
L1
L99
R14
L82`

func TestExample(t *testing.T) {
	dial := main.NewDial()
	dial.Parse(strings.NewReader(exampleInput))

	require.Equal(t, 32, dial.Value)
	require.Equal(t, 3, dial.Strictzero)
	require.Equal(t, 6, dial.Zero)
}
