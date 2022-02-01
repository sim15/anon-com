package slotlist

import (
	"math/rand"
	"testing"
	"time"
)

const PRGsize = 10

func BenchmarkPRG(b *testing.B) {
	s1 := rand.NewSource(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		PRG((rand.New(s1).Int63n(100)), PRGsize)
	}
}
