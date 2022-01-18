package sposs

import "github.com/sim15/anon-com/algebra"

type PublicParams struct {
	Group    *algebra.Group
	ExpField *algebra.Field
	RandSeed *algebra.FieldElement
}

type ProofShare struct {
	AdditiveShareX *algebra.FieldElement // additive share of x such that g^x =
	BeaverConstant *algebra.FieldElement // a or b
	AdditiveShareC *algebra.FieldElement // [c]
}

type PublicAuditShare struct {
	MulConstant *algebra.FieldElement // d or e
}

type PrivateAuditShare struct {
	AdditiveShareY      *algebra.FieldElement
	MulShareY           *algebra.FieldElement
	KnownBeaverConstant *algebra.FieldElement
	AdditiveShareC      *algebra.FieldElement
	KnownMulConstant    *algebra.FieldElement
	UnknownMulConstant  *algebra.FieldElement
}

type PublicVerificationShare struct {
	BeaverOutputShare *algebra.FieldElement
}

type PrivateVerificationShare struct {
	KnownBeaverOutputShare   *algebra.FieldElement
	UnknownBeaverOutputShare *algebra.FieldElement
}

func NewPublicParams(g *algebra.Group) *PublicParams {
	f := algebra.NewField(g.Field.Pminus1()) // TODO: [priority minor] p-1 is not prime; shouldn't call this a field
	return &PublicParams{g, f, nil}
}

func (pp *PublicParams) SetRandSeed(seed *algebra.FieldElement) {
	pp.RandSeed = seed
}
