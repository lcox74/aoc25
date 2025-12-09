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

// JunctionBox represents a junction box position in 3D space.
type JunctionBox struct {
	X, Y, Z int
}

// Edge represents a pair of junction boxes and their squared distance.
type Edge struct {
	I, J   int
	DistSq int64
}

// Playground connects junction boxes with light strings to form circuits.
// Part 1: After connecting 1000 closest pairs, multiply sizes of 3 largest circuits.
// Part 2: Connect until one circuit; return product of X coords of last connection.
type Playground struct {
	boxes  []JunctionBox
	parent []int
	rank   []int

	ResultPart1 int
	ResultPart2 int
}

// NewPlayground creates a new Playground instance.
func NewPlayground() *Playground {
	return &Playground{}
}

// String implements fmt.Stringer for output.
func (p *Playground) String() string {
	return fmt.Sprintf("part1: %d, part2: %d", p.ResultPart1, p.ResultPart2)
}

// Parse reads junction box coordinates from an io.Reader.
func (p *Playground) Parse(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue
		}
		x, err1 := strconv.Atoi(parts[0])
		y, err2 := strconv.Atoi(parts[1])
		z, err3 := strconv.Atoi(parts[2])
		if err1 != nil || err2 != nil || err3 != nil {
			continue
		}
		p.boxes = append(p.boxes, JunctionBox{X: x, Y: y, Z: z})
	}
	p.solve()
}

// Solve connects junction boxes and computes circuit sizes using brute force.
// numConnections specifies how many pairs to connect for Part 1.
func (p *Playground) Solve(numConnections int) {
	n := len(p.boxes)
	if n == 0 {
		return
	}

	p.initUnionFind(n)

	edges := p.buildAllEdges()
	sortEdgesByDistance(edges)

	// Track connections and circuits
	connected := 0
	circuits := n
	var lastMerge Edge
	part1Computed := false

	for _, e := range edges {
		if p.union(e.I, e.J) {
			circuits--
			lastMerge = e
		}
		connected++

		// Part 1: compute after numConnections
		if connected == numConnections && !part1Computed {
			sizes := p.computeCircuitSizes()
			p.ResultPart1 = topNProduct(sizes, 3)
			part1Computed = true
		}

		// Part 2: stop when all in one circuit
		if circuits == 1 {
			break
		}
	}

	// Part 2: product of X coordinates of last connection
	if circuits == 1 {
		p.ResultPart2 = p.boxes[lastMerge.I].X * p.boxes[lastMerge.J].X
	}
}

// solve is called after parsing with default 1000 connections.
func (p *Playground) solve() {
	p.Solve(1000)
}

func main() {
	var inputFile string

	flag.StringVar(&inputFile, "input", "day08/input.txt", "input file path")
	flag.StringVar(&inputFile, "i", "day08/input.txt", "input file path (shorthand)")
	flag.Parse()

	if inputFile == "" {
		log.Fatal("no input file specified")
	}

	f, err := os.Open(filepath.Clean(inputFile))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	solver := NewPlayground()
	solver.Parse(f)
	fmt.Println(solver)
}
