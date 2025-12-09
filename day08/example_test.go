package main_test

import (
	"io"
	"os"
	"strings"
	"testing"

	main "github.com/lcox74/aoc25/day08"
	"github.com/stretchr/testify/require"
)

// exampleInput is the sample input from the problem description.
// 20 junction boxes, connecting 10 closest pairs.
const exampleInput = `162,817,812
57,618,57
906,360,560
592,479,940
352,342,300
466,668,158
542,29,236
431,825,988
739,650,466
52,470,668
216,146,977
819,987,18
117,168,530
805,96,715
346,949,466
970,615,88
941,993,340
862,61,35
984,92,344
425,690,689`

func TestExample(t *testing.T) {
	solver := main.NewPlayground()
	solver.Parse(strings.NewReader(exampleInput))

	// The example uses 10 connections instead of 1000
	solver.Solve(10)

	// Part 1: After 10 connections, the 3 largest circuits have sizes
	// 5, 4, and 2 -> 5 * 4 * 2 = 40
	require.Equal(t, 40, solver.ResultPart1)

	// Part 2: Last connection to form single circuit is between
	// 216,146,977 and 117,168,530 -> 216 * 117 = 25272
	require.Equal(t, 25272, solver.ResultPart2)
}

func TestExampleWithKDTree(t *testing.T) {
	solver := main.NewPlayground()
	solver.Parse(strings.NewReader(exampleInput))

	// Test KD-Tree solution with 10 connections
	solver.SolveWithKDTree(10)

	// Should produce same result as brute force
	require.Equal(t, 40, solver.ResultPart1)
}

func TestKDTreeMatchesBruteForce(t *testing.T) {
	// Test that KD-Tree solution matches brute force on example input
	solver1 := main.NewPlayground()
	solver1.Parse(strings.NewReader(exampleInput))
	solver1.Solve(10)

	solver2 := main.NewPlayground()
	solver2.Parse(strings.NewReader(exampleInput))
	solver2.SolveWithKDTree(10)

	require.Equal(t, solver1.ResultPart1, solver2.ResultPart1,
		"KD-Tree result should match brute force result")
}

func TestKDTreeOnRealInput(t *testing.T) {
	// This test verifies KD-Tree matches brute force on the real input
	// by reading from the input file
	f, err := os.Open("input.txt")
	if err != nil {
		t.Skip("input.txt not found, skipping real input test")
	}
	defer f.Close()

	// Read input into buffer so we can parse twice
	content, _ := io.ReadAll(f)

	solver1 := main.NewPlayground()
	solver1.Parse(strings.NewReader(string(content)))
	solver1.Solve(1000)

	solver2 := main.NewPlayground()
	solver2.Parse(strings.NewReader(string(content)))
	solver2.SolveWithKDTree(1000)

	require.Equal(t, solver1.ResultPart1, solver2.ResultPart1,
		"KD-Tree result should match brute force on real input")
}

// loadRealInput loads the real input for benchmarks.
func loadRealInput() []byte {
	content, err := os.ReadFile("input.txt")
	if err != nil {
		return nil
	}
	return content
}

func BenchmarkBruteForce(b *testing.B) {
	content := loadRealInput()
	if content == nil {
		b.Skip("input.txt not found")
	}

	for b.Loop() {
		solver := main.NewPlayground()
		solver.Parse(strings.NewReader(string(content)))
		solver.Solve(1000)
	}
}

func BenchmarkKDTree(b *testing.B) {
	content := loadRealInput()
	if content == nil {
		b.Skip("input.txt not found")
	}

	for b.Loop() {
		solver := main.NewPlayground()
		solver.Parse(strings.NewReader(string(content)))
		solver.SolveWithKDTree(1000)
	}
}
