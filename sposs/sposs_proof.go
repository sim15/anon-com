package sposs

import (
	"github.com/sim15/anon-com/algebra"
)

func (pp *PublicParams) GenProof(x *algebra.FieldElement) (*ProofShare, *ProofShare) {
	xA, xB := pp.ExpLinearShares(x)
	a, b := pp.Group.Field.RandomElement(), pp.Group.Field.RandomElement()
	cA, cB := pp.LinearShares(pp.Group.Field.Mul(a, b))
	return &ProofShare{xA, a, cA}, &ProofShare{xB, b, cB}
}
