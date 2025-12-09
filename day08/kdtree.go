package main

import (
	"container/heap"
	"sort"
)

// KDTree is a 3-dimensional k-d tree for efficient nearest neighbor queries.
// A k-d tree recursively partitions space by cycling through dimensions.
// At depth 0, we split on X; depth 1 on Y; depth 2 on Z; depth 3 on X again, etc.
type KDTree struct {
	root *kdNode
}

// kdNode represents a node in the KD-Tree.
type kdNode struct {
	point JunctionBox // The point stored at this node
	index int         // Original index in the input slice
	left  *kdNode     // Points with smaller value in split dimension
	right *kdNode     // Points with larger value in split dimension
}

// indexedPoint pairs a point with its original index for tree construction.
type indexedPoint struct {
	point JunctionBox
	index int
}

// NewKDTree builds a KD-Tree from a slice of junction boxes.
// Time complexity: O(n log n) for balanced tree construction.
func NewKDTree(boxes []JunctionBox) *KDTree {
	if len(boxes) == 0 {
		return &KDTree{}
	}

	// Create indexed points to track original indices
	points := make([]indexedPoint, len(boxes))
	for i, b := range boxes {
		points[i] = indexedPoint{point: b, index: i}
	}

	return &KDTree{
		root: buildKDTree(points, 0),
	}
}

// buildKDTree recursively constructs the tree.
// depth determines which axis to split on (depth % 3).
func buildKDTree(points []indexedPoint, depth int) *kdNode {
	if len(points) == 0 {
		return nil
	}

	// Choose axis based on depth: 0=X, 1=Y, 2=Z
	axis := depth % 3

	// Sort points by the current axis
	sort.Slice(points, func(i, j int) bool {
		return getAxis(points[i].point, axis) < getAxis(points[j].point, axis)
	})

	// Select median as pivot
	median := len(points) / 2

	return &kdNode{
		point: points[median].point,
		index: points[median].index,
		left:  buildKDTree(points[:median], depth+1),
		right: buildKDTree(points[median+1:], depth+1),
	}
}

// getAxis returns the value of the specified axis (0=X, 1=Y, 2=Z).
func getAxis(p JunctionBox, axis int) int {
	switch axis {
	case 0:
		return p.X
	case 1:
		return p.Y
	default:
		return p.Z
	}
}

// Neighbor represents a potential nearest neighbor with its distance.
type Neighbor struct {
	Index  int
	DistSq int64
}

// KNNBuffer holds reusable state for k-nearest neighbor queries.
// This avoids allocations by reusing the same buffer across queries.
type KNNBuffer struct {
	neighbors []Neighbor // sorted by distance descending (farthest first)
	count     int        // number of valid neighbors in buffer
	k         int        // max neighbors to track
}

// NewKNNBuffer creates a buffer for k-nearest neighbor queries.
func NewKNNBuffer(k int) *KNNBuffer {
	return &KNNBuffer{
		neighbors: make([]Neighbor, k),
		k:         k,
	}
}

// Reset clears the buffer for a new query.
func (b *KNNBuffer) Reset() {
	b.count = 0
}

// Results returns the neighbors found, sorted by distance ascending (closest first).
func (b *KNNBuffer) Results() []Neighbor {
	// Buffer is sorted descending, reverse for ascending order
	result := b.neighbors[:b.count]
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result
}

// MaxDist returns the distance to the farthest neighbor, or max int64 if not full.
func (b *KNNBuffer) MaxDist() int64 {
	if b.count < b.k {
		return 1<<63 - 1
	}
	return b.neighbors[0].DistSq
}

// TryAdd attempts to add a neighbor. Returns true if added.
// Uses insertion sort to maintain descending order by distance.
func (b *KNNBuffer) TryAdd(index int, distSq int64) bool {
	// If buffer not full, always add
	if b.count < b.k {
		// Insert in sorted position (descending by distance)
		pos := b.count
		for pos > 0 && b.neighbors[pos-1].DistSq < distSq {
			pos--
		}
		// Shift elements to make room
		copy(b.neighbors[pos+1:b.count+1], b.neighbors[pos:b.count])
		b.neighbors[pos] = Neighbor{Index: index, DistSq: distSq}
		b.count++
		return true
	}

	// Buffer full - only add if closer than farthest
	if distSq >= b.neighbors[0].DistSq {
		return false
	}

	// Find insertion position (descending order)
	pos := 1
	for pos < b.k && b.neighbors[pos].DistSq > distSq {
		pos++
	}
	// Shift elements left (dropping the farthest at index 0)
	copy(b.neighbors[0:pos-1], b.neighbors[1:pos])
	b.neighbors[pos-1] = Neighbor{Index: index, DistSq: distSq}
	return true
}

// maxHeap implements a max-heap of neighbors for k-nearest search.
// We use a max-heap so we can efficiently remove the farthest neighbor
// when we find a closer one.
type maxHeap []Neighbor

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].DistSq > h[j].DistSq } // Max-heap
func (h maxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *maxHeap) Push(x any) {
	*h = append(*h, x.(Neighbor))
}

func (h *maxHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

// KNearestInto finds the k nearest neighbors using a reusable buffer.
// The buffer is reset and filled with results. Returns slice view into buffer.
// This method avoids allocations by reusing the buffer across queries.
func (t *KDTree) KNearestInto(query JunctionBox, excludeIndex int, buf *KNNBuffer) []Neighbor {
	if t.root == nil || buf.k <= 0 {
		return nil
	}

	buf.Reset()
	t.kNearestSearchBuf(t.root, query, 0, buf, excludeIndex)
	return buf.Results()
}

// KNearest finds the k nearest neighbors to the query point.
// Returns neighbors sorted by distance (closest first).
// Time complexity: O(k log k + log n) average case.
//
// The algorithm:
//  1. Maintain a max-heap of size k with the best candidates so far
//  2. Recursively search the tree, always exploring the side containing the query first
//  3. Only explore the other side if it could contain closer points
//     (i.e., if the splitting plane is closer than our k-th best distance)
func (t *KDTree) KNearest(query JunctionBox, k, excludeIndex int) []Neighbor {
	if t.root == nil || k <= 0 {
		return nil
	}

	h := &maxHeap{}
	heap.Init(h)

	t.kNearestSearch(t.root, query, k, 0, h, excludeIndex)

	// Convert heap to sorted slice (closest first)
	result := make([]Neighbor, h.Len())
	for i := len(result) - 1; i >= 0; i-- {
		result[i] = heap.Pop(h).(Neighbor)
	}
	return result
}

// kNearestSearch recursively searches for k nearest neighbors.
func (t *KDTree) kNearestSearch(node *kdNode, query JunctionBox, k, depth int, h *maxHeap, excludeIndex int) {
	if node == nil {
		return
	}

	// Calculate distance to current node
	if node.index != excludeIndex {
		d := distSq(query, node.point)

		if h.Len() < k {
			// Haven't found k neighbors yet, add this one
			heap.Push(h, Neighbor{Index: node.index, DistSq: d})
		} else if d < (*h)[0].DistSq {
			// Found a closer neighbor than our current k-th best
			heap.Pop(h)
			heap.Push(h, Neighbor{Index: node.index, DistSq: d})
		}
	}

	// Determine which side to search first
	axis := depth % 3
	queryVal := getAxis(query, axis)
	nodeVal := getAxis(node.point, axis)

	// Search the side containing the query point first
	var first, second *kdNode
	if queryVal < nodeVal {
		first, second = node.left, node.right
	} else {
		first, second = node.right, node.left
	}

	t.kNearestSearch(first, query, k, depth+1, h, excludeIndex)

	// Check if we need to search the other side
	// We only need to if the splitting plane is closer than our k-th best distance
	// (or if we haven't found k neighbors yet)
	planeDist := int64(queryVal - nodeVal)
	planeDistSq := planeDist * planeDist

	if h.Len() < k || planeDistSq < (*h)[0].DistSq {
		t.kNearestSearch(second, query, k, depth+1, h, excludeIndex)
	}
}

// kNearestSearchBuf recursively searches using the buffer instead of heap.
func (t *KDTree) kNearestSearchBuf(node *kdNode, query JunctionBox, depth int, buf *KNNBuffer, excludeIndex int) {
	if node == nil {
		return
	}

	// Calculate distance to current node
	if node.index != excludeIndex {
		d := distSq(query, node.point)
		buf.TryAdd(node.index, d)
	}

	// Determine which side to search first
	axis := depth % 3
	queryVal := getAxis(query, axis)
	nodeVal := getAxis(node.point, axis)

	// Search the side containing the query point first
	var first, second *kdNode
	if queryVal < nodeVal {
		first, second = node.left, node.right
	} else {
		first, second = node.right, node.left
	}

	t.kNearestSearchBuf(first, query, depth+1, buf, excludeIndex)

	// Check if we need to search the other side
	planeDist := int64(queryVal - nodeVal)
	planeDistSq := planeDist * planeDist

	if planeDistSq < buf.MaxDist() {
		t.kNearestSearchBuf(second, query, depth+1, buf, excludeIndex)
	}
}

// EdgeCandidate represents a potential edge between two junction boxes.
type EdgeCandidate struct {
	I, J   int
	DistSq int64
}

// SolveWithKDTree uses a KD-Tree for spatial queries instead of brute force.
// This approach:
// 1. Builds a KD-Tree from all points - O(n log n)
// 2. For each point, finds k nearest neighbors - O(n * k log n) average
// 3. Collects candidate edges and processes them with Union-Find
//
// For small k relative to n, this is much faster than O(n^2) brute force.
func (p *Playground) SolveWithKDTree(numConnections int) {
	n := len(p.boxes)
	if n == 0 {
		return
	}

	tree := NewKDTree(p.boxes)
	p.initUnionFind(n)

	k := min(50, n-1)
	edges := p.collectCandidateEdges(tree, k)

	// Process edges in order of distance
	connected := 0
	for _, e := range edges {
		if connected >= numConnections {
			break
		}
		p.union(e.I, e.J)
		connected++
	}

	sizes := p.computeCircuitSizes()
	p.ResultPart1 = topNProduct(sizes, 3)
}

// collectCandidateEdges finds candidate edges using KNN queries.
func (p *Playground) collectCandidateEdges(tree *KDTree, k int) []EdgeCandidate {
	buf := NewKNNBuffer(k)
	edgeSet := make(map[[2]int]int64)

	for i, box := range p.boxes {
		neighbors := tree.KNearestInto(box, i, buf)
		for _, nb := range neighbors {
			key := normalizeEdgeKey(i, nb.Index)
			if existing, ok := edgeSet[key]; !ok || nb.DistSq < existing {
				edgeSet[key] = nb.DistSq
			}
		}
	}

	edges := make([]EdgeCandidate, 0, len(edgeSet))
	for key, d := range edgeSet {
		edges = append(edges, EdgeCandidate{I: key[0], J: key[1], DistSq: d})
	}
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].DistSq < edges[j].DistSq
	})

	return edges
}

// normalizeEdgeKey returns edge key with smaller index first.
func normalizeEdgeKey(i, j int) [2]int {
	if i < j {
		return [2]int{i, j}
	}
	return [2]int{j, i}
}
