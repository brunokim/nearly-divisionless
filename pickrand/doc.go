// This package implements the 'nearly divisionless' algorithm for random
// integer generation in an interval. For more details, benchmarks and
// comparisons, see the paper 'Fast Random Integer Generation in an
// Interval', Daniel Lemire (https://dx.doi.org/10.1145/3230636).
//
// Consider that x is a random uint in the interval [0,M), n < M. To get a
// random uint in the [0,n) interval instead, we may wish to do the
// following, considering that '·' and '÷' have infinite precision.
//
//   i = ⌊n·(x÷M)⌋    // (x÷M) is in [0,1)
//     = ⌊(n·x)÷M⌋
//
// However, this would introduce a bias where some numbers appears more
// often then others. The truncated division maps the multiples m = n·x in
// [0,n·M) to [0,n). All multiples in the interval [i*M, (i+1)*M) are
// mapped to i, and some intervals will have more multiples than others.
// This is shown as follows, with M=0xF and n=0x7; if the marked multiples
// were rejected, then all i's would be equiprobable:
//
//    +---+----+---+    +---+----+---+    +---+----+---+    +---+----+---+
//    | x |  m | i |    | x |  m | i |    | x |  m | i |    | x |  m | i |
//    +---+----+---+    +---+----+---+    +---+----+---+    +---+----+---+
//    | 0 | 00*| 0 |    | 4 | 1C | 1 |    | 8 | 38 | 3 |    | C | 54 | 5 |
//    | 1 | 07 | 0 |    | 5 | 24 | 2 |    | 9 | 39 | 3 |    | D | 57 | 5 |
//    | 2 | 0E | 0 |    | 6 | 2A | 2 |    | A | 46 | 4 |    | E | 62 | 6 |
//    | 3 | 15 | 1 |    | 7 | 31*| 3 |    | B | 4D | 4 |    | F | 69 | 6 |
//    +---+----+---+    +---+----+---+    +---+----+---+    +---+----+---+
//
// The solution is to reject multiples that fall in the subinterval
// [i*M, i*M + (M%n)) and accepting only those that fall in the subinterval
// [i*M + (M%n), (i+1)*M). The accepting subinterval has length M-(M%n),
// which is divided exactly by n; there are always M/n multiples within
// them. This is illustrated as follows:
//
//     - M = 0xF, n = 0x7
//     - x in {0x0..0xF} is shown at position (i,j), where
//     - m = x*n = i*0xF + j
//     - that is, i and j are the high and low halfs of m
//
//                          i,j | 0      7 8      F
//                         -----+--------- --------
//            These lines →  0  | 0......1 ......2.
//             have three    1  | .....3.. ....4...
//        multiples while    2  | ...5.... ..6.....
//        others have two →  3  | .7...... 8......9
//                           4  | ......A. .....B..
//                           5  | ....C... ...D....
//                           6  | ..E..... .F......
//                                ↑↑
//                   By rejecting these 2 columns,
//                  we will get an unbiased sample.
//                       Note that 2 = 16 % 7.
//
// With M a power of 2, the division can be translated to a bit shift.
// m needs to be of an extended precision relative to n (if n is an uint32,
// than m is an uint64). A direct translation of this algorithm to code can
// be found in the unoptimizedUint32n function.
package pickrand
