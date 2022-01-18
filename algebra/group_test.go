package algebra

import (
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func setupGroup(p *big.Int, n int) (*Group, []*GroupElement) {
	rand.Seed(time.Now().UnixNano())

	field := NewField(p)
	g := findRandomGenerator(field)
	group := NewGroup(field, g)
	elements := make([]*GroupElement, n)
	for i := 0; i < n; i++ {
		val := big.NewInt(rand.Int63())
		sign := rand.Int() % 2
		if sign == 0 {
			val.Sub(big.NewInt(0), val)
		}
		elements[i] = group.NewElement(val)
	}
	return group, elements
}
func TestAddGroup(t *testing.T) {

	n := 100
	p := big.NewInt(1523) // 1523 is a safe prime
	group, elements := setupGroup(p, n)

	sumInt := group.Field.MulIdentity()
	sum := group.Identity()
	for i := 0; i < n; i++ {
		sum = group.Mul(sum, elements[i])
		sumInt = group.Field.Mul(sumInt, elements[i].Value)
	}

	if sum.Value.Int.Cmp(big.NewInt(0)) == 0 {
		t.Fatalf("Group element should never be zero!")
	}

	if sum.Value.Int.Cmp(sumInt.Int) != 0 {
		t.Fatalf("Sum over group is not correct. expected: %v, got: %v", sumInt, sum.Value.Int)
	}
}

func TestInverseGroup(t *testing.T) {

	n := 100
	p := big.NewInt(1523) // 1523 is a safe prime
	group, elements := setupGroup(p, n)

	sum := group.Identity()
	sumInv := group.Identity()
	for i := 0; i < n; i++ {
		sum = group.Mul(sum, elements[i])
		sumInv = group.Mul(sumInv, group.MulInv(elements[i]))
	}

	test := group.Mul(sum, sumInv)

	if test.Value.Int.Cmp(big.NewInt(0)) == 0 {
		t.Fatalf("Group element should never be zero!")
	}

	if test.Cmp(group.Identity()) != 0 {
		t.Fatalf("Multiplicative inverse is incorrect!")
	}
}
