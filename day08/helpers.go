package main

import "sort"

// Edge represents a pair of junction boxes and their squared distance.
type Edge struct {
	I, J   int
	DistSq int64
}

// buildEdges creates all pairwise edges between junction boxes.
func buildEdges(boxes []JunctionBox) []Edge {
	n := len(boxes)
	edges := make([]Edge, 0, n*(n-1)/2)

	for i := range n {
		for j := i + 1; j < n; j++ {
			dx := int64(boxes[i].X - boxes[j].X)
			dy := int64(boxes[i].Y - boxes[j].Y)
			dz := int64(boxes[i].Z - boxes[j].Z)
			edges = append(edges, Edge{I: i, J: j, DistSq: dx*dx + dy*dy + dz*dz})
		}
	}

	sort.Slice(edges, func(a, b int) bool {
		return edges[a].DistSq < edges[b].DistSq
	})

	return edges
}

// initUnionFind initializes parent and rank slices for Union-Find.
func (p *Playground) initUnionFind(n int) {
	p.parent = make([]int, n)
	p.rank = make([]int, n)

	for i := range n {
		p.parent[i] = i
	}
}

// find returns the root of the set containing i, with path compression.
func (p *Playground) find(i int) int {
	if p.parent[i] != i {
		p.parent[i] = p.find(p.parent[i])
	}

	return p.parent[i]
}

// union merges the sets containing i and j. Returns true if they were separate.
func (p *Playground) union(i, j int) bool {
	ri, rj := p.find(i), p.find(j)
	if ri == rj {
		return false
	}

	if p.rank[ri] < p.rank[rj] {
		ri, rj = rj, ri
	}

	p.parent[rj] = ri
	if p.rank[ri] == p.rank[rj] {
		p.rank[ri]++
	}
	return true
}

// topCircuitProduct returns the product of the n largest circuit sizes.
func (p *Playground) topCircuitProduct(n int) int {
	sizes := make(map[int]int)
	for i := range len(p.boxes) {
		sizes[p.find(i)]++
	}
	sizeList := make([]int, 0, len(sizes))
	for _, s := range sizes {
		sizeList = append(sizeList, s)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sizeList)))

	result := 1
	for i := 0; i < n && i < len(sizeList); i++ {
		result *= sizeList[i]
	}
	return result
}
