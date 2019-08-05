package pickrand

import (
	"fmt"
	"math/rand"
)

// split64 splits the high and low halfs of a 64-bit uint, returning them as
// a pair of uint64, e.g.,
//
//     var hi, lo uint64 = split64(0x1111111122222222)
//     fmt.Printf("%016x %016x", hi, lo)  // 0000000011111111 0000000022222222
func split64(x uint64) (hi, lo uint64) {
	return x >> 32, x & 0x00000000FFFFFFFF
}

// split64to32 splits the high and low halfs of a 64-bit uint, returning them as
// a pair of uint32, e.g.,
//
//     var hi, lo uint32 = split64to32(0x1111111122222222)
//     fmt.Printf("%08x %08x", hi, lo)  // 11111111 22222222
func split64to32(x uint64) (hi, lo uint32) {
	return uint32(x >> 32), uint32(x)
}

// Uint128 represents a 128-bit unsigned int with a pair of 64-bit unsigned ints.
type Uint128 struct {
	Hi, Lo uint64
}

func (u Uint128) String() string {
	return fmt.Sprintf("%016x%016x", u.Hi, u.Lo)
}

// Mul64 multiplies two 64-bit uints keeping full precision, returning the
// result in a 128-bit uint.
func Mul64(x, y uint64) Uint128 {
	// Full precision is accomplished by working on the halfs of each number,
	// performing multiplications with 32 significant bits (resulting in 64 bit
	// numbers) and addition in 64 bits.
	//
	// x = xHi * 2^32 + xLo
	// y = yHi * 2^32 + yLo
	// x*y = (xHi * 2^32 + xLo) * (yHi * 2^32 + yLo)
	//     = xHi * yHi * 2^64 + (xLo * yHi + xHi * yLo) * 2^32 + xLo * yLo
	//       ---------           ---------   ---------           ---------
	//           d                   c           b                   a

	xHi, xLo := split64(x)
	yHi, yLo := split64(y)

	a := xLo * yLo
	b := xHi * yLo
	c := xLo * yHi
	d := xHi * yHi

	// The result number m is composed by adding each hi and lo part of the
	// summation terms. Note that there may be a carry from the pointed column
	// to the next.
	//
	//                  ↓
	//               [ aHi | aLo ]
	//         [ bHi | bLo ]
	//         [ cHi | cLo ]
	// + [ dHi | dLo ]
	//   -------------------------
	//   [    mHi    |    mLo    ]

	aHi := a >> 32
	bHi, bLo := split64(b)
	cHi, cLo := split64(c)

	carry := (aHi + bLo + cLo) >> 32

	mLo := a + bLo<<32 + cLo<<32
	mHi := d + bHi + cHi + carry

	return Uint128{mHi, mLo}
}

// unoptimizedRandUint32n picks a random uint32 in [0,n), without bias.
//
// The version is simply a direct translation of the algorithm in the package
// documentation, and always performs 1 mod operation. The optimized version
// Uint32n almost never does.
func unoptimizedUint32n(n uint32) uint32 {
	minLo := uint64(1<<32) % uint64(n)
	for {
		x := rand.Uint32()
		m := uint64(x) * uint64(n)
		hi, lo := split64(m)
		if lo >= minLo {
			return uint32(hi)
		}
	}
}

// Uint32n returns a random uint32 in the [0,n) range, without bias.
func Uint32n(n uint32) uint32 {
	x := rand.Uint32()
	m := uint64(x) * uint64(n)
	hi, lo := split64to32(m)
	if lo >= n {
		// Divisionless case: since n > minLo, if lo >= n then it won't be
		// rejected, and we avoid the (relatively) expensive mod operation.
		return hi
	}
	minLo := -n % n // == (2^32-n) % n == 2^32 % n
	for {
		if lo >= minLo {
			return hi
		}
		x = rand.Uint32()
		m = uint64(x) * uint64(n)
		hi, lo = split64to32(m)
	}
}

// Uint64n returns a random uint64 in the [0,n) range, without bias.
func Uint64n(n uint64) uint64 {
	x := rand.Uint64()
	m := Mul64(x, n)
	if m.Lo >= n {
		// Divisionless case: since n > minLo, if m.Lo >= n then it won't be
		// rejected and we avoid the (relatively) expensive mod operation.
		return m.Hi
	}
	minLo := -n % n // == (2^64 - n) % n == 2^64 % n
	for {
		if m.Lo >= minLo {
			return m.Hi
		}
		x = rand.Uint64()
		m = Mul64(x, n)
	}
}
