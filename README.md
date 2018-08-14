# mph

[![GoDoc](https://godoc.org/github.com/cespare/mph?status.svg)](https://godoc.org/github.com/cespare/mph)

mph is a Go package for that implements a
[minimal perfect hash table](https://en.wikipedia.org/wiki/Perfect_hash_function#Minimal_perfect_hash_function)
over strings. It uses the
["Hash, displace, and compress" algorithm](http://cmph.sourceforge.net/papers/esa09.pdf)
and the [Murmur3 hash function](https://en.wikipedia.org/wiki/MurmurHash).

Some quick benchmark results (this is on a 2.5 GHz Skylake CPU):

* `Build` constructs a minimal perfect hash table from a 99k word dictionary in
  34ms (construction time is linear in the size of the input).
* `Lookup`s on that dictionary take <50ns and are 28% faster than a
  `map[string]uint32`:

    ```
    BenchmarkTable          300000000               48.6 ns/op
    BenchmarkTableMap       200000000               67.6 ns/op
    ```
