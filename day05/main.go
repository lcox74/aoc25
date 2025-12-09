package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
)

// Cafeteria checks which ingredient IDs are fresh.
// Part 1: Count how many available ingredient IDs fall within any fresh range.
// Part 2: Count total unique fresh IDs across all ranges.
type Cafeteria struct {
	Ranges      [][2]int
	Ingredients []int
	FreshCount  int // Part 1: fresh available ingredients
	TotalFresh  int // Part 2: total unique IDs in all ranges
}

func NewCafeteria() *Cafeteria {
	return &Cafeteria{}
}

func (c *Cafeteria) String() string {
	return fmt.Sprintf("Fresh Ingredients:\n\tPart 1: %d\n\tPart 2: %d", c.FreshCount, c.TotalFresh)
}

// Parse reads the database from r.
// First section contains fresh ID ranges (e.g., "3-5").
// After a blank line, the second section contains available ingredient IDs.
func (c *Cafeteria) Parse(r io.Reader) {
	scanner := bufio.NewScanner(r)
	parsingRanges := true

	for scanner.Scan() {
		line := scanner.Text()

		// Blank line switches from ranges to ingredients
		if line == "" {
			parsingRanges = false
			continue
		}

		if parsingRanges {
			parts := strings.Split(line, "-")
			if len(parts) != 2 {
				continue
			}

			start, err1 := strconv.Atoi(parts[0])
			end, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				continue
			}

			c.Ranges = append(c.Ranges, [2]int{start, end})
		} else {
			id, err := strconv.Atoi(line)
			if err != nil {
				continue
			}

			c.Ingredients = append(c.Ingredients, id)
		}
	}

	// Part 1: Count fresh available ingredients
	for _, id := range c.Ingredients {
		if c.isFresh(id) {
			c.FreshCount++
		}
	}

	// Part 2: Count total unique IDs across all ranges
	c.TotalFresh = c.countTotalFreshIDs()
}

// isFresh returns true if the ID falls within any fresh range.
func (c *Cafeteria) isFresh(id int) bool {
	for _, r := range c.Ranges {
		if id >= r[0] && id <= r[1] {
			return true
		}
	}
	return false
}

// countTotalFreshIDs merges overlapping ranges and counts total unique IDs.
func (c *Cafeteria) countTotalFreshIDs() int {
	if len(c.Ranges) == 0 {
		return 0
	}

	// Sort ranges by start value
	sorted := make([][2]int, len(c.Ranges))
	copy(sorted, c.Ranges)
	slices.SortFunc(sorted, func(a, b [2]int) int {
		return a[0] - b[0]
	})

	// Merge overlapping ranges
	merged := [][2]int{sorted[0]}
	for _, r := range sorted[1:] {
		last := &merged[len(merged)-1]
		if r[0] <= last[1]+1 {
			// Overlapping or adjacent, extend the range
			if r[1] > last[1] {
				last[1] = r[1]
			}
		} else {
			// Non-overlapping, add new range
			merged = append(merged, r)
		}
	}

	// Count total IDs in merged ranges
	total := 0
	for _, r := range merged {
		total += r[1] - r[0] + 1
	}
	return total
}

func main() {
	var inputFile string

	flag.StringVar(&inputFile, "input", "day05/input.txt", "input file path")
	flag.StringVar(&inputFile, "i", "day05/input.txt", "input file path (shorthand)")
	flag.Parse()

	if inputFile == "" {
		log.Fatal("no input file specified")
	}

	f, err := os.Open(filepath.Clean(inputFile))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	cafe := NewCafeteria()
	cafe.Parse(f)
	fmt.Println(cafe)
}
