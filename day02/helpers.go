package main

import "strconv"

// pow10 returns 10^n as an int64.
func pow10(n int) int64 {
	result := int64(1)
	for range n {
		result *= 10
	}
	return result
}

// patternBounds returns the min and max values for a number with the given digit length.
// For length 1: returns (1, 9)
// For length 2: returns (10, 99)
// etc.
func patternBounds(length int) (minVal, maxVal int64) {
	if length == 1 {
		return 1, 9
	}

	return pow10(length - 1), pow10(length) - 1
}

// buildRepeatedID constructs a number by repeating a pattern a given number of times.
// For example, buildRepeatedID(12, 2, 3) returns 121212.
func buildRepeatedID(pattern int64, patternLen, repetitions int) int64 {
	var result, multiplier int64 = 0, 1

	// Repeat the pattern the specified number of times
	for range repetitions {
		result += pattern * multiplier
		multiplier *= pow10(patternLen)
	}

	return result
}

// digitLength returns the number of digits in n.
func digitLength(n int64) int {
	return len(strconv.FormatInt(n, 10))
}

// sum returns the sum of a slice of numbers.
func sum[T ~int | ~int64 | ~float64](nums []T) T {
	var s T
	for _, n := range nums {
		s += n
	}
	return s
}
