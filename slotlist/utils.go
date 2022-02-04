package slotlist

import (
	"math/rand"
)

func PRG(seed int64, size uint64) []byte {
	rand.Seed(seed)
	res := make([]byte, size)
	rand.Read(res)
	return res
}

// add prg data to clients and servers. store its size and initial state
func PRGAES(c)
