package construction1

import (
	"math"
	"math/big"

	"github.com/sim15/anon-com/dpf"
	"github.com/sim15/anon-com/slotlist"
	"github.com/sim15/anon-com/sposs"
)

type Client struct {
	ProofPP     *sposs.PublicParams
	MessageSize int64
	x           *big.Int
	m           []byte
}

type ClientQuery struct {
	MaskedM     []byte
	DPFKey      *dpf.DPFKey
	PrfKey      dpf.PrfKey
	HKey1       dpf.HashKey
	HKey2       dpf.HashKey
	ShareNumber uint
	SPoSSProof  *sposs.ProofShare
}

func NewClient(pp *sposs.PublicParams, xVal *big.Int, mSize int64, message []byte) *Client {
	return &Client{
		ProofPP:     pp,
		x:           xVal,
		MessageSize: mSize,
		m:           message}
}

func (c *Client) NewClientQuery(idx, numBoxes uint64, x *big.Int) []*ClientQuery {
	pf := dpf.ClientVDPFInitialize()

	vdpfKeyA, vdpfKeyB := pf.GenVDPFKeys(idx, uint(math.Ceil(math.Log2(float64(numBoxes)))))

	xVal := new(big.Int).Set(x)
	resA, _, _ := pf.BatchVerEval(vdpfKeyA, []uint64{idx})
	resB, bitsB, _ := pf.BatchVerEval(vdpfKeyB, []uint64{idx})

	if bitsB[0] {
		q := c.ProofPP.Group.Field.Pminus1()
		q.Div(q, big.NewInt(2))
		xVal = new(big.Int).Add(xVal, q)
		xVal.Mod(xVal, c.ProofPP.ExpField.P)
	}

	authProofA, authProofB := c.ProofPP.GenProof(c.ProofPP.ExpField.NewElement(xVal))

	shares := make([]*ClientQuery, 2)

	mMask := slotlist.XorByteArray(slotlist.PRG(resA[0], c.MessageSize), slotlist.PRG(resB[0], c.MessageSize))
	mMask = slotlist.XorByteArray(mMask, c.m)

	// share for auth server A
	shares[0] = &ClientQuery{}
	shares[0].MaskedM = mMask
	shares[0].ShareNumber = 0
	shares[0].PrfKey = pf.PrfKey
	shares[0].DPFKey = vdpfKeyA
	shares[0].HKey1 = pf.H1Key
	shares[0].HKey2 = pf.H2Key
	shares[0].SPoSSProof = authProofA

	// share for auth server B
	shares[1] = &ClientQuery{}
	shares[1].MaskedM = mMask
	shares[1].ShareNumber = 1
	shares[1].PrfKey = pf.PrfKey
	shares[1].DPFKey = vdpfKeyB
	shares[1].HKey1 = pf.H1Key
	shares[1].HKey2 = pf.H2Key
	shares[1].SPoSSProof = authProofB

	return shares
}
