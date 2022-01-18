package construction1

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/sim15/anon-com/dpf"
	"github.com/sim15/anon-com/slotlist"
	"github.com/sim15/anon-com/sposs"
)

const mSize = 5
const NumMailboxes = 512
const NumQueries = 100

func TestMessageEncode(t *testing.T) {
	group, _ := DefaultSetup()

	pp := sposs.NewPublicParams(group)

	x := big.NewInt(4)

	// sA := NewServer(false, mSize, NumMailboxes, pp)
	// sB := NewServer(true, mSize, NumMailboxes, pp)

	// sA.InitTestSlotList(x)
	// sB.InitTestSlotList(pp.ExpField.Add(pp.ExpField.NewElement(x), query).Int)

	for i := 0; i < NumQueries; i++ {
		// x := big.NewInt(4)
		message := []byte{101, 100}
		c := NewClient(pp, x, mSize, message)

		query := c.NewClientQuery(0, NumMailboxes, x)

		pf1 := dpf.ServerVDPFInitialize(query[0].PrfKey, query[0].HKey1, query[0].HKey2)
		pf2 := dpf.ServerVDPFInitialize(query[1].PrfKey, query[1].HKey1, query[1].HKey2)

		resA, _, _ := pf1.BatchVerEval(query[0].DPFKey, []uint64{0})
		resB, _, _ := pf2.BatchVerEval(query[1].DPFKey, []uint64{0})
		// fmt.Println("---")
		// fmt.Println(resA)
		// fmt.Println(resB)
		// fmt.Println(resA[0] - resB[0])

		// fmt.Println(slotlist.PRG(resA[0], c.MessageSize))
		// fmt.Println(slotlist.PRG(resB[0], c.MessageSize))

		mMask := slotlist.XorByteArray(slotlist.PRG(resA[0], c.MessageSize), slotlist.PRG(resB[0], c.MessageSize))

		unmaskedMessage := slotlist.XorByteArray(mMask, query[0].MaskedM)

		// fmt.Println(unmaskedMessage)

		if bytes.Compare(unmaskedMessage, message) == 1 {
			t.Fatalf(
				"Message shares do not recover message correctly",
			)
		}

	}

}

func TestFullClientAuth(t *testing.T) {
	group, q := DefaultSetup()

	pp := sposs.NewPublicParams(group)

	serverTestX := big.NewInt(4)
	serverTestAltX := pp.ExpField.Add(pp.ExpField.NewElement(serverTestX), q).Int

	sA := NewServer(false, mSize, NumMailboxes, pp)
	sB := NewServer(true, mSize, NumMailboxes, pp)

	sA.InitTestSlotList(group.NewElement(serverTestX))
	sB.InitTestSlotList(group.NewElement(serverTestAltX))

	for i := 0; i < NumQueries; i++ {

		x := pp.ExpField.RandomElement().Int
		altX := pp.ExpField.Add(pp.ExpField.NewElement(x), q).Int

		sA.Boxes.Slots[0].SPOSSKey = pp.Group.NewElement(x)
		sB.Boxes.Slots[0].SPOSSKey = pp.Group.NewElement(altX)

		sA.SPOSSKeys[0] = pp.Group.NewElement(x)
		sB.SPOSSKeys[0] = pp.Group.NewElement(altX)

		message := []byte{101, 100}
		c := NewClient(pp, x, mSize, message)

		query := c.NewClientQuery(0, NumMailboxes, x)

		sA.StartSession(query[0])
		sB.StartSession(query[1])

		sA.ComputePrepareAuthAudit()
		sB.ComputePrepareAuthAudit()

		// fmt.Printf("%v\n", sA.SPOSSKeys[0].Value.Int)
		// fmt.Printf("%v		%v\n", sA.CurrentSession.Pfbits, sA.CurrentSession.PfShare)
		// fmt.Printf("%v		%v\n", sB.CurrentSession.Pfbits, sB.CurrentSession.PfShare)
		// fmt.Println(sA.CurrentSession.QueryShare.Share)
		// fmt.Println(sB.CurrentSession.QueryShare.Share)

		if !bytes.Equal(sA.CurrentSession.QueryShare.Pi, sB.CurrentSession.QueryShare.Pi) {
			t.Fatalf("pi0 =/= p1: Got: %v and %v", sA.CurrentSession.QueryShare.Pi, sB.CurrentSession.QueryShare.Pi)
		}

		recoveredSPOSSkey := group.Field.Add(sA.CurrentSession.QueryShare.Share, sB.CurrentSession.QueryShare.Share)

		if recoveredSPOSSkey.Cmp(group.NewElement(x).Value) == recoveredSPOSSkey.Cmp(group.NewElement(altX).Value) {
			t.Fatalf("Recovered g^x (sposs key) does not match user's sposs key \n Got: %v \n Expected: %v or %v \n",
				recoveredSPOSSkey.Int,
				group.NewElement(x).Value.Int,
				group.NewElement(altX).Value.Int,
			)
		}

		recievedShare := group.NewElement(sA.CurrentSession.RecievedQuery.SPoSSProof.AdditiveShareX.Int)
		recievedShare = group.Mul(recievedShare, group.NewElement(sB.CurrentSession.RecievedQuery.SPoSSProof.AdditiveShareX.Int))

		if recievedShare.Value.Cmp(recoveredSPOSSkey) != 0 {
			t.Fatalf("Recovered g^x (sposs key) does not match user's SHARED proof value \n Recovered: %v \n Recieved: %v \n",
				recoveredSPOSSkey.Int,
				recievedShare.Value.Int,
			)
		}

		rand := sA.Boxes.ProofParams.Group.Field.RandomElement()
		sA.Boxes.ProofParams.SetRandSeed(rand)
		sB.Boxes.ProofParams.SetRandSeed(rand)

		// fmt.Println(c.ProofPP.Group.NewElement(x).Value.Int)
		// fmt.Println(recoveredSPOSSkey.Int)
		// fmt.Println(c.ProofPP.RandSeed)
		// fmt.Println(sB.Boxes.ProofParams.RandSeed)
		// fmt.Println(sA.CurrentSession.QueryShare.Share)
		// fmt.Println(sB.CurrentSession.QueryShare.Share)
		// fmt.Println(sA.CurrentSession.RecievedQuery.SPoSSProof)
		// fmt.Println(sB.CurrentSession.RecievedQuery.SPoSSProof)

		// setp 1 of sposs
		pubAuditShareA, privAuditShareA := sA.Boxes.ProofParams.PrepareAudit(
			sA.CurrentSession.RecievedQuery.SPoSSProof,
			sA.CurrentSession.QueryShare.Share,
			false)
		pubAuditShareB, privAuditShareB := sB.Boxes.ProofParams.PrepareAudit(
			sB.CurrentSession.RecievedQuery.SPoSSProof,
			sB.CurrentSession.QueryShare.Share,
			true)

		// step 2 of sposs
		pubVerificationShareA, privVerificationShareA := sA.Boxes.ProofParams.Audit(
			pubAuditShareB,
			privAuditShareA,
			false)

		pubVerificationShareB, privVerificationShareB := sB.Boxes.ProofParams.Audit(
			pubAuditShareA,
			privAuditShareB,
			true)

		// step 3 of sposs
		okA := sA.Boxes.ProofParams.VerifyAudit(pubVerificationShareB, privVerificationShareA)
		okB := sB.Boxes.ProofParams.VerifyAudit(pubVerificationShareA, privVerificationShareB)

		if !okA || !okB {
			t.Fatalf("Auth failed on query %v (okA = %v | okB = %v)", i, okA, okB)
		}

	}

}
