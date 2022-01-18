package slotlist

import (
	"math/big"

	"github.com/sim15/anon-com/algebra"
)

type Slot struct {
	PhysicalAddress uint64
	// VirtualAddress  *big.int
	SPOSSKey  *algebra.GroupElement
	DataShare []byte
}

type ExpressSlot struct {
	*Slot
	DecryptionKey *big.Int //temporary
}

func (sl *SlotList) NewSlot(phys uint64, gx *algebra.GroupElement) *Slot {
	return &Slot{
		phys,
		// virtual,
		gx.Copy(),
		make([]byte, sl.SlotSize)}
}

func XorSlots(a, b *Slot) {

	if len(a.DataShare) < len(b.DataShare) {
		for j := 0; j < len(a.DataShare); j++ {
			a.DataShare[j] ^= b.DataShare[j]
		}
	} else {
		for j := 0; j < len(b.DataShare); j++ {
			a.DataShare[j] ^= b.DataShare[j]
		}
	}
}

func XorByteArray(a, b []byte) []byte {

	if len(a) < len(b) {
		res := make([]byte, len(a))
		for j := 0; j < len(a); j++ {
			res[j] = a[j] ^ b[j]
		}
		return res
	} else {
		res := make([]byte, len(b))
		for j := 0; j < len(b); j++ {
			res[j] = a[j] ^ b[j]
		}
		return res
	}
}
