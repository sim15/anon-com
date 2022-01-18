package construction1

import (
	"errors"
	"fmt"
	"math/big"
)

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
