#!/usr/bin/env python3

import math
import random

from collections import Counter

def samples(n, s, bias=0):
    for _ in range(n):
        yield random.randrange(s % (s+1 - bias))


def samples_float(n, s, mode="floor"):
    rounding = {
        "floor": math.floor,
        "ceil": math.ceil,
        "round": round,
    }[mode]
    for _ in range(n):
        yield int(rounding(random.random() * s))


def cumulative(xs, s):
    total = len(xs)
    freq = sorted(Counter(xs).items())

    acc = 0
    yield 0.0, 0.0
    for x, y in freq:
        acc += y
        yield (x+1)/s, acc/total
    yield 1.0, 1.0
    

def ks_dist_unif(cum_hist):
    return max(abs(y-x) for x, y in cum_hist)


def cramer_von_mises_dist_unif(cum_hist):
    n = len(cum_hist)
    diff_squared = [(x, (y-x)**2) for x, y, in cum_hist]
    
    # Numerical integration
    s = 0.0
    for i in range(1, n-1):
        x1, y1 = diff_squared[i-1]
        x2, y2 = diff_squared[i]
        s += (x2-x1) * (y2+y1) / 2
    return math.sqrt(s)


def ks_metric(xs, s):
    cum_hist = list(cumulative(xs, s))
    dist = ks_dist_unif(cum_hist)
    return dist * math.sqrt(len(xs))


def cvm_metric(xs, s):
    cum_hist = list(cumulative(xs, s))
    dist = cramer_von_mises_dist_unif(cum_hist)
    return dist * math.sqrt(len(xs))


def main():
    n = 1000
    s = 1 << 32
    for bias in range(0, 32, 4):
        print(f"bias: {bias}")
        xs = list(samples(n, s, bias))

        for sqrt in range(2, int(math.sqrt(n)+1)):
            i = sqrt*sqrt
            ks = ks_metric(xs[:i], s)
            cvm = cvm_metric(xs[:i], s)
            print(f"{i:>4} {ks:>.3} {cvm:>.3}")
        print()
    
    for mode in ["round", "floor", "ceil"]:
        print(f"float bias: {mode}")
        xs = list(samples_float(n, s, mode))

        for sqrt in range(2, int(math.sqrt(n)+1)):
            i = sqrt*sqrt
            ks = ks_metric(xs[:i], s)
            cvm = cvm_metric(xs[:i], s)
            print(f"{i:>4} {ks:>.3} {cvm:>.3}")
        print()

    ks_metrics, cvm_metrics = [], []
    trials = 100000
    for _ in range(trials):
        ks_metrics.append(ks_metric(list(samples(n, s)), s))
        cvm_metrics.append(cvm_metric(list(samples(n, s)), s))
    ks_metrics.sort()
    print(f"p25: {ks_metrics[25*trials//100]:.3}")
    print(f"p50: {ks_metrics[50*trials//100]:.3}")
    print(f"p90: {ks_metrics[90*trials//100]:.3}")
    print(f"p99: {ks_metrics[99*trials//100]:.3}")
    print(f"p99.9: {ks_metrics[999*trials//1000]:.3}")
    
    cvm_metrics.sort()
    print(f"p25: {cvm_metrics[25*trials//100]:.3}")
    print(f"p50: {cvm_metrics[50*trials//100]:.3}")
    print(f"p90: {cvm_metrics[90*trials//100]:.3}")
    print(f"p99: {cvm_metrics[99*trials//100]:.3}")
    print(f"p99.9: {cvm_metrics[999*trials//1000]:.3}")


if __name__ == '__main__':
    main()
