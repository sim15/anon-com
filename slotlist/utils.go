package slotlist

import (
	"math/rand"
)

func PRG(seed, size int64) []byte {
	rand.Seed(seed)
	res := make([]byte, size)
	rand.Read(res)
	return res
}
