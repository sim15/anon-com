package slotlist

import (

	// "crypto/rand"

	"github.com/sim15/anon-com/sposs"
)

type SlotList struct {
	ProofParams *sposs.PublicParams
	SlotSize    uint64
	Slots       []*Slot
}

func NewSlotList(pp *sposs.PublicParams, size, numSlots uint64) *SlotList {
	return &SlotList{
		pp,
		size,
		make([]*Slot, numSlots)}
}
