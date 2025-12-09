package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// MathWorksheet solves cephalopod math homework problems.
// Part 1: Sum of all problem solutions reading numbers horizontally.
// Part 2: Sum of all problem solutions reading numbers vertically (columns) right-to-left.
type MathWorksheet struct {
	ResultPart1 int
	ResultPart2 int
}

// NewMathWorksheet creates a new MathWorksheet instance.
func NewMathWorksheet() *MathWorksheet {
	return &MathWorksheet{}
}

// String implements fmt.Stringer for output.
func (m *MathWorksheet) String() string {
	return fmt.Sprintf("part1: %d, part2: %d", m.ResultPart1, m.ResultPart2)
}

// Parse reads the worksheet and solves all problems.
func (m *MathWorksheet) Parse(r io.Reader) {
	lines := readLines(r)
	if len(lines) < 2 {
		return
	}

	// Part 1: horizontal reading
	m.ResultPart1 = solveHorizontal(lines)

	// Part 2: vertical reading (columns as numbers, right-to-left)
	m.ResultPart2 = solveVertical(lines)
}

// solveHorizontal reads numbers horizontally (left-to-right on each row)
// and applies operations column by column across all rows.
func solveHorizontal(lines []string) int {
	opLine := lines[len(lines)-1]
	numLines := lines[:len(lines)-1]

	// Parse each row into numbers
	var numRows [][]int
	for _, line := range numLines {
		var row []int
		for field := range strings.FieldsSeq(line) {
			if num, err := strconv.Atoi(field); err == nil {
				row = append(row, num)
			}
		}
		numRows = append(numRows, row)
	}

	// Parse operators
	ops := strings.Fields(opLine)

	// Solve each problem (column) by combining numbers vertically with the operator
	total := 0
	for col, op := range ops {
		var numbers []int
		for _, row := range numRows {
			if col < len(row) {
				numbers = append(numbers, row[col])
			}
		}
		total += applyOperation(numbers, op[0])
	}
	return total
}

func solveVertical(lines []string) int {
	opRow := lines[len(lines)-1]
	numLines := lines[:len(lines)-1]

	// Find max width
	maxWidth := 0
	for _, line := range lines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	// Find problem boundaries and operators by scanning right-to-left
	total := 0
	col := maxWidth - 1

	for col >= 0 {
		// Skip separator columns (all spaces in number rows)
		if isColumnAllSpaces(numLines, col) {
			col--
			continue
		}

		// Collect all digit columns for this problem (until we hit a separator)
		// Also find the operator within this problem's columns
		var numbers []int
		op := byte('*') // default
		for col >= 0 && !isColumnAllSpaces(numLines, col) {
			// Check for operator in this column
			if col < len(opRow) && (opRow[col] == '+' || opRow[col] == '*') {
				op = opRow[col]
			}
			num := readColumnAsNumber(numLines, col)
			numbers = append(numbers, num)
			col--
		}

		// Apply operation to numbers
		total += applyOperation(numbers, op)
	}

	return total
}

func main() {
	var inputFile string

	flag.StringVar(&inputFile, "input", "day06/input.txt", "input file path")
	flag.StringVar(&inputFile, "i", "day06/input.txt", "input file path (shorthand)")
	flag.Parse()

	if inputFile == "" {
		log.Fatal("no input file specified")
	}

	f, err := os.Open(filepath.Clean(inputFile))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	solver := NewMathWorksheet()
	solver.Parse(f)
	fmt.Println(solver)
}
