package stats

import (
	"math"
	"sort"
)

type point struct {
	x, y float64
}

// Compute the relative frequency of the given samples in the [0, maximum) domain.
func pdf(samples []uint64, maximum uint64) []point {
	total := len(samples)
	// Count number of occurrences per sample.
	freq := make(map[uint64]int)
	for _, x := range samples {
		freq[x]++
	}
	// Normalize frequency and range to obtain a PDF over [0,1).
	arr := make([]point, len(freq))
	i := 0
	for x, count := range freq {
		arr[i] = point{float64(x+1) / float64(maximum), float64(count) / float64(total)}
		i++
	}
	return arr
}

// Compute the cumulative density function of the given samples in the [0, maximum) domain.
func cdf(samples []uint64, maximum uint64) []point {
	arr := pdf(samples, maximum)
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].x < arr[j].x
	})
	cumulative := make([]point, len(arr))
	acc := 0.0
	for i, e := range arr {
		cumulative[i] = point{e.x, acc + e.y}
		acc += e.y
	}
	return cumulative
}

// The Kolmorogov-Smirnov (K-S) statistic computes the maximum y difference
// between a sampled CDF and its expected CDF, adjusted by the (sqrt of) number of samples.
//
// In this case, the expected CDF is y = x, for the uniform distribution in [0,1].
func KSUnitUniformStatistic(samples []uint64, maximum uint64) float64 {
	n := len(samples)
	cdf := cdf(samples, maximum)

	maxDiff := 0.0
	for _, f := range cdf {
		diff := math.Abs(f.y - f.x) // f.y = sampled, f.x = expected
		if diff > maxDiff {
			maxDiff = diff
		}
	}
	return math.Sqrt(float64(n)) * maxDiff
}
