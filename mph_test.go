package mph

import (
	"bufio"
	"os"
	"strconv"
	"sync"
	"testing"
)

func TestBuild_simple(t *testing.T) {
	testTable(t, []string{"foo", "foo2", "bar", "baz"}, []string{"quux"})
}

func TestBuild_stress(t *testing.T) {
	var keys, extra []string
	for i := 0; i < 20000; i++ {
		s := strconv.Itoa(i)
		if i < 10000 {
			keys = append(keys, s)
		} else {
			extra = append(extra, s)
		}
	}
	testTable(t, keys, extra)
}

func testTable(t *testing.T, keys []string, extra []string) {
	table := Build(keys)
	for i, key := range keys {
		n, ok := table.Lookup(key)
		if !ok {
			t.Errorf("Lookup(%s): got !ok; want ok", key)
			continue
		}
		if int(n) != i {
			t.Errorf("Lookup(%s): got n=%d; want %d", n, i)
		}
	}
	for _, key := range extra {
		if _, ok := table.Lookup(key); ok {
			t.Errorf("Lookup(%s): got ok; want !ok", key)
		}
	}
}

var (
	words      []string
	wordsOnce  sync.Once
	benchTable *Table
)

func BenchmarkBuild(b *testing.B) {
	wordsOnce.Do(loadBenchTable)
	if len(words) == 0 {
		b.Skip("unable to load dictionary file")
	}
	for i := 0; i < b.N; i++ {
		Build(words)
	}
}

func BenchmarkTable(b *testing.B) {
	wordsOnce.Do(loadBenchTable)
	if len(words) == 0 {
		b.Skip("unable to load dictionary file")
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := i % len(words)
		n, ok := benchTable.Lookup(words[j])
		if !ok {
			b.Fatal("missing key")
		}
		if n != uint32(j) {
			b.Fatal("bad result index")
		}
	}
}

// For comparison against BenchmarkTable.
func BenchmarkTableMap(b *testing.B) {
	wordsOnce.Do(loadBenchTable)
	if len(words) == 0 {
		b.Skip("unable to load dictionary file")
	}
	m := make(map[string]uint32)
	for i, word := range words {
		m[word] = uint32(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		j := i % len(words)
		n, ok := m[words[j]]
		if !ok {
			b.Fatal("missing key")
		}
		if n != uint32(j) {
			b.Fatal("bad result index")
		}
	}
}

func loadBenchTable() {
	for _, dict := range []string{"/usr/share/dict/words", "/usr/dict/words"} {
		var err error
		words, err = loadDict(dict)
		if err == nil {
			break
		}
	}
	if len(words) > 0 {
		benchTable = Build(words)
	}
}

func loadDict(dict string) ([]string, error) {
	f, err := os.Open(dict)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var words []string
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return words, nil
}
