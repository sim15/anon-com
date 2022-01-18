package sposs

import (
	"math/rand"
	"testing"
	"time"

	"github.com/sim15/anon-com/algebra"
)

// 2048 bit first type Cunnigham Chain of length 3
// origin prime of the form p = 2q+1
// q = 2k+1
const primeHexP = "88426e468d0e90c43ac3d7ff2713ec3e341b1ff2dbdc0f9ef8e7067e5e95d73ab553ffb19d094cae390bb2f1e0c28c4cbbaf3858f071568b120b10a36c9d058b5a219e5842a8ac8c59c8a787b353322e26ee80275fb0d6b39133d7250b9dbd570ea457ad766539196dd93017ecb117e65590422ac309415931554b0e71d6b96008f216782f082cbddfdb7f79b37ace203da13cfe072df9291501efd0edd280c739a7e01010e8782e78ebc556ce7c2a4b54c338d4ee5cc5e2fb668ba6d0a793ea345559768ea104b1b984118b47ea2e8670f722db9d6cdb0e802b79b0c1daa48160308bda2bba41adcc2b884a31a6274be34e11bda421dde626de94a1dc522d47"
const generatorG = "5"

// for testing:
// const primeHexP = "17"

// for testing:
// const generatorG = "7"

func TestingGroup() *algebra.Group {
	rand.Seed(time.Now().Unix())

	p := FromSafeHex(primeHexP)
	g := FromSafeHex(generatorG)

	// Initialize field values
	baseField := algebra.NewField(p)
	group := algebra.NewGroup(baseField, baseField.NewElement(g))

	return group
}

func TestFullSPoSS(t *testing.T) {
	group := TestingGroup()
	pp := NewPublicParams(group)

	for i := 0; i < 100; i++ {
		r := group.Field.RandomElement()
		pp.SetRandSeed(r)

		x := pp.ExpField.RandomElement()

		// generate additive shares of g^x
		gX := pp.Group.NewElement(x.Int).Value
		additiveShareA, additiveShareB := pp.LinearShares(gX)

		// client proof of knowledge
		proofA, proofB := pp.GenProof(x)

		// step 1: each server computes a proof audit share that is sent to the other server
		pubAuditShareA, privAuditShareA := pp.PrepareAudit(proofA, additiveShareA, false)
		pubAuditShareB, privAuditShareB := pp.PrepareAudit(proofB, additiveShareB, true)

		// step 2: each server uses the received audit share to update the private and public audits
		pubVerificationShareA, privVerificationShareA := pp.Audit(pubAuditShareB, privAuditShareA, false)
		pubVerificationShareB, privVerificationShareB := pp.Audit(pubAuditShareA, privAuditShareB, true)

		// step 3: check that all the values are correct (i.e., the client didn't provide a bad proof)
		okA := pp.VerifyAudit(pubVerificationShareB, privVerificationShareA)
		okB := pp.VerifyAudit(pubVerificationShareA, privVerificationShareB)

		if !okA || !okB {
			t.Fatalf("SPoSS audit and verification test failed")
		}

	}
}

func BenchmarkClientProof(b *testing.B) {
	group := TestingGroup()
	pp := NewPublicParams(group)

	x := pp.Group.Field.RandomElement()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pp.GenProof(x)
	}

}

func BenchmarkServer(b *testing.B) {
	group := TestingGroup()
	pp := NewPublicParams(group)

	r := group.Field.RandomElement()
	pp.SetRandSeed(r)

	x := pp.ExpField.RandomElement()

	// generate additive shares of g^x
	gX := pp.Group.NewElement(x.Int).Value
	additiveShareA, _ := pp.LinearShares(gX)

	// client proof of knowledge
	proofA, _ := pp.GenProof(x)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// step 1: each server computes a proof audit share that is sent to the other server
		pubAuditShareA, privAuditShareA := pp.PrepareAudit(proofA, additiveShareA, false)
		pubAuditShareB, _ := pubAuditShareA, privAuditShareA

		// step 2: each server uses the received audit share to update the private and public audits
		pubVerificationShareA, privVerificationShareA := pp.Audit(pubAuditShareB, privAuditShareA, false)
		pubVerificationShareB, _ := pubVerificationShareA, privVerificationShareA

		// step 3: check that all the values are correct (i.e., the client didn't provide a bad proof)
		pp.VerifyAudit(pubVerificationShareB, privVerificationShareA)

	}
}
