package algebra

import (
	"crypto/rand"
	"math/big"

	"github.com/ncw/gmp"
)

// Get all prime factors of a given number n
// taken from https://siongui.github.io/2017/05/09/go-find-all-prime-factors-of-integer-number/
func PrimeFactors(n *big.Int) []*big.Int {

	pfs := make([]*big.Int, 0)
	two := big.NewInt(2)
	zero := big.NewInt(0)

	// Get the number of 2s that divide n
	for new(big.Int).Mod(n, two).Cmp(zero) == 0 {
		pfs = append(pfs, two)
		n.Div(n, two)
	}

	// n must be odd at this point. so we can skip one element
	// (note i = i + 2)
	i := big.NewInt(3)
	for {
		// while i divides n, append i and divide n
		for big.NewInt(0).Mod(n, i).Cmp(zero) == 0 {
			pfs = append(pfs, big.NewInt(0).SetBytes(i.Bytes()))
			n = n.Div(n, i)
		}

		test := big.NewInt(0)
		if test.Mul(i, i).Cmp(n) <= 0 {
			i.Add(i, two)
		} else {
			break
		}
	}

	// This condition is to handle the case when n is a prime number
	// greater than 2
	if n.Cmp(two) > 0 {
		pfs = append(pfs, n)
	}

	return pfs
}

func findRandomGenerator(field *Field) *FieldElement {

	found := false
	one := big.NewInt(1)
	factors := PrimeFactors(field.Pminus1())
	g := randomInt(field.P)
	for {

		// test if g is a generator
		for i := 0; i < len(factors); i++ {
			pow := new(big.Int).Div(field.Pminus1(), factors[i])
			if new(big.Int).Exp(g, pow, field.P).Cmp(one) == 0 {
				break
			}
			if i+1 == len(factors) {
				found = true
			}
		}

		if found {
			break
		}

		// try a new candidate
		g = randomInt(field.P)
	}

	return field.NewElement(g)
}

func findRandomQuadraticGenerator(field *Field) *FieldElement {

	found := false
	one := big.NewInt(1)
	factors := PrimeFactors(field.Pminus1())
	g := randomInt(field.P)
	for {

		// test if g is a generator
		for i := 0; i < len(factors); i++ {
			pow := new(big.Int).Div(field.Pminus1(), factors[i])
			if new(big.Int).Exp(g, pow, field.P).Cmp(one) == 0 {
				break
			}
			if i+1 == len(factors) {
				found = true
			}
		}

		if found {
			break
		}

		// try a new candidate
		g = randomInt(field.P)
	}

	g.Mul(g, g) // make it a generator of the quadratic residues

	return field.NewElement(g)
}

func randomInt(max *big.Int) *big.Int {
	randomBig, _ := rand.Int(rand.Reader, new(big.Int).SetBytes(max.Bytes()))
	return new(big.Int).SetBytes(randomBig.Bytes())
}

func gmpExp(a, b, n *big.Int) *big.Int {
	aGmp := new(gmp.Int).SetBytes(a.Bytes())
	bGmp := new(gmp.Int).SetBytes(b.Bytes())
	nGmp := new(gmp.Int).SetBytes(n.Bytes())
	resGmp := new(gmp.Int).Exp(aGmp, bGmp, nGmp)
	return new(big.Int).SetBytes(resGmp.Bytes())
}

func gmpExpInplace(a, b, n *big.Int) {
	aGmp := new(gmp.Int).SetBytes(a.Bytes())
	bGmp := new(gmp.Int).SetBytes(b.Bytes())
	nGmp := new(gmp.Int).SetBytes(n.Bytes())
	resGmp := new(gmp.Int).Exp(aGmp, bGmp, nGmp)
	a.SetBytes(resGmp.Bytes())
}
