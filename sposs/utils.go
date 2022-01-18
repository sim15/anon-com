package sposs

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"math/big"

	"github.com/sim15/anon-com/algebra"
)

// Return a pair of linear shares for toShare, s.t. share1 + share2 = toShare
func (pp *PublicParams) LinearShares(toShare *algebra.FieldElement) (*algebra.FieldElement, *algebra.FieldElement) {
	share1 := pp.Group.Field.RandomElement()
	share2 := pp.Group.Field.Sub(toShare, share1)
	return share1, share2
}

// Return a pair of linear shares for toShare, s.t. share1 + share2 = toShare
// the field is the *exponent field* of the group
func (pp *PublicParams) ExpLinearShares(toShare *algebra.FieldElement) (*algebra.FieldElement, *algebra.FieldElement) {
	share1 := pp.ExpField.RandomElement()
	share2 := pp.ExpField.Sub(toShare, share1)
	return share1, share2
}

func RandomInt(max *big.Int) *big.Int {
	randomBig, _ := rand.Int(rand.Reader, new(big.Int).SetBytes(max.Bytes()))
	return new(big.Int).SetBytes(randomBig.Bytes())
}

func (pp *PublicParams) EvalPoly(poly []*algebra.FieldElement, x *algebra.FieldElement) *algebra.FieldElement {
	res := pp.Group.Field.NewElement(big.NewInt(0))
	for i := 0; i < len(poly); i++ {
		exponent := big.NewInt((int64)(len(poly) - i - 1))
		expVal := pp.Group.Field.Exp(x, exponent)
		expVal = pp.Group.Field.Mul(poly[i], expVal)

		res = pp.Group.Field.Add(res, expVal)
	}
	return res
}

func concatenateNumbers(numbers ...*big.Int) []byte {
	var bs []byte
	for _, n := range numbers {
		bs = append(bs, n.Bytes()...)
	}
	return bs
}

func hashIntoBytesSHA256(numbers ...*big.Int) []byte {
	toBeHashed := concatenateNumbers(numbers...)
	sha256 := sha256.New()
	sha256.Write(toBeHashed)
	hashBytes := sha256.Sum(nil)
	return hashBytes
}

func (pp *PublicParams) CalculateHashShare(xShare, gxShare *algebra.FieldElement) (hashVal *algebra.FieldElement) {
	bytesHash := hashIntoBytesSHA256(xShare.Int, gxShare.Int)
	return pp.ExpField.NewElement(new(big.Int).SetBytes(bytesHash))
}

// from https://github.com/didiercrunch/elgamal/blob/master/elgamal_test.go
func fromHex(hex string) (*big.Int, error) {
	n, err := new(big.Int).SetString(hex, 16)
	if !err {
		msg := fmt.Sprintf("Cannot convert %s to int as hexadecimal", hex)
		return nil, errors.New(msg)
	}
	return n, nil
}

// from https://github.com/didiercrunch/elgamal/blob/master/elgamal_test.go
func FromSafeHex(s string) *big.Int {
	ret, err := fromHex(s)
	if err != nil {
		panic(err)
	}
	return ret
}
