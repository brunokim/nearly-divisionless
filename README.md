# Nearly-divisionless algorithm for picking an integer in an interval

This package is my entry in [@colmmacc's contest](https://twitter.com/colmmacc/status/1153727018896244736) for a readable
implementation of Daniel Lemire's algorithm for picking a random number between 0 and a given N.

- `pickrand` contains the documentation, test and implementation for the algorithm. [Godoc](https://godoc.org/github.com/brunokim/nearly-divisionless/pickrand)
  is the best option to read the docs; I hope I did a good job in distilling the [original paper](https://dx.doi.org/10.1145/3230636)
  and why it works.
- `main.go` has a sample usage for `pickrand.Uint64`: generating a random network using [Barabási-Albert's model](https://en.wikipedia.org/wiki/Barab%C3%A1si%E2%80%93Albert_model),
  where most nodes have few connections and a few nodes have MOST connections. I also compute a vertex cover for the generated
  network, so it's clear that from a few nodes one can reach all others in one hop.
- `stats` has an attempt on creating statistical tests to validate that the produced numbers are indeed random and not biased.
  I can't say that I was successful, because even a clearly biased algorithm (simply `$RANDOM % $N`) didn't produce enough of
  a signal for [Kolmogorov-Smirnov](https://en.wikipedia.org/wiki/Kolmogorov%E2%80%93Smirnov_test) or Cramer-Von Mises tests,
  in the ranges that I tried, to clearly indicate a biased distribution. Most likely I did something wrong, don't use this code.

Enjoy!
