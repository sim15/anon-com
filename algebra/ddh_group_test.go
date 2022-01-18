package algebra

import (
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func setupDDHGroup(p *big.Int, q *big.Int, n int) (*DDHGroup, []*GroupElement) {
	rand.Seed(time.Now().UnixNano())

	field := NewField(p)
	expField := NewField(q)
	g := findRandomQuadraticGenerator(field)
	h := findRandomQuadraticGenerator(field)
	group := NewGroup(field, g)
	expGroup := NewGroup(expField, h)

	ddhGroup := NewDDHGroup(group, expGroup)

	elements := make([]*GroupElement, n)
	for i := 0; i < n; i++ {
		val := big.NewInt(rand.Int63())
		sign := rand.Int() % 2
		if sign == 0 {
			val.Sub(big.NewInt(0), val)
		}
		elements[i] = group.NewElement(val)
	}
	return ddhGroup, elements
}
func TestDDHGroup(t *testing.T) {

	n := 100
	p := big.NewInt(1523) // 1523 is a safe prime
	q := big.NewInt(761)  // 1523 = 2*761 + 1
	setupDDHGroup(p, q, n)
}
