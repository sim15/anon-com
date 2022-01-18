package construction1

import (
	"github.com/sim15/anon-com/algebra"
	"github.com/sim15/anon-com/dpf"
	"github.com/sim15/anon-com/slotlist"
	"github.com/sim15/anon-com/sposs"
)

type Server struct {
	NumBoxes       uint64
	Boxes          *slotlist.SlotList
	SPOSSKeys      []*algebra.GroupElement
	ServerID       bool
	CurrentSession *Session
}

type QueryResult struct {
	Share *algebra.FieldElement
	Pi    []byte
}

type Session struct {
	RecievedQuery *ClientQuery
	PfShare       []int64
	Pfbits        []bool
	QueryShare    *QueryResult
}

func NewServer(id bool, boxsize, numboxes uint64, pp *sposs.PublicParams) *Server {

	return &Server{
		NumBoxes:  numboxes,
		ServerID:  id,
		Boxes:     slotlist.NewSlotList(pp, boxsize, numboxes),
		SPOSSKeys: make([]*algebra.GroupElement, numboxes)}
}

func (s *Server) InitTestSlotList(testgxval *algebra.GroupElement) {
	for i := 0; i < int(s.NumBoxes); i += 1 {
		s.Boxes.Slots[i] = s.Boxes.NewSlot(uint64(i), testgxval)
		s.SPOSSKeys[i] = s.Boxes.Slots[i].SPOSSKey
	}
}

func (s *Server) StartSession(recieved *ClientQuery) {
	s.CurrentSession = &Session{
		RecievedQuery: recieved}
	s.CurrentSession.QueryShare = &QueryResult{
		Share: s.Boxes.ProofParams.Group.Field.AddIdentity(),
	}

}

func (s *Server) expandVDPF(userQuery *ClientQuery) ([]int64, []bool, []byte) {

	pf := dpf.ServerVDPFInitialize(userQuery.PrfKey, userQuery.HKey1, userQuery.HKey2)

	// var indices []uint64

	// for i := uint64(0); i < s.NumBoxes; i++ {
	// 	indices = append(indices, i)
	// }

	return pf.FullDomainVerEval(userQuery.DPFKey)
	// consider redoing dpf code to convert into bitlist on first loop (so you dont need to first expand fully and then into bits)
}

func (s *Server) ComputePrepareAuthAudit() {
	// expand bits and vdpf into slotlist
	s.CurrentSession.PfShare, s.CurrentSession.Pfbits, s.CurrentSession.QueryShare.Pi = s.expandVDPF(s.CurrentSession.RecievedQuery)

	for i := uint64(0); i < s.NumBoxes; i++ {
		if s.CurrentSession.Pfbits[i] {
			s.Boxes.ProofParams.Group.Field.AddInplace(s.CurrentSession.QueryShare.Share, s.SPOSSKeys[i].Value)
		}
	}

}
