package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// TachyonManifold simulates tachyon beams passing through a manifold with splitters.
// Part 1: Count how many times beams are split by splitters (^).
// Part 2: Count distinct timelines (paths) through the manifold.
type TachyonManifold struct {
	splitters [][]uint64 // bitmask of splitter positions per row (only non-empty rows)
	numWords  int
	width     int
	startCol  int

	ResultPart1 int
	ResultPart2 int
}

// NewTachyonManifold creates a new TachyonManifold instance.
func NewTachyonManifold() *TachyonManifold {
	return &TachyonManifold{}
}

// String implements fmt.Stringer for output.
func (t *TachyonManifold) String() string {
	return fmt.Sprintf("part1: %d, part2: %d", t.ResultPart1, t.ResultPart2)
}

// Parse reads the manifold diagram from an io.Reader.
func (t *TachyonManifold) Parse(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if t.numWords == 0 {
			t.width = len(line)
			t.numWords = (len(line) + 63) / 64
		}

		rowMask := make([]uint64, t.numWords)
		hasSplitters := false
		for i := range len(line) {
			switch line[i] {
			case 'S':
				t.startCol = i
			case '^':
				rowMask[i/64] |= 1 << (i % 64)
				hasSplitters = true
			}
		}
		// Only store rows that have splitters
		if hasSplitters {
			t.splitters = append(t.splitters, rowMask)
		}
	}
	t.solve()
}

// solve simulates beams through the manifold, computing both parts in one pass.
// Part 1: count splitter hits. Part 2: count distinct timelines.
func (t *TachyonManifold) solve() {
	if len(t.splitters) == 0 {
		t.ResultPart2 = 1
		return
	}

	timelines := make([]int, t.width)
	next := make([]int, t.width)
	timelines[t.startCol] = 1

	splitCount := 0
	for _, splitterMask := range t.splitters {
		clear(next)
		for col, count := range timelines {
			if count == 0 {
				continue
			}
			if t.hasSplitter(splitterMask, col) {
				splitCount++
				if col > 0 {
					next[col-1] += count
				}
				if col+1 < t.width {
					next[col+1] += count
				}
			} else {
				next[col] += count
			}
		}
		timelines, next = next, timelines
	}

	t.ResultPart1 = splitCount
	t.ResultPart2 = sum(timelines)
}

// hasSplitter checks if there's a splitter at the given column in the mask.
func (t *TachyonManifold) hasSplitter(mask []uint64, col int) bool {
	return (mask[col/64] & (1 << (col % 64))) != 0
}

// sum returns the sum of all values in the slice.
func sum(vals []int) int {
	total := 0
	for _, v := range vals {
		total += v
	}
	return total
}

func main() {
	var inputFile string

	flag.StringVar(&inputFile, "input", "day07/input.txt", "input file path")
	flag.StringVar(&inputFile, "i", "day07/input.txt", "input file path (shorthand)")
	flag.Parse()

	if inputFile == "" {
		log.Fatal("no input file specified")
	}

	f, err := os.Open(filepath.Clean(inputFile))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	solver := NewTachyonManifold()
	solver.Parse(f)
	fmt.Println(solver)
}
