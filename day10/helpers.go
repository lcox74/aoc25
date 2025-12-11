package main

func patternToMask(pattern string) int {
	mask := 0
	for i, c := range pattern {
		if c == '#' {
			mask |= 1 << i
		}
	}
	return mask
}

func buttonsToMasks(buttons [][]int, n int) []int {
	masks := make([]int, len(buttons))
	for i, btn := range buttons {
		for _, idx := range btn {
			if idx < n {
				masks[i] |= 1 << idx
			}
		}
	}
	return masks
}

func findPivotRow(mat [][]float64, start, m, col int) int {
	for r := start; r < m; r++ {
		if abs(mat[r][col]) > 1e-9 {
			return r
		}
	}
	return -1
}

func eliminateCol(mat [][]float64, pivotRow, m, col int) {
	for r := range m {
		if r != pivotRow && abs(mat[r][col]) > 1e-9 {
			f := mat[r][col]
			for j := range mat[r] {
				mat[r][j] -= f * mat[pivotRow][j]
			}
		}
	}
}

func gaussElim(mat [][]float64, m, n int) []int {
	var pivots []int
	row := 0
	for col := 0; col < n && row < m; col++ {
		pivot := findPivotRow(mat, row, m, col)
		if pivot < 0 {
			continue
		}
		mat[row], mat[pivot] = mat[pivot], mat[row]
		scale := mat[row][col]
		for j := range mat[row] {
			mat[row][j] /= scale
		}
		eliminateCol(mat, row, m, col)
		pivots = append(pivots, col)
		row++
	}
	return pivots
}

func findFree(pivots []int, n int) []int {
	pset := make(map[int]bool)
	for _, c := range pivots {
		pset[c] = true
	}
	var free []int
	for c := range n {
		if !pset[c] {
			free = append(free, c)
		}
	}
	return free
}

func sumSolution(mat [][]float64, pivots []int, n int) int {
	total := 0
	for i := range pivots {
		v := int(mat[i][n] + 0.5)
		if v < 0 {
			return 0
		}
		total += v
	}
	return total
}

func extractCoefs(mat [][]float64, pivots, freeVars []int, n int) ([][]float64, []float64) {
	coefs := make([][]float64, len(pivots))
	targets := make([]float64, len(pivots))
	for i := range pivots {
		coefs[i] = make([]float64, len(freeVars))
		targets[i] = mat[i][n]
		for j, fc := range freeVars {
			coefs[i][j] = mat[i][fc]
		}
	}
	return coefs, targets
}

func evalFreeVars(coefs [][]float64, targets []float64, freeVals []int) int {
	total := 0
	for _, v := range freeVals {
		total += v
	}
	for i := range targets {
		val := targets[i]
		for j, fv := range freeVals {
			val -= coefs[i][j] * float64(fv)
		}
		rounded := int(val + 0.5)
		if rounded < 0 || abs(val-float64(rounded)) > 1e-6 {
			return -1
		}
		total += rounded
	}
	return total
}

func searchMin(mat [][]float64, pivots, freeVars []int, n int) int {
	coefs, targets := extractCoefs(mat, pivots, freeVars, n)

	maxFree := 0
	for _, t := range targets {
		if int(t) > maxFree {
			maxFree = int(t)
		}
	}
	maxFree += 50

	minTotal := -1
	var search func(idx int, freeVals []int)
	search = func(idx int, freeVals []int) {
		if idx == len(freeVars) {
			if total := evalFreeVars(coefs, targets, freeVals); total >= 0 {
				if minTotal < 0 || total < minTotal {
					minTotal = total
				}
			}
			return
		}
		for v := 0; v <= maxFree; v++ {
			freeVals[idx] = v
			search(idx+1, freeVals)
		}
	}
	search(0, make([]int, len(freeVars)))

	if minTotal < 0 {
		return 0
	}
	return minTotal
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
