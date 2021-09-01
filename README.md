# mph


[![Go Reference](https://pkg.go.dev/badge/github.com/cespare/mph.svg)](https://pkg.go.dev/github.com/cespare/mph)

mph is a Go package for that implements a [minimal perfect hash table][mph] over
strings. It uses the ["Hash, displace, and compress" algorithm][algo]  and the
[Murmur3 hash function][murmur3].

Some quick benchmark results (this is on an i7-8700K):

* `Build` constructs a minimal perfect hash table from a 102k word dictionary in
  18ms (construction time is linear in the size of the input).
* `Lookup`s on that dictionary take about 30ns and are 27% faster than a
  `map[string]uint32`:

    ```
    BenchmarkTable-12               199293806               29.99 ns/op
    BenchmarkTableMap-12            145449822               40.92 ns/op
    ```

[mph]: https://en.wikipedia.org/wiki/Perfect_hash_function#Minimal_perfect_hash_function
[algo]: http://cmph.sourceforge.net/papers/esa09.pdf
[murmur3]: https://en.wikipedia.org/wiki/MurmurHash
