package stats_test

import (
	"math"
	"math/rand"
	"testing"
	"testing/quick"

	"brunokim.xyz/stats"
)

const ksThreshold = 2

func TestKSStatistic(t *testing.T) {
	n := 1000
	samples := make([]uint64, n)
	for i := 0; i < n; i++ {
		samples[i] = rand.Uint64()
	}
	if ks := stats.KSUniformStatistic(samples, math.MaxUint64); ks > ksThreshold {
		for _, p := range stats.Frequency(samples) {
			t.Logf("%.0f: %.4f", p.X, p.Y)
		}
		t.Errorf("K-S test failed with statistic %.4f", ks)
	}
}

func TestQuickFailingKSStatistic(t *testing.T) {
	const n = 1000
	const runs = 10
	samples := make([]uint64, n)

	test := func(input uint64) bool {
		max := 3 * input / 4
		if max == 0 {
			// Avoid an unlikely division by zero.
			return true
		}
		for run := 0; run < runs; run++ {
			// Biased generator with simple modulus operation.
			for i := 0; i < n; i++ {
				samples[i] = rand.Uint64() % max
			}
			ks := stats.KSUniformStatistic(samples, max)
			t.Logf("max: %d, ks: %.4f", max, ks)
			if ks > 0.8 {
				return true
			}
		}
		return false
	}
	if err := quick.Check(test, nil); err != nil {
		inputs := err.(*quick.CheckError).In
		input := inputs[0].(uint64)
		t.Errorf("K-S test passed %d times for biased generator with max %d", runs, 3*input/4)
	}
}
