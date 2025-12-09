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

// PrintDept finds accessible paper rolls in the printing department.
// A roll is accessible if fewer than 4 rolls are in adjacent positions.
type PrintDept struct {
	Grid            []byte
	Width           int
	Height          int
	Kernel          []int
	KernelSize      int
	AccessibleRolls int // Part 1: initial accessible count
	TotalRemoved    int // Part 2: total removed after iterative removal
}

func NewPrintDept() *PrintDept {
	return &PrintDept{
		// Default 3x3 kernel: check all 8 neighbors, skip center
		Kernel:     []int{1, 1, 1, 1, 0, 1, 1, 1, 1},
		KernelSize: 3,
	}
}

// Parse reads the grid from r and counts accessible rolls.
func (p *PrintDept) Parse(r io.Reader) {
	scanner := bufio.NewScanner(r)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		return
	}

	p.Width = len(lines[0])
	p.Height = len(lines)
	p.Grid = make([]byte, p.Width*p.Height)

	for y, line := range lines {
		for x := range p.Width {
			if x < len(line) {
				p.Grid[p.Width*y+x] = line[x]
			}
		}
	}

	p.AccessibleRolls = p.countAccessibleRolls()
	p.removeAllAccessible()
}

func (p *PrintDept) String() string {
	return fmt.Sprintf(
		"Paper Rolls:\n\tAccessible: %d\n\tTotal Removed: %d",
		p.AccessibleRolls,
		p.TotalRemoved,
	)
}

// countAccessibleRolls counts rolls with fewer than 4 adjacent rolls.
func (p *PrintDept) countAccessibleRolls() int {
	count := 0
	for y := range p.Height {
		for x := range p.Width {
			idx := p.Width*y + x
			if p.Grid[idx] != '@' {
				continue
			}

			neighborCount := p.getNeighborCount(x, y)
			if neighborCount < 4 {
				count++
			}
		}
	}
	return count
}

// removeAllAccessible iteratively removes accessible rolls until none remain.
func (p *PrintDept) removeAllAccessible() {
	for {
		toRemove := p.findAccessiblePositions()
		if len(toRemove) == 0 {
			break
		}

		for _, idx := range toRemove {
			p.Grid[idx] = '.'
		}
		p.TotalRemoved += len(toRemove)
	}
}

// findAccessiblePositions returns indices of all currently accessible rolls.
func (p *PrintDept) findAccessiblePositions() []int {
	var positions []int
	for y := range p.Height {
		for x := range p.Width {
			idx := p.Width*y + x
			if p.Grid[idx] != '@' {
				continue
			}

			neighborCount := p.getNeighborCount(x, y)
			if neighborCount < 4 {
				positions = append(positions, idx)
			}
		}
	}
	return positions
}

// getNeighborCount applies the kernel to count adjacent rolls.
func (p *PrintDept) getNeighborCount(x, y int) int {
	count := 0
	halfK := p.KernelSize / 2

	for ky := range p.KernelSize {
		for kx := range p.KernelSize {
			if p.Kernel[ky*p.KernelSize+kx] == 0 {
				continue
			}

			nx := x + kx - halfK
			ny := y + ky - halfK

			if nx < 0 || nx >= p.Width || ny < 0 || ny >= p.Height {
				continue
			}

			if p.Grid[p.Width*ny+nx] == '@' {
				count++
			}
		}
	}

	return count
}

func main() {
	var inputFile string

	flag.StringVar(&inputFile, "input", "day04/input.txt", "input file path")
	flag.StringVar(&inputFile, "i", "day04/input.txt", "input file path (shorthand)")
	flag.Parse()

	if inputFile == "" {
		log.Fatal("no input file specified")
	}

	f, err := os.Open(filepath.Clean(inputFile))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	dept := NewPrintDept()
	dept.Parse(f)
	fmt.Println(dept)
}
