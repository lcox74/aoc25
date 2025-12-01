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
)

// Dial tracks a rotating dial that wraps at 0-99.
// It counts how many times the dial passes through zero.
type Dial struct {
	Value      int
	Strictzero int // times landed exactly on zero
	Zero       int // times passed through zero
}

func NewDial() *Dial {
	return &Dial{Value: 50}
}

// Parse reads rotation instructions from r.
// Each line is a direction (L/R) followed by a number, e.g. "L68" or "R30".
// L rotates left (counter-clockwise), R rotates right (clockwise).
func (d *Dial) Parse(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 2 {
			continue
		}

		dir := line[0]
		n, err := strconv.Atoi(line[1:])
		if err != nil {
			continue
		}

		switch dir {
		case 'R':
			d.rotate(n)
		case 'L':
			d.rotate(-n)
		}
	}
}

func (d *Dial) String() string {
	return fmt.Sprintf("value: %d, part1: %d, part2: %d", d.Value, d.Strictzero, d.Zero)
}

func (d *Dial) rotate(n int) {
	// Count zero crossings based on direction
	if n >= 0 {
		d.Zero += (d.Value + n) / 100
	} else {
		d.Zero += ((100-d.Value)%100 - n) / 100
	}

	// Update value with wrap-around
	d.Value = ((d.Value+n)%100 + 100) % 100
	if d.Value == 0 {
		d.Strictzero++
	}
}

func main() {
	var inputFile string

	flag.StringVar(&inputFile, "input", "day01/input.txt", "input file path")
	flag.StringVar(&inputFile, "i", "day01/input.txt", "input file path (shorthand)")
	flag.Parse()

	// Validate input file
	if inputFile == "" {
		log.Fatal("no input file specified")
	}

	// Open the input file
	f, err := os.Open(filepath.Clean(inputFile))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Process the dial instructions
	dial := NewDial()
	dial.Parse(f)
	fmt.Println(dial)
}
