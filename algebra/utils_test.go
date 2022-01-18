package algebra

import (
	"fmt"
	"math/big"
	"testing"
)

func TestPrimeFactors(t *testing.T) {
	if fmt.Sprintf("%v", PrimeFactors(big.NewInt(23))) != `[23]` {
		t.Error(23)
	}
	if fmt.Sprintf("%v", PrimeFactors(big.NewInt(12))) != `[2 2 3]` {
		t.Error(12)
	}
	if fmt.Sprintf("%v", PrimeFactors(big.NewInt(360))) != `[2 2 2 3 3 5]` {
		t.Error(PrimeFactors(big.NewInt(360)))
	}
	if fmt.Sprintf("%v", PrimeFactors(big.NewInt(97))) != `[97]` {
		t.Error(97)
	}
}
