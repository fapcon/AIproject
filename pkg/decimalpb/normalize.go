package decimalpb

import (
	"math"
	"math/big"

	"github.com/shopspring/decimal"
)

const tenNum = 10

var /* const */ ten = big.NewInt(tenNum) //nolint:gochecknoglobals // const
var /* const */ zero = decimal.New(0, 0) //nolint:gochecknoglobals // const

// Normalize returns [decimal.Decimal] equal to input d, but rescaled.
// [decimal.Decimal.Exponent] is changed such a way to remove all trailing zeroes
// from decimal representation of [decimal.Decimal.Coefficient].
// As a special case zero decimal has zero [decimal.Decimal.Exponent].
//
// Examples
//
//	Normalize(decimal.New(100, 0)) -> decimal.New(1, 2)
//	Normalize(decimal.New(100, -3)) -> decimal.New(1, -1)
//	Normalize(decimal.New(123, -2)) -> decimal.New(123, -2)
//	Normalize(decimal.New(0, 1)) -> decimal.New(0, 0)
func Normalize(d decimal.Decimal) decimal.Decimal {
	if d.IsZero() {
		return zero
	}

	value := d.Coefficient()
	remainder := new(big.Int)
	mayValue := new(big.Int)
	n := int64(0)
	maxN := math.MaxInt32 - int64(d.Exponent())
	for {
		if n == maxN {
			break // avoid int32 overflow
		}

		mayValue.QuoRem(value, ten, remainder)
		if remainder.Sign() != 0 {
			break
		}

		value.Set(mayValue)
		n++
	}

	if n == 0 {
		return d
	}

	return decimal.NewFromBigInt(value, d.Exponent()+int32(n))
}
