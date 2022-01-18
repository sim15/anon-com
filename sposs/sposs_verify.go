package sposs

import (
	"math/big"

	"github.com/sim15/anon-com/algebra"
)

func (pp *PublicParams) PrepareAudit(pi *ProofShare, additiveShareY *algebra.FieldElement, serverID bool) (*PublicAuditShare, *PrivateAuditShare) {
	mulShareY := pp.Group.NewElement(pi.AdditiveShareX.Int)
	mulConstant := mulShareY.Copy().Value
	if !serverID {
		mulConstant = pp.Group.Field.Mul(pp.RandSeed, mulConstant)
	}
	mulConstant = pp.Group.Field.Sub(mulConstant, pi.BeaverConstant)

	return &PublicAuditShare{mulConstant},
		&PrivateAuditShare{
			additiveShareY,
			mulShareY.Value,
			pi.BeaverConstant,
			pi.AdditiveShareC,
			mulConstant,
			nil}
}

func (pp *PublicParams) Audit(sharedAudit *PublicAuditShare, preparedAudit *PrivateAuditShare, serverID bool) (*PublicVerificationShare, *PrivateVerificationShare) {
	preparedAudit.UnknownMulConstant = sharedAudit.MulConstant

	f := pp.Group.Field.Mul(preparedAudit.KnownMulConstant, preparedAudit.UnknownMulConstant)
	oneHalf := pp.Group.Field.MulInv(pp.Group.Field.NewElement(big.NewInt(2)))
	f = pp.Group.Field.Mul(f, oneHalf)

	vShare := pp.Group.Field.Mul(preparedAudit.UnknownMulConstant, preparedAudit.KnownBeaverConstant)
	vShare = pp.Group.Field.Add(vShare, f)
	vShare = pp.Group.Field.Add(vShare, preparedAudit.AdditiveShareC)

	wShare := pp.Group.Field.Mul(pp.RandSeed, preparedAudit.AdditiveShareY)
	wShare = pp.Group.Field.Sub(vShare, wShare)

	if serverID {
		wShare = pp.Group.Field.Sub(pp.Group.Field.AddIdentity(), wShare)
	}

	return &PublicVerificationShare{wShare}, &PrivateVerificationShare{wShare, nil}
}

func (pp *PublicParams) VerifyAudit(sharedVerification *PublicVerificationShare, privateVerification *PrivateVerificationShare) bool {
	privateVerification.UnknownBeaverOutputShare = sharedVerification.BeaverOutputShare
	return privateVerification.KnownBeaverOutputShare.Cmp(privateVerification.UnknownBeaverOutputShare) == 0
}
