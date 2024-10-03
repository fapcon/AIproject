package decimalpb_test

import (
	"math"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"studentgit.kata.academy/quant/torque/pkg/decimalpb"
)

func FuzzNormalize(f *testing.F) {
	for _, val := range []int64{0, 1, 10, 100, 1000, 1_000_000_000} {
		for _, valS := range []int64{1, -1} {
			for _, exp := range []int32{0, 1, 10, 100, 1000, 1_000_000_000} {
				for _, expS := range []int32{1, -1} {
					f.Add(val*valS, exp*expS)
				}
			}
		}
	}
	for _, val := range []int64{math.MinInt64, math.MaxInt64, 1, 10, 100, 1000, -1, -10, -100, -1000} {
		for _, exp := range []int32{math.MinInt32, math.MaxInt32, math.MinInt32 + 1, math.MaxInt32 - 1} {
			f.Add(val, exp)
		}
	}
	f.Fuzz(func(t *testing.T, val int64, exp int32) { //nolint:thelper // false positive
		x := decimal.New(val, exp)
		y := decimalpb.Normalize(x)

		xc := x.Coefficient()
		xe := x.Exponent()
		yc := y.Coefficient()
		ye := y.Exponent()

		if x.IsZero() {
			// Using y.IsZero() because decimal.New(0, 1_000_000).Equal(decimal.Zero) takes too long to execute
			assert.Truef(t, y.IsZero(), "x = %s * 10 ^ %d\ny = %s * 10 ^ %d", xc, xe, yc, ye)

			assert.Zerof(t, ye, "x = %s * 10 ^ %d\ny = %s * 10 ^ %d", xc, xe, yc, ye)
			return
		}

		assert.Truef(t, y.Equal(x), "x = %s * 10 ^ %d\ny = %s * 10 ^ %d", xc, xe, yc, ye)
		assert.GreaterOrEqualf(t, ye, xe, "x = %s * 10 ^ %d\ny = %s * 10 ^ %d", xc, xe, yc, ye)

		if new(big.Int).Rem(yc, big.NewInt(10)).Sign() == 0 {
			assert.EqualValuesf(t, math.MaxInt32, ye, "x = %s * 10 ^ %d\ny = %s * 10 ^ %d", xc, xe, yc, ye)
		}
	})
}
