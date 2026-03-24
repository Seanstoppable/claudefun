package main

import (
	"hash/fnv"
	"math"
	"sort"
)

// Edge represents a line connecting two stars in the constellation.
type Edge struct {
	From, To int     // indices into StarMap.Stars
	Length   float64 // Euclidean distance
}

// Constellation holds the star map and the edges that connect them.
type Constellation struct {
	Stars *StarMap
	Edges []Edge
}

// Connect builds constellation lines between stars using MST + artistic extras.
func Connect(sm *StarMap) *Constellation {
	c := &Constellation{Stars: sm}
	n := len(sm.Stars)
	if n < 2 {
		return c
	}

	// Build MST using Prim's algorithm.
	mstEdges := primMST(sm.Stars)
	c.Edges = mstEdges

	// Add 1-2 artistic edges to create small triangles/loops.
	addArtisticEdges(c, sm)

	return c
}

// primMST returns a minimum spanning tree over the stars using Prim's algorithm.
func primMST(stars []Star) []Edge {
	n := len(stars)
	if n < 2 {
		return nil
	}

	inMST := make([]bool, n)
	// cheapest[i] = lowest distance from star i to any star already in the MST.
	cheapest := make([]float64, n)
	cheapFrom := make([]int, n) // which MST node offered the cheapest edge
	for i := range cheapest {
		cheapest[i] = math.Inf(1)
		cheapFrom[i] = -1
	}

	// Start from star 0.
	inMST[0] = true
	for i := 1; i < n; i++ {
		d := starDist(stars[0], stars[i])
		cheapest[i] = d
		cheapFrom[i] = 0
	}

	edges := make([]Edge, 0, n-1)

	for count := 0; count < n-1; count++ {
		// Pick the non-MST vertex with the smallest edge to the MST.
		next := -1
		best := math.Inf(1)
		for i := 0; i < n; i++ {
			if !inMST[i] && cheapest[i] < best {
				best = cheapest[i]
				next = i
			}
		}
		if next == -1 {
			break
		}

		inMST[next] = true
		edges = append(edges, Edge{
			From:   cheapFrom[next],
			To:     next,
			Length: best,
		})

		// Update cheapest edges for remaining vertices.
		for i := 0; i < n; i++ {
			if !inMST[i] {
				d := starDist(stars[next], stars[i])
				if d < cheapest[i] {
					cheapest[i] = d
					cheapFrom[i] = next
				}
			}
		}
	}

	return edges
}

// addArtisticEdges adds 1-2 short non-MST edges to create visual loops.
func addArtisticEdges(c *Constellation, sm *StarMap) {
	n := len(sm.Stars)
	if n < 3 {
		return
	}

	// Compute average MST edge length.
	avgLen := 0.0
	for _, e := range c.Edges {
		avgLen += e.Length
	}
	avgLen /= float64(len(c.Edges))
	threshold := avgLen * 1.5

	// Build set of existing edges for fast lookup.
	existing := make(map[[2]int]bool)
	for _, e := range c.Edges {
		a, b := e.From, e.To
		if a > b {
			a, b = b, a
		}
		existing[[2]int{a, b}] = true
	}

	// Gather candidate edges: not in MST and shorter than threshold.
	type candidate struct {
		edge Edge
	}
	var candidates []candidate
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if existing[[2]int{i, j}] {
				continue
			}
			d := starDist(sm.Stars[i], sm.Stars[j])
			if d < threshold {
				candidates = append(candidates, candidate{
					edge: Edge{From: i, To: j, Length: d},
				})
			}
		}
	}

	if len(candidates) == 0 {
		return
	}

	// Sort by length so shorter edges are preferred.
	sort.Slice(candidates, func(a, b int) bool {
		return candidates[a].edge.Length < candidates[b].edge.Length
	})

	// Deterministic RNG from the input to pick 1-2 extras.
	h := fnv.New64a()
	h.Write([]byte(sm.Seed))
	seed := h.Sum64()

	extraCount := 1 + int(seed%2) // 1 or 2 extras
	if extraCount > len(candidates) {
		extraCount = len(candidates)
	}

	// Pick edges using the seeded value to skip around the sorted list.
	for i := 0; i < extraCount; i++ {
		idx := int((seed / uint64(i+2)) % uint64(len(candidates)))
		c.Edges = append(c.Edges, candidates[idx].edge)

		// Remove picked candidate to avoid duplicates.
		candidates[idx] = candidates[len(candidates)-1]
		candidates = candidates[:len(candidates)-1]
		if len(candidates) == 0 {
			break
		}
	}
}

func starDist(a, b Star) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// BoundingBox returns the axis-aligned bounding box of the constellation's stars.
func (c *Constellation) BoundingBox() (minX, minY, maxX, maxY float64) {
	if len(c.Stars.Stars) == 0 {
		return 0, 0, 0, 0
	}
	minX, minY = math.Inf(1), math.Inf(1)
	maxX, maxY = math.Inf(-1), math.Inf(-1)
	for _, s := range c.Stars.Stars {
		if s.X < minX {
			minX = s.X
		}
		if s.Y < minY {
			minY = s.Y
		}
		if s.X > maxX {
			maxX = s.X
		}
		if s.Y > maxY {
			maxY = s.Y
		}
	}
	return
}

// StarCount returns the number of stars in the constellation.
func (c *Constellation) StarCount() int {
	return len(c.Stars.Stars)
}

// EdgeCount returns the number of edges (lines) in the constellation.
func (c *Constellation) EdgeCount() int {
	return len(c.Edges)
}
