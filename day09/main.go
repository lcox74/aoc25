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

// MovieTheater finds the largest rectangle using red tiles as opposite corners.
// Part 1: Find maximum rectangle area between any two red tiles.
// Part 2: Find maximum rectangle area using only red and green tiles.
type MovieTheater struct {
	TilesX      []int // X coordinates of red tiles
	TilesY      []int // Y coordinates of red tiles
	ResultPart1 int   // Maximum rectangle area (any two red tiles)
	ResultPart2 int   // Maximum rectangle area (only red/green tiles)
}

func NewMovieTheater() *MovieTheater {
	return &MovieTheater{}
}

func (m *MovieTheater) String() string {
	return fmt.Sprintf("Movie Theater:\n\tPart 1: %d\n\tPart 2: %d", m.ResultPart1, m.ResultPart2)
}

// Parse reads coordinate pairs from r and finds the maximum rectangle area.
func (m *MovieTheater) Parse(r io.Reader) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			continue
		}

		x, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		y, err := strconv.Atoi(parts[1])
		if err != nil {
			continue
		}

		m.TilesX = append(m.TilesX, x)
		m.TilesY = append(m.TilesY, y)
	}

	m.solve()
}

// solve finds the maximum rectangle areas for both parts.
func (m *MovieTheater) solve() {
	n := len(m.TilesX)
	if n < 2 {
		return
	}

	// Build polygon interior map for Part 2
	outside, xIdx, yIdx := m.buildPolygonMap()

	// Single pass over all tile pairs
	for i := range n {
		for j := i + 1; j < n; j++ {
			x1, y1 := m.TilesX[i], m.TilesY[i]
			x2, y2 := m.TilesX[j], m.TilesY[j]

			xMin, xMax := min(x1, x2), max(x1, x2)
			yMin, yMax := min(y1, y2), max(y1, y2)
			area := (xMax - xMin + 1) * (yMax - yMin + 1)

			// Part 1: Any rectangle
			m.ResultPart1 = max(m.ResultPart1, area)

			// Part 2: Only rectangles inside polygon
			if isInsidePolygon(outside, xIdx, yIdx, xMin, xMax, yMin, yMax) {
				m.ResultPart2 = max(m.ResultPart2, area)
			}
		}
	}
}

func main() {
	var inputFile string

	flag.StringVar(&inputFile, "input", "day09/input.txt", "input file path")
	flag.StringVar(&inputFile, "i", "day09/input.txt", "input file path (shorthand)")
	flag.Parse()

	if inputFile == "" {
		log.Fatal("no input file specified")
	}

	f, err := os.Open(filepath.Clean(inputFile))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	theater := NewMovieTheater()
	theater.Parse(f)
	fmt.Println(theater)
}
