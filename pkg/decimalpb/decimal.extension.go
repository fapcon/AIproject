package decimalpb

import (
	"math"
	"math/big"

	"github.com/shopspring/decimal"
)

// AsDecimal returns normalized [decimal.Decimal].
func (x *Decimal) AsDecimal() decimal.Decimal {
	if x == nil {
		return zero
	}
	return Normalize(decimal.New(x.Unscaled, -x.Scale))
}

// Zero returns normalized zero [Decimal].
func Zero() *Decimal {
	return &Decimal{Unscaled: 0, Scale: 0}
}

// New returns normalized [Decimal] representing [decimal.Decimal] d.
// If [decimal.Decimal.Coefficient] cannot be represented in an int64, the result will be undefined.
func New(d decimal.Decimal) *Decimal {
	d = RoundByInt64(Normalize(d))
	return &Decimal{
		Unscaled: d.CoefficientInt64(), // If coefficient cannot be represented in an int64, the result will be undefined
		Scale:    -d.Exponent(),
	}
}

// NewFromInt returns normalized [Decimal] representing integer v.
func NewFromInt(v int) *Decimal {
	return New(decimal.NewFromInt(int64(v)))
}

// NewFromFloat returns normalized [Decimal] representing float64 v.
func NewFromFloat(v float64) *Decimal {
	return New(decimal.NewFromFloat(v))
}

// RoundByInt64 rounds [decimal.Decimal.Coefficient] to be representable as int64.
// It makes calling [decimal.Decimal.CoefficientInt64()] safe in almost all cases.
// [decimal.Decimal.Exponent] is changed in such a way to remove all extra magnitude orders
// from [decimal.Decimal.Coefficient] bigger than int64. Rounding is always down.
//
// Examples
//
//	input:  9223372036854775817 * 10^-18 == 9.223372036854775817  (9223372036854775817 == maxInt64 + 10)
//	output: 922337203685477581  * 10^-17 == 9.22337203685477581
//
// Note: This can represent numbers with a maximum of 2^31-1 exponent.
// If the exponent is exceeded and coefficient is greater than int64,
// then the result of [decimal.Decimal.CoefficientInt64()] will be undefined.
func RoundByInt64(d decimal.Decimal) decimal.Decimal {
	value := d.Coefficient()
	if value.IsInt64() {
		return d
	}

	mayValue := new(big.Int)
	n := int64(0)
	maxN := math.MaxInt32 - int64(d.Exponent())
	for {
		if n == maxN {
			break // avoid int32 overflow
		}

		mayValue.Quo(value, ten)

		value.Set(mayValue)
		n++

		if value.IsInt64() {
			break
		}
	}

	return decimal.NewFromBigInt(value, d.Exponent()+int32(n))
}
