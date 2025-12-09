package main

import (
	"bufio"
	"io"
)

// readLines reads all non-empty lines from the provided reader.
func readLines(r io.Reader) []string {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		if line := scanner.Text(); len(line) > 0 {
			lines = append(lines, line)
		}
	}
	return lines
}

// isColumnAllSpaces returns true if the specified column index contains only
// space characters (or is out of bounds) across all provided lines.
// Used to identify separator columns between problems in the worksheet.
func isColumnAllSpaces(lines []string, col int) bool {
	for _, line := range lines {
		if col < len(line) && line[col] != ' ' {
			return false
		}
	}

	return true
}

// readColumnAsNumber reads a single column of digits from top to bottom and
// interprets them as a decimal number. The topmost digit is the most significant.
// Non-digit characters are skipped.
// For example, if column 5 contains '1', '2', '3' from top to bottom, returns 123.
func readColumnAsNumber(lines []string, col int) int {
	num := 0

	for _, line := range lines {
		if col < len(line) && line[col] >= '0' && line[col] <= '9' {
			num = num*10 + int(line[col]-'0')
		}
	}

	return num
}

// applyOperation applies the given operator (+, *) to a slice of numbers.
// The first number is the initial value, and subsequent numbers are combined
// using the operator. Returns 0 for empty slices.
func applyOperation(numbers []int, op byte) int {
	if len(numbers) == 0 {
		return 0
	}

	result := numbers[0]
	for i := 1; i < len(numbers); i++ {
		switch op {
		case '+':
			result += numbers[i]
		default:
			result *= numbers[i]
		}
	}

	return result
}
