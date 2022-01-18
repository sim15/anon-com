package construction1

import (
	"math/big"
	"math/rand"
	"time"

	"github.com/sim15/anon-com/algebra"
)

const primeHexP = "88426e468d0e90c43ac3d7ff2713ec3e341b1ff2dbdc0f9ef8e7067e5e95d73ab553ffb19d094cae390bb2f1e0c28c4cbbaf3858f071568b120b10a36c9d058b5a219e5842a8ac8c59c8a787b353322e26ee80275fb0d6b39133d7250b9dbd570ea457ad766539196dd93017ecb117e65590422ac309415931554b0e71d6b96008f216782f082cbddfdb7f79b37ace203da13cfe072df9291501efd0edd280c739a7e01010e8782e78ebc556ce7c2a4b54c338d4ee5cc5e2fb668ba6d0a793ea345559768ea104b1b984118b47ea2e8670f722db9d6cdb0e802b79b0c1daa48160308bda2bba41adcc2b884a31a6274be34e11bda421dde626de94a1dc522d47"
const primeHexQ = "44213723468748621d61ebff9389f61f1a0d8ff96dee07cf7c73833f2f4aeb9d5aa9ffd8ce84a6571c85d978f06146265dd79c2c7838ab4589058851b64e82c5ad10cf2c215456462ce453c3d9a9991713774013afd86b59c899eb9285cedeab87522bd6bb329c8cb6ec980bf6588bf32ac821156184a0ac98aaa58738eb5cb004790b3c1784165eefedbfbcd9bd67101ed09e7f0396fc948a80f7e876e940639cd3f00808743c173c75e2ab673e1525aa619c6a772e62f17db345d36853c9f51a2aacbb47508258dcc208c5a3f51743387b916dceb66d874015bcd860ed5240b01845ed15dd20d6e615c42518d313a5f1a708ded210eef3136f4a50ee2916a3"
const generatorG = "5"

// TESTING VALS
// const primeHexP = "17"
// const primeHexQ = "B"
// const generatorG = "7"

func DefaultSetup() (*algebra.Group, *algebra.FieldElement) {
	rand.Seed(time.Now().Unix())

	p := FromSafeHex(primeHexP)
	q := FromSafeHex(primeHexQ)
	g := FromSafeHex(generatorG)

	if big.NewInt(0).Exp(g, q, p).Cmp(big.NewInt(1)) == 0 {
		panic("g isn't a generator of order 2q")
	}

	baseField := algebra.NewField(p)

	group := algebra.NewGroup(baseField, baseField.NewElement(g))

	return group, baseField.NewElement(q)
}
