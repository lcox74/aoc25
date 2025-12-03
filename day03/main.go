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

// BatteryBank finds the maximum joltage from each bank of batteries.
// Part 1: Select exactly 2 batteries to form a two-digit number.
// Part 2: Select exactly 12 batteries to form a twelve-digit number.
type BatteryBank struct {
	TotalJoltage2Bat  int64
	TotalJoltage12Bat int64
}

func NewBatteryBank() *BatteryBank {
	return &BatteryBank{}
}

// Parse reads battery banks from r, one per line.
// For each bank, finds the maximum joltage by selecting batteries.
func (b *BatteryBank) Parse(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 2 {
			continue
		}

		b.TotalJoltage2Bat += b.findMaxJoltageN(line, 2)   // Part 1: select 2 batteries
		b.TotalJoltage12Bat += b.findMaxJoltageN(line, 12) // Part 2: select 12 batteries
	}
}

func (b *BatteryBank) String() string {
	return fmt.Sprintf(
		"Total Joltage: \n\tSmall Bat: %d jolts\n\tBig Bat: %d jolts",
		b.TotalJoltage2Bat,
		b.TotalJoltage12Bat,
	)
}

// findMaxJoltageN finds the maximum number by selecting exactly n digits.
// I'm doing it the lazy way in brute-force fashion, since I am lazy and
// need to get back to work.
func (b *BatteryBank) findMaxJoltageN(bank string, n int) int64 {
	if len(bank) < n {
		return 0
	}

	currentMax := int64(0)
	pos := 0 // current position in bank

	for i := range n {
		maxPos := len(bank) - (n - i - 1) - 1 // furthest we can look ahead

		// Find the largest digit from pos to maxPos
		bestDigit := byte('0')
		bestIdx := pos
		for j := pos; j <= maxPos; j++ {
			if bank[j] > bestDigit {
				bestDigit = bank[j]
				bestIdx = j
			}
		}

		// Append bestDigit to currentMax
		currentMax = currentMax*10 + int64(bestDigit-'0')
		pos = bestIdx + 1
	}

	return currentMax
}

func main() {
	var inputFile string

	flag.StringVar(&inputFile, "input", "day03/input.txt", "input file path")
	flag.StringVar(&inputFile, "i", "day03/input.txt", "input file path (shorthand)")
	flag.Parse()

	if inputFile == "" {
		log.Fatal("no input file specified")
	}

	f, err := os.Open(filepath.Clean(inputFile))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	bank := NewBatteryBank()
	bank.Parse(f)
	fmt.Println(bank)
}
