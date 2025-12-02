package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// GiftShop checks product ID ranges for invalid IDs.
// Part 1: Invalid IDs are numbers made of a digit sequence repeated exactly twice (e.g., 55, 6464, 123123).
// Part 2: Invalid IDs are numbers made of a digit sequence repeated at least twice.
type GiftShop struct {
	Ranges      [][2]int64
	InvalidSum1 int64 // Part 1: exactly twice
	InvalidSum2 int64 // Part 2: at least twice
	Verbose     bool
}

func NewGiftShop() *GiftShop {
	return &GiftShop{}
}

// Parse reads product ID ranges from r.
// Each range is formatted as "start-end" and separated by commas.
func (g *GiftShop) Parse(r io.Reader) {
	scanner := bufio.NewScanner(r)
	var input strings.Builder
	for scanner.Scan() {
		input.WriteString(scanner.Text())
	}

	// Parse the bounds
	parts := strings.SplitSeq(strings.TrimSpace(input.String()), ",")
	for part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		bounds := strings.Split(part, "-")
		if len(bounds) != 2 {
			continue
		}

		start, err1 := strconv.ParseInt(bounds[0], 10, 64)
		end, err2 := strconv.ParseInt(bounds[1], 10, 64)
		if err1 != nil || err2 != nil {
			continue
		}

		g.Ranges = append(g.Ranges, [2]int64{start, end})
	}

	// Consolidate invalid IDs for all ranges
	g.InvalidSum1 = 0
	g.InvalidSum2 = 0
	for _, r := range g.Ranges {
		ids1 := g.findInvalidIDsInRange(r[0], r[1], false)
		ids2 := g.findInvalidIDsInRange(r[0], r[1], true)

		g.InvalidSum1 += sum(ids1)
		g.InvalidSum2 += sum(ids2)

		if g.Verbose {
			fmt.Printf("  %d-%d: part1=%v part2=%v\n", r[0], r[1], ids1, ids2)
		}
	}
}

func (g *GiftShop) String() string {
	return fmt.Sprintf("ranges: %d, part1: %d, part2: %d", len(g.Ranges), g.InvalidSum1, g.InvalidSum2)
}

// findInvalidIDsInRange returns all invalid IDs within the given range.
func (g *GiftShop) findInvalidIDsInRange(start, end int64, atLeastTwice bool) []int64 {
	startLen := digitLength(start)
	endLen := digitLength(end)

	if atLeastTwice {
		return g.findAtLeastTwiceInRange(start, end, startLen, endLen)
	}
	return g.findExactlyTwiceInRange(start, end, startLen, endLen)
}

// findAtLeastTwiceInRange finds IDs with a pattern repeated 2 or more times.
func (g *GiftShop) findAtLeastTwiceInRange(start, end int64, startLen, endLen int) []int64 {
	var ids []int64
	seen := make(map[int64]bool)

	for numDigits := startLen; numDigits <= endLen; numDigits++ {
		g.collectAtLeastTwiceIDs(numDigits, start, end, seen, &ids)
	}
	return ids
}

// collectAtLeastTwiceIDs collects IDs of a specific digit length with repeated patterns.
func (g *GiftShop) collectAtLeastTwiceIDs(numDigits int, start, end int64, seen map[int64]bool, ids *[]int64) {
	for patternLen := 1; patternLen <= numDigits/2; patternLen++ {
		if numDigits%patternLen != 0 {
			continue
		}

		// Calculate the number of repetitions for the pattern
		repetitions := numDigits / patternLen
		if repetitions < 2 {
			continue
		}

		// Generate all patterns of the given length
		minPattern, maxPattern := patternBounds(patternLen)
		for pattern := minPattern; pattern <= maxPattern; pattern++ {
			invalidID := buildRepeatedID(pattern, patternLen, repetitions)
			if invalidID >= start && invalidID <= end && !seen[invalidID] {
				*ids = append(*ids, invalidID)
				seen[invalidID] = true
			}
		}
	}
}

// findExactlyTwiceInRange finds IDs with a pattern repeated exactly twice.
func (g *GiftShop) findExactlyTwiceInRange(start, end int64, startLen, endLen int) []int64 {
	var ids []int64

	for numDigits := startLen; numDigits <= endLen; numDigits++ {
		if numDigits%2 != 0 {
			continue
		}

		halfLen := numDigits / 2
		minHalf, maxHalf := patternBounds(halfLen)

		// Generate all halves and form the invalid ID by repeating the half
		for half := minHalf; half <= maxHalf; half++ {
			invalidID := buildRepeatedID(half, halfLen, 2)
			if invalidID >= start && invalidID <= end {
				ids = append(ids, invalidID)
			}
		}
	}
	return ids
}

func main() {
	var inputFile string
	var verbose bool

	flag.StringVar(&inputFile, "input", "day02/input.txt", "input file path")
	flag.StringVar(&inputFile, "i", "day02/input.txt", "input file path (shorthand)")
	flag.BoolVar(&verbose, "verbose", false, "print invalid IDs found")
	flag.BoolVar(&verbose, "v", false, "print invalid IDs found (shorthand)")
	flag.Parse()

	if inputFile == "" {
		log.Fatal("no input file specified")
	}

	f, err := os.Open(filepath.Clean(inputFile))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	shop := NewGiftShop()
	shop.Verbose = verbose
	shop.Parse(f)
	fmt.Println(shop)
}
