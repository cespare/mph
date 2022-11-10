// Package mph implements a minimal perfect hash table over strings.
package mph

import (
	"github.com/zeebo/xxh3"
	"golang.org/x/exp/slices"
)

// A Table is an immutable hash table that provides constant-time lookups of key
// indices using a minimal perfect hash.
type Table struct {
	keys       []string
	level0     []uint32 // power of 2 size
	level0Mask int      // len(Level0) - 1
	level1     []uint32 // power of 2 size >= len(keys)
	level1Mask int      // len(Level1) - 1
}

// Build builds a Table from keys using the "Hash, displace, and compress"
// algorithm described in http://cmph.sourceforge.net/papers/esa09.pdf.
func Build(keys []string) *Table {
	var (
		level0        = make([]uint32, nextPow2(len(keys)/4))
		level0Mask    = len(level0) - 1
		level1        = make([]uint32, nextPow2(len(keys)))
		level1Mask    = len(level1) - 1
		sparseBuckets = make([][]int, len(level0))
	)
	for i, s := range keys {
		n := int(xxh3.HashStringSeed(s, 0)) & level0Mask
		sparseBuckets[n] = append(sparseBuckets[n], i)
	}
	var buckets []indexBucket
	for n, vals := range sparseBuckets {
		if len(vals) > 0 {
			buckets = append(buckets, indexBucket{n, vals})
		}
	}
	slices.SortFunc(buckets, func(a, b indexBucket) bool {
		return len(a.vals) > len(b.vals)
	})

	occ := make([]bool, len(level1))
	var tmpOcc []int
	for _, bucket := range buckets {
		var seed uint32
	trySeed:
		tmpOcc = tmpOcc[:0]
		for _, i := range bucket.vals {
			n := int(xxh3.HashStringSeed(keys[i], uint64(seed))) & level1Mask
			if occ[n] {
				for _, n := range tmpOcc {
					occ[n] = false
				}
				seed++
				goto trySeed
			}
			occ[n] = true
			tmpOcc = append(tmpOcc, n)
			level1[n] = uint32(i)
		}
		level0[int(bucket.n)] = uint32(seed)
	}

	return &Table{
		keys:       keys,
		level0:     level0,
		level0Mask: level0Mask,
		level1:     level1,
		level1Mask: level1Mask,
	}
}

func nextPow2(n int) int {
	for i := 1; ; i *= 2 {
		if i >= n {
			return i
		}
	}
}

// Lookup searches for s in t and returns its index and whether it was found.
func (t *Table) Lookup(s string) (n uint32, ok bool) {
	i0 := int(xxh3.HashStringSeed(s, 0)) & t.level0Mask
	seed := t.level0[i0]
	i1 := int(xxh3.HashStringSeed(s, uint64(seed))) & t.level1Mask
	n = t.level1[i1]
	return n, s == t.keys[int(n)]
}

type indexBucket struct {
	n    int
	vals []int
}
