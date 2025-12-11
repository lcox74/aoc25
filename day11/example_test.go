package main_test

import (
	"strings"
	"testing"

	main "github.com/lcox74/aoc25/day11"
	"github.com/stretchr/testify/require"
)

// exampleInput is the sample input from the problem description (Part 1).
const exampleInput = `aaa: you hhh
you: bbb ccc
bbb: ddd eee
ccc: ddd eee fff
ddd: ggg
eee: out
fff: out
ggg: out
hhh: ccc fff iii
iii: out`

// exampleInputPart2 is the sample input for Part 2.
const exampleInputPart2 = `svr: aaa bbb
aaa: fft
fft: ccc
bbb: tty
tty: ccc
ccc: ddd eee
ddd: hub
hub: fff
eee: dac
dac: fff
fff: ggg hhh
ggg: out
hhh: out`

func TestExample(t *testing.T) {
	reactor := main.NewReactor()
	reactor.Parse(strings.NewReader(exampleInput))

	// Part 1: Count paths from 'you' to 'out'
	// Path 1: you -> bbb -> ddd -> ggg -> out
	// Path 2: you -> bbb -> eee -> out
	// Path 3: you -> ccc -> ddd -> ggg -> out
	// Path 4: you -> ccc -> eee -> out
	// Path 5: you -> ccc -> fff -> out
	// Total: 5 paths
	require.Equal(t, 5, reactor.ResultPart1)
}

func TestExamplePart2(t *testing.T) {
	reactor := main.NewReactor()
	reactor.Parse(strings.NewReader(exampleInputPart2))

	// Part 2: Count paths from 'svr' to 'out' that visit both 'dac' and 'fft'
	// Only 2 paths visit both dac and fft:
	// svr -> aaa -> fft -> ccc -> eee -> dac -> fff -> ggg -> out
	// svr -> aaa -> fft -> ccc -> eee -> dac -> fff -> hhh -> out
	require.Equal(t, 2, reactor.ResultPart2)
}
