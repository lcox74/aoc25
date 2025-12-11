package main

import "slices"

// buildPolygonMap creates a compressed grid marking cells outside the polygon.
func (m *MovieTheater) buildPolygonMap() (outside [][]bool, xIdx, yIdx map[int]int) {
	xCoords, yCoords, xIdx, yIdx := m.buildCoordinateMaps()
	boundary := m.markBoundary(xCoords, yCoords, xIdx, yIdx)
	outside = floodFillOutside(boundary, len(xCoords), len(yCoords))
	return outside, xIdx, yIdx
}

// buildCoordinateMaps creates sorted coordinate lists and index maps.
func (m *MovieTheater) buildCoordinateMaps() (xCoords, yCoords []int, xIdx, yIdx map[int]int) {
	n := len(m.TilesX)

	xSet := make(map[int]struct{})
	ySet := make(map[int]struct{})
	minX, maxX := m.TilesX[0], m.TilesX[0]
	minY, maxY := m.TilesY[0], m.TilesY[0]

	for i := range n {
		x, y := m.TilesX[i], m.TilesY[i]
		xSet[x] = struct{}{}
		ySet[y] = struct{}{}
		minX, maxX = min(minX, x), max(maxX, x)
		minY, maxY = min(minY, y), max(maxY, y)
	}

	// Add boundary coordinates outside the polygon
	xSet[minX-1] = struct{}{}
	xSet[maxX+1] = struct{}{}
	ySet[minY-1] = struct{}{}
	ySet[maxY+1] = struct{}{}

	xCoords = slices.Sorted(mapKeys(xSet))
	yCoords = slices.Sorted(mapKeys(ySet))

	xIdx = make(map[int]int)
	yIdx = make(map[int]int)
	for i, x := range xCoords {
		xIdx[x] = i
	}
	for i, y := range yCoords {
		yIdx[y] = i
	}

	return xCoords, yCoords, xIdx, yIdx
}

// markBoundary marks polygon edge cells on the compressed grid.
func (m *MovieTheater) markBoundary(xCoords, yCoords []int, xIdx, yIdx map[int]int) [][]bool {
	n := len(m.TilesX)
	boundary := make([][]bool, len(yCoords))
	for i := range boundary {
		boundary[i] = make([]bool, len(xCoords))
	}

	for i := range n {
		x1, y1 := m.TilesX[i], m.TilesY[i]
		x2, y2 := m.TilesX[(i+1)%n], m.TilesY[(i+1)%n]

		if x1 == x2 {
			xi := xIdx[x1]
			yMin, yMax := min(y1, y2), max(y1, y2)
			for _, y := range yCoords {
				if y >= yMin && y <= yMax {
					boundary[yIdx[y]][xi] = true
				}
			}
		} else {
			yi := yIdx[y1]
			xMin, xMax := min(x1, x2), max(x1, x2)
			for _, x := range xCoords {
				if x >= xMin && x <= xMax {
					boundary[yi][xIdx[x]] = true
				}
			}
		}
	}

	return boundary
}

// floodFillOutside marks all cells reachable from outside the polygon using BFS.
func floodFillOutside(boundary [][]bool, width, height int) [][]bool {
	outside := make([][]bool, height)
	for i := range outside {
		outside[i] = make([]bool, width)
	}

	type point struct{ x, y int }
	queue := []point{{0, 0}}
	outside[0][0] = true

	dx := []int{-1, 1, 0, 0}
	dy := []int{0, 0, -1, 1}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for d := range 4 {
			nx, ny := cur.x+dx[d], cur.y+dy[d]
			if nx < 0 || nx >= width || ny < 0 || ny >= height {
				continue
			}
			if outside[ny][nx] || boundary[ny][nx] {
				continue
			}
			outside[ny][nx] = true
			queue = append(queue, point{nx, ny})
		}
	}

	return outside
}

// isInsidePolygon checks if all cells in the rectangle are inside the polygon.
func isInsidePolygon(outside [][]bool, xIdx, yIdx map[int]int, xMin, xMax, yMin, yMax int) bool {
	for yi := yIdx[yMin]; yi <= yIdx[yMax]; yi++ {
		for xi := xIdx[xMin]; xi <= xIdx[xMax]; xi++ {
			if outside[yi][xi] {
				return false
			}
		}
	}
	return true
}

// mapKeys returns an iterator over map keys.
func mapKeys[K comparable, V any](m map[K]V) func(func(K) bool) {
	return func(yield func(K) bool) {
		for k := range m {
			if !yield(k) {
				return
			}
		}
	}
}
