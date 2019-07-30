package pickrand_test

import (
	"brunokim.xyz/pickrand"
	"testing"
)

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
	}
}
