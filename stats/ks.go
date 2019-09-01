package stats

import (
	"math"
	"sort"
)

// Point represents a cartesian point.
type Point struct {
	X, Y float64
}

// Frequency returns the relative frequency of the given samples, sorted by x.
func Frequency(samples []uint64) []Point {
	total := len(samples)
	// Count occurrences in array.
	freq := make(map[uint64]int)
	for _, x := range samples {
		freq[x]++
	}
	// Compute relative frequency.
	arr := make([]Point, len(freq))
	i := 0
	for x, count := range freq {
		arr[i] = Point{float64(x + 1), float64(count) / float64(total)}
		i++
	}
	// Sort slice for better presentation, and to ease computing the CDF.
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].X < arr[j].X
	})
	return arr
}

// Cumulative returns the cumulative frequency of the given samples, sorted by x.
func Cumulative(samples []uint64) []Point {
	arr := Frequency(samples)
	cumulative := make([]Point, len(arr))
	acc := 0.0
	for i, p := range arr {
		cumulative[i] = Point{p.X, acc + p.Y}
		acc += p.Y
	}
	return cumulative
}

// The Kolmorogov-Smirnov (K-S) statistic computes the maximum y difference
// between a sampled CDF and its expected CDF, adjusted by the (sqrt of) number of samples.
//
// In this case, the expected CDF is y = x/max, for the uniform distribution in [0, max].
func KSUniformStatistic(samples []uint64, max uint64) float64 {
	n := len(samples)
	cdf := Cumulative(samples)

	maxDiff := 0.0
	for _, f := range cdf {
		expected := f.X / float64(max)
		diff := math.Abs(f.Y - expected)
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	return math.Sqrt(float64(n)) * maxDiff
}
