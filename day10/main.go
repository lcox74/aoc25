package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Factory struct {
	ResultPart1 int
	ResultPart2 int
}

func NewFactory() *Factory { return &Factory{} }

func (f *Factory) String() string {
	return fmt.Sprintf("Factory:\n\tPart 1: %d\n\tPart 2: %d", f.ResultPart1, f.ResultPart2)
}

func (f *Factory) Parse(r io.Reader) {
	scanner := bufio.NewScanner(r)
	patternRe := regexp.MustCompile(`\[([.#]+)\]`)
	buttonRe := regexp.MustCompile(`\(([0-9,]*)\)`)
	joltageRe := regexp.MustCompile(`\{([0-9,]+)\}`)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		pm := patternRe.FindStringSubmatch(line)
		if pm == nil {
			continue
		}

		buttons := parseButtons(buttonRe.FindAllStringSubmatch(line, -1))
		f.ResultPart1 += solveXOR(pm[1], buttons)
		if jm := joltageRe.FindStringSubmatch(line); jm != nil {
			f.ResultPart2 += solveAdd(parseInts(jm[1]), buttons)
		}
	}
}

func parseButtons(matches [][]string) [][]int {
	var buttons [][]int

	for _, m := range matches {
		buttons = append(buttons, parseInts(m[1]))
	}

	return buttons
}

func parseInts(s string) []int {
	if s == "" {
		return nil
	}

	var result []int
	for p := range strings.SplitSeq(s, ",") {
		if v, err := strconv.Atoi(p); err == nil {
			result = append(result, v)
		}
	}

	return result
}

func solveXOR(pattern string, buttons [][]int) int {
	target := patternToMask(pattern)
	if target == 0 {
		return 0
	}

	btnMasks := buttonsToMasks(buttons, len(pattern))

	type state struct {
		val     int
		presses int
	}
	visited := map[int]bool{0: true}
	queue := []state{{0, 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, mask := range btnMasks {

			next := cur.val ^ mask
			if next == target {
				return cur.presses + 1
			}

			if !visited[next] {
				visited[next] = true
				queue = append(queue, state{next, cur.presses + 1})
			}
		}
	}
	return 0
}

func solveAdd(joltages []int, buttons [][]int) int {
	m, numBtn := len(joltages), len(buttons)
	if m == 0 || numBtn == 0 {
		return 0
	}

	mat := make([][]float64, m)
	for i := range m {
		mat[i] = make([]float64, numBtn+1)
		mat[i][numBtn] = float64(joltages[i])
	}

	for j, btn := range buttons {
		for _, idx := range btn {
			if idx < m {
				mat[idx][j] = 1
			}
		}
	}

	pivots := gaussElim(mat, m, numBtn)

	for r := len(pivots); r < m; r++ {
		if abs(mat[r][numBtn]) > 1e-9 {
			return 0
		}
	}

	freeVars := findFree(pivots, numBtn)
	if len(freeVars) == 0 {
		return sumSolution(mat, pivots, numBtn)
	}

	return searchMin(mat, pivots, freeVars, numBtn)
}

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "input", "day10/input.txt", "input file path")
	flag.StringVar(&inputFile, "i", "day10/input.txt", "input file path (shorthand)")
	flag.Parse()

	if inputFile == "" {
		log.Fatal("no input file specified")
	}

	f, err := os.Open(filepath.Clean(inputFile))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	factory := NewFactory()
	factory.Parse(f)
	fmt.Println(factory)
}
