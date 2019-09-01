package pickrand_test

import (
	"testing"
	"testing/quick"

	"brunokim.xyz/pickrand"
	"brunokim.xyz/stats"
)

// Naive, inefficient, but probably correct multiplication implementation.
func naiveMul64(x, y uint64) pickrand.Uint128 {
	m := pickrand.Uint128{0, 0}

	for i := uint(0); i < 64; i++ {
		lsb := y % 2
		y = y >> 1
		if lsb == 0 {
			continue
		}
		newLo := m.Lo + (x << i)
		if newLo < m.Lo {
			// Overflow; there's a carry from lo to hi.
			m.Hi += 1
		}
		m.Lo = newLo

		m.Hi += x >> (64 - i)
	}
	return m
}

func TestMul64(t *testing.T) {
	tests := []struct {
		desc string
		x, y uint64
		want pickrand.Uint128
	}{
		{"Minimum       ", 0x0000000000000001, 0x0000000000000001, pickrand.Uint128{0x0000000000000000, 0x0000000000000001}},
		{"Last just low ", 0x0000000000000002, 0x7FFFFFFFFFFFFFFF, pickrand.Uint128{0x0000000000000000, 0xFFFFFFFFFFFFFFFE}},
		{"First w/ high ", 0x0000000000000002, 0x8000000000000000, pickrand.Uint128{0x0000000000000001, 0x0000000000000000}},
		{"Maximum       ", 0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF, pickrand.Uint128{0xFFFFFFFFFFFFFFFE, 0x0000000000000001}},
	}

	for _, test := range tests {
		got := pickrand.Mul64(test.x, test.y)
		if got != test.want {
			t.Errorf("%s: %016x * %016x = %v (want %v)", test.desc, test.x, test.y, got, test.want)
		}
		got = naiveMul64(test.x, test.y)
		if got != test.want {
			t.Errorf("%s (naive): %016x * %016x = %v (want %v)", test.desc, test.x, test.y, got, test.want)
		}
	}
}

func TestQuickMul64(t *testing.T) {
	if err := quick.CheckEqual(pickrand.Mul64, naiveMul64, nil); err != nil {
		err := err.(*quick.CheckEqualError)
		x, y := err.In[0], err.In[1]
		out1, out2 := err.Out1[0], err.Out2[0]
		t.Errorf("quick check #%d failure: pickrand.Mul64(%016x, %016x) = %v, naiveMul64(%016x, %016x) = %v", err.Count, x, y, out1, x, y, out2)
	}
}

// ksThreshold is the maximum value of the K-S statistic we accept from a sample
// of 1000 values to be still classified as uniform, with 99.9% confidence.
//
// This thresold was obtained empirically, by running 100,000 trials where on
// each one we generated a list of 1000 random integers between [0, 2^32) using
// Python's "random" package, and then computed their K-S statistic. The distribution
// of the statistic is summarized in the following percentiles:
//
//     p25:   0.655
//     p50:   0.806
//     p90:   1.20
//     p99:   1.61
//     p99.9: 1.95
const ksThreshold = 2.0

func TestUint32n(t *testing.T) {
	tests := []struct {
		desc string
		s    uint32
	}{
		{"Constant", 1},
		{"1 bit", 2},
		{"1/2 of space", 1<<31 + 1},
		{"2/3 of space", 0xAAAAAAAA},
	}

	for _, test := range tests {
		n := 1000
		samples := make([]uint64, n)
		for i := 0; i < n; i++ {
			samples[i] = uint64(pickrand.Uint32n(test.s))
		}
		ks := stats.KSUniformStatistic(samples, uint64(test.s))
		t.Logf("%d samples from [0, %#x) K-S statistic = %.4f\n", n, test.s, ks)

		if ks > ksThreshold {
			for _, p := range stats.Frequency(samples) {
				t.Logf("%.0f: %.4f", p.X, p.Y)
			}
			t.Errorf("%s: K-S test failed", test.desc)
		}
	}
}

func TestUint64n(t *testing.T) {
	tests := []struct {
		desc string
		s    uint64
	}{
		{"Constant", 1},
		{"1 bit", 2},
		{"1/2 of space", 1<<63 + 1},
		{"2/3 of space", 0xAAAAAAAAAAAAAAAA},
	}

	for _, test := range tests {
		n := 1000
		samples := make([]uint64, n)
		for i := 0; i < n; i++ {
			samples[i] = pickrand.Uint64n(test.s)
		}
		ks := stats.KSUniformStatistic(samples, uint64(test.s))
		t.Logf("%d samples from [0, %#x) K-S statistic = %.4f\n", n, test.s, ks)

		if ks > ksThreshold {
			for _, p := range stats.Frequency(samples) {
				t.Logf("%.0f: %.4f", p.X, p.Y)
			}
			t.Errorf("%s: K-S test failed", test.desc)
		}
	}
}
