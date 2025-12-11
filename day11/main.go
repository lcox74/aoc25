package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Reactor solves device path counting puzzles.
// Part 1: Count all paths from 'you' to 'out'.
// Part 2: Count paths from 'svr' to 'out' that visit both 'dac' and 'fft'.
type Reactor struct {
	graph       map[string][]string
	ResultPart1 int
	ResultPart2 int
}

func NewReactor() *Reactor {
	return &Reactor{
		graph: make(map[string][]string),
	}
}

func (r *Reactor) String() string {
	return fmt.Sprintf("Reactor:\n\tPart 1: %d\n\tPart 2: %d", r.ResultPart1, r.ResultPart2)
}

func (r *Reactor) Parse(rd io.Reader) {
	scanner := bufio.NewScanner(rd)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}

		device := parts[0]
		targets := strings.Fields(parts[1])
		r.graph[device] = targets
	}

	r.ResultPart1 = r.countPaths("you", make(map[string]int))
	r.ResultPart2 = r.countPathsWithCheckpoints("svr", false, false, make(map[string]int))
}

// countPaths counts all paths from current node to "out" using memoized DFS.
func (r *Reactor) countPaths(current string, memo map[string]int) int {
	if current == "out" {
		return 1
	}
	if cached, ok := memo[current]; ok {
		return cached
	}
	count := 0
	for _, next := range r.graph[current] {
		count += r.countPaths(next, memo)
	}
	memo[current] = count
	return count
}

// countPathsWithCheckpoints counts paths from current to "out" that visit both dac and fft.
func (r *Reactor) countPathsWithCheckpoints(current string, visitedDac, visitedFft bool, memo map[string]int) int {
	if current == "dac" {
		visitedDac = true
	}
	if current == "fft" {
		visitedFft = true
	}

	if current == "out" {
		if visitedDac && visitedFft {
			return 1
		}
		return 0
	}

	key := current
	if visitedDac {
		key += ":d"
	}
	if visitedFft {
		key += ":f"
	}
	if cached, ok := memo[key]; ok {
		return cached
	}

	count := 0
	for _, next := range r.graph[current] {
		count += r.countPathsWithCheckpoints(next, visitedDac, visitedFft, memo)
	}
	memo[key] = count
	return count
}

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "input", "day11/input.txt", "input file path")
	flag.StringVar(&inputFile, "i", "day11/input.txt", "input file path (shorthand)")
	flag.Parse()

	if inputFile == "" {
		log.Fatal("no input file specified")
	}

	f, err := os.Open(filepath.Clean(inputFile))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reactor := NewReactor()
	reactor.Parse(f)
	fmt.Println(reactor)
}
