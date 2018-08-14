package mph

import (
	"fmt"
	"strings"
	"testing"
)

var murmurTestCases = []struct {
	input string
	seed  murmurSeed
	want  uint32
}{
	{"", 0, 0},
	{"", 1, 0x514e28b7},
	{"", 0xffffffff, 0x81f16f39},
	{"\xff\xff\xff\xff", 0, 0x76293b50},
	{"!Ce\x87", 0, 0xf55b516b},
	{"!Ce\x87", 0x5082edee, 0x2362f9de},
	{"!Ce", 0, 0x7e4a8634},
	{"!C", 0, 0xa0f7b07a},
	{"!", 0, 0x72661cf4},
	{"\x00\x00\x00\x00", 0, 0x2362f9de},
	{"\x00\x00\x00", 0, 0x85f0b427},
	{"\x00\x00", 0, 0x30f4c306},
	{"Hello, world!", 0x9747b28c, 0x24884CBA},
	{"ππππππππ", 0x9747b28c, 0xD58063C1},
	{"abc", 0, 0xb3dd93fa},
	{"abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq", 0, 0xee925b90},
	{"The quick brown fox jumps over the lazy dog", 0x9747b28c, 0x2fa826cd},
	{strings.Repeat("a", 256), 0x9747b28c, 0x37405bdc},
}

func TestMurmur(t *testing.T) {
	for _, tt := range murmurTestCases {
		got := tt.seed.hash(tt.input)
		if got != tt.want {
			t.Errorf("hash(%q, seed=0x%x): got 0x%x; want %x",
				tt.input, tt.seed, got, tt.want)
		}
	}
}

func BenchmarkMurmur(b *testing.B) {
	for _, size := range []int{1, 4, 8, 16, 32, 50, 500} {
		b.Run(fmt.Sprint(size), func(b *testing.B) {
			s := strings.Repeat("a", size)
			b.SetBytes(int64(size))
			var seed murmurSeed
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				seed.hash(s)
			}
		})
	}
}
