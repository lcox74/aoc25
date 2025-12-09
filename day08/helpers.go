package main

import "sort"

// Union-Find (Disjoint Set Union) operations

// initUnionFind initializes the Union-Find data structure.
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

// Distance calculations

// distSq returns the squared Euclidean distance between two junction boxes.
func distSq(a, b JunctionBox) int64 {
	dx := int64(a.X - b.X)
	dy := int64(a.Y - b.Y)
	dz := int64(a.Z - b.Z)
	return dx*dx + dy*dy + dz*dz
}

// Circuit size calculations

// computeCircuitSizes counts the size of each connected component.
func (p *Playground) computeCircuitSizes() map[int]int {
	sizes := make(map[int]int)
	for i := range len(p.boxes) {
		root := p.find(i)
		sizes[root]++
	}
	return sizes
}

// topNProduct returns the product of the n largest values in the sizes map.
func topNProduct(sizes map[int]int, n int) int {
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

// Edge building

// buildAllEdges creates edges between all pairs of junction boxes.
func (p *Playground) buildAllEdges() []Edge {
	n := len(p.boxes)
	edges := make([]Edge, 0, n*(n-1)/2)
	for i := range n {
		for j := i + 1; j < n; j++ {
			edges = append(edges, Edge{I: i, J: j, DistSq: distSq(p.boxes[i], p.boxes[j])})
		}
	}
	return edges
}

// sortEdgesByDistance sorts edges by squared distance ascending.
func sortEdgesByDistance(edges []Edge) {
	sort.Slice(edges, func(a, b int) bool {
		return edges[a].DistSq < edges[b].DistSq
	})
}
