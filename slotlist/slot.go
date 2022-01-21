package slotlist

import (
	"math/big"

	"github.com/sim15/anon-com/algebra"
)

type Slot struct {
	PhysicalAddress uint64
	// VirtualAddress  *big.int
	// SPOSSKey  *algebra.GroupElement
	DataShare []byte
}

type ExpressSlot struct {
	*Slot
	DecryptionKey *big.Int //temporary
}

type SpossSlot struct {
	*Slot
	SPOSSKey *algebra.GroupElement
}

func (sl *SlotList) NewSpossSlot(phys uint64, gx *algebra.GroupElement) *SpossSlot {
	return &SpossSlot{
		&Slot{
			phys,
			make([]byte, sl.SlotSize)},
		gx.Copy()}
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

func XorToSlot(a *Slot, b []byte) {

	if len(a.DataShare) < len(b) {
		for j := 0; j < len(a.DataShare); j++ {
			a.DataShare[j] ^= b[j]
		}
	} else {
		for j := 0; j < len(b); j++ {
			a.DataShare[j] ^= b[j]
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

func (sl *Slot) WriteMessageToSlot(maskedMessage []byte, slotSize uint64, prgKey int64, bit bool) {

	if bit {
		XorToSlot(sl, PRG(prgKey, slotSize))
	}

	XorToSlot(sl, maskedMessage)
}
