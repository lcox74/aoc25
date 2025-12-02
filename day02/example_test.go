package main

import (
	"strings"
	"testing"
)

// exampleInput is the sample input from the problem description.
const exampleInput = `11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124`

func TestExample(t *testing.T) {
	expectedPart1 := int64(1227775554)
	expectedPart2 := int64(4174379265)

	shop := NewGiftShop()
	shop.Parse(strings.NewReader(exampleInput))

	if shop.InvalidSum1 != expectedPart1 {
		t.Errorf("Part 1: expected %d, got %d", expectedPart1, shop.InvalidSum1)
	}
	if shop.InvalidSum2 != expectedPart2 {
		t.Errorf("Part 2: expected %d, got %d", expectedPart2, shop.InvalidSum2)
	}
}
