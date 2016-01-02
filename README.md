# mph

[![GoDoc](https://godoc.org/github.com/cespare/mph?status.svg)](https://godoc.org/github.com/cespare/mph)

mph is a Go package for that implements a
[minimal perfect hash table](https://en.wikipedia.org/wiki/Perfect_hash_function#Minimal_perfect_hash_function)
over strings. It uses the
["Hash, displace, and compress" algorithm](http://cmph.sourceforge.net/papers/esa09.pdf) and the
[Murmur3 hash function](https://en.wikipedia.org/wiki/MurmurHash).

Very rough benchmarks from my laptop:

* `Build` constructed a minimal perfect hash table from a 236k word dictionary in <200ms (construction time is
  linear in the size of the input).
* `Lookup`s on that dictionary took ~85ns, almost twice as fast as lookups in a `map[string]uint32`.
