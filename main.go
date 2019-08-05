package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"brunokim.xyz/pickrand"
)

// Edge represents a pair of vertices in a graph.
type Edge [2]uint64

var (
	n = flag.Uint64("n", 1000, "Number of vertices in the graph")
	k = flag.Uint64("k", 5, "Minimum number of edges per vertex")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	fmt.Printf("n = %d, k = %d\n\n", *n, *k)
	edges := scaleFreeNetwork(*n, *k)

	degrees := make([]int, *n)
	for _, edge := range edges {
		degrees[edge[0]]++
		degrees[edge[1]]++
	}

	degreeFreq := make([]int, *n)
	for _, degree := range degrees {
		degreeFreq[degree]++
	}

	w1 := len("degree")
	w2 := len("frequency")
	fmt.Println("+--------+-----------+")
	fmt.Println("| degree | frequency |")
	fmt.Println("+--------+-----------+")
	for degree, freq := range degreeFreq {
		if freq == 0 {
			continue
		}
		fmt.Printf("| %*d | %*d |\n", w1, degree, w2, freq)
	}
	fmt.Println("+--------+-----------+")

	coverSet := cover(edges)
	fmt.Println()
	fmt.Printf("%d vertices cover the entire graph:\n", len(coverSet))
	entries := sortCover(coverSet)
	fmt.Printf("  Vertex #%d covers %d vertices\n", entries[0].key, entries[0].val)
	for _, entry := range entries[1:] {
		fmt.Printf("  Vertex #%d covers %d other vertices\n", entry.key, entry.val)
	}
}

// Create a random scale-free network with n vertices and k*n edges.
//
// The degree distribution of a scale-free network follows a power law, with
// remarkable inequality in edge assignment. The algorithm used for this network
// is the Barabási-Albert model, with resulting degree distribution of P(k) ~ k^-3
// (or: there are 1000x less vertices with degree 10*d than vertices with degree d).
func scaleFreeNetwork(n, k uint64) []Edge {
	edges := make([]Edge, n*k)
	i := uint64(0)

	// Create initial (k+1)-clique
	for u := uint64(0); u <= k; u++ {
		for v := u + 1; v <= k; v++ {
			edges[i] = Edge{u, v}
			i++
		}
	}

	// For each new vertex, add k edges to it. The opposite endpoint is selected
	// by picking a random existing edge, and then selecting one of its
	// endpoints at random. This favors picking vertices that already have large
	// degree.
	for u := k + 1; u < n; u++ {
		numEdges := i
		alreadySelected := make(map[uint64]bool)
		for cnt := uint64(0); cnt < k; cnt++ {
			var v uint64
			for {
				edgeIdx := pickrand.Uint64n(numEdges)
				edge := edges[edgeIdx]
				endpoint := pickrand.Uint64n(2)
				v = edge[endpoint]
				if !alreadySelected[v] {
					break
				}
			}
			selected[v] = true
			edges[i] = Edge{u, v}
			i++
		}
	}

	return edges
}

// Greedy algorithm to obtain a set cover.
func cover(edges []Edge) map[uint64]int {
	neighbors := func(u uint64) []uint64 {
		vs := make([]uint64, 0)
		for _, edge := range edges {
			if edge[0] == u {
				vs = append(vs, edge[1])
			}
			if edge[1] == u {
				vs = append(vs, edge[0])
			}
		}
		return vs
	}

	covered := make(map[uint64]bool)
	coverCount := func() map[uint64]int {
		covering := make(map[uint64]int)
		for _, edge := range edges {
			if covered[edge[0]] && covered[edge[1]] {
				continue
			}
			if !covered[edge[0]] {
				covering[edge[1]]++
			}
			if !covered[edge[1]] {
				covering[edge[0]]++
			}
		}
		return covering
	}

	coverSet := make(map[uint64]int)
	for {
		u, k := argmax(coverCount())
		if k == 0 {
			return coverSet
		}
		coverSet[u] = k

		covered[u] = true
		for _, v := range neighbors(u) {
			covered[v] = true
		}
	}
}

func argmax(m map[uint64]int) (uint64, int) {
	imax := uint64(0)
	for i, k := range m {
		if k > m[imax] {
			imax = i
		}
	}
	return imax, m[imax]
}

type entry struct {
	key uint64
	val int
}

func sortCover(m map[uint64]int) []entry {
	a := make([]entry, len(m))
	i := 0
	for k, v := range m {
		a[i] = entry{k, v}
		i++
	}
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].key < a[j].key
	})
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].val >= a[j].val
	})
	return a
}
