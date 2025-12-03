package main_test

import (
	"strings"
	"testing"

	main "github.com/lcox74/aoc25/day03"
	"github.com/stretchr/testify/require"
)

// exampleInput is the sample input from the problem description.
const exampleInput = `987654321111111
811111111111119
234234234234278
818181911112111`

func TestExample(t *testing.T) {
	bank := main.NewBatteryBank()
	bank.Parse(strings.NewReader(exampleInput))

	// Part 1: 98 + 89 + 78 + 92 = 357
	require.Equal(t, int64(357), bank.TotalJoltage2Bat)

	// Part 2: 987654321111 + 811111111119 + 434234234278 + 888911112111 = 3121910778619
	require.Equal(t, int64(3121910778619), bank.TotalJoltage12Bat)
}
