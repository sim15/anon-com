package algebra

import (
	"math/big"
)

type Field struct {
	P *big.Int // field modulus
}

type FieldElement struct {
	Int *big.Int
}

// new field of order p
func NewField(p *big.Int) *Field {
	return &Field{P: p}
}

// add modulo P
func (f *Field) Add(a, b *FieldElement) *FieldElement {
	newValue := new(big.Int).Mod(new(big.Int).Add(a.Int, b.Int), f.P)
	return f.NewElement(newValue)
}

func (f *Field) AddInplace(a *FieldElement, b *FieldElement) {
	a.Int.Add(a.Int, b.Int).Mod(a.Int, f.P)
}

// sub modulo P
func (f *Field) Sub(a, b *FieldElement) *FieldElement {
	newValue := new(big.Int).Mod(new(big.Int).Sub(a.Int, b.Int), f.P)
	return f.NewElement(newValue)
}

func (f *Field) SubInplace(a *FieldElement, b *FieldElement) {
	a.Int.Sub(a.Int, b.Int).Mod(a.Int, f.P)
}

// return multiplicative inverse with mod P
func (f *Field) MulInv(a *FieldElement) *FieldElement {
	newValue := new(big.Int).ModInverse(a.Int, f.P)
	return f.NewElement(newValue)
}

// multiply mod P
func (f *Field) Mul(a, b *FieldElement) *FieldElement {
	newValue := new(big.Int).Mul(a.Int, b.Int)
	return f.NewElement(newValue)
}

func (f *Field) MulInplace(a *FieldElement, b *FieldElement) {
	a.Int.Mul(a.Int, b.Int).Mod(a.Int, f.P)
}

// exponentiation mod P
func (f *Field) Exp(a *FieldElement, c *big.Int) *FieldElement {
	newValue := gmpExp(a.Int, c, f.P)
	return f.NewElement(newValue)
}

func (f *Field) ExpInplace(a *FieldElement, c *big.Int) {
	gmpExpInplace(a.Int, c, f.P)
}

// new element (a mod P)
func (f *Field) NewElement(a *big.Int) *FieldElement {
	newValue := new(big.Int).Mod(a, f.P)
	return &FieldElement{newValue}
}

// returns a random element in the field
func (f *Field) RandomElement() *FieldElement {
	a := randomInt(f.P)
	return f.NewElement(a)
}

func (f *Field) AddIdentity() *FieldElement {
	return &FieldElement{big.NewInt(0)}
}

func (f *Field) MulIdentity() *FieldElement {
	return &FieldElement{big.NewInt(1)}
}

func (f *Field) IsAddIdentity(e *FieldElement) bool {
	return f.AddIdentity().Cmp(e) == 0
}

func (f *Field) IsMulIdentity(e *FieldElement) bool {
	return f.MulIdentity().Cmp(e) == 0
}

func (f *Field) IsZero(e *FieldElement) bool {
	return f.AddIdentity().Cmp(e) == 0
}

func (elem *FieldElement) Cmp(b *FieldElement) int {
	return elem.Int.Cmp(b.Int)
}
