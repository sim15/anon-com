package slotlist

import (

	// "crypto/rand"

	"github.com/sim15/anon-com/sposs"
)

type SlotList struct {
	ProofParams *sposs.PublicParams
	SlotSize    uint64
	// Slots       []*Slot
}

type SpossSlotList struct {
	*SlotList
	Slots []*SpossSlot
}

func NewSpossSlotList(pp *sposs.PublicParams, size, numSlots uint64) *SpossSlotList {
	return &SpossSlotList{
		&SlotList{
			pp,
			size},
		make([]*SpossSlot, numSlots)}
}
