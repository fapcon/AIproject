package decimalpb_test

import (
	"math"
	"math/big"
	"testing"

	"studentgit.kata.academy/quant/torque/pkg/decimalpb"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestRound(t *testing.T) {
	t.Parallel()

	// positives

	t.Run("big positive value, positive exp", func(t *testing.T) {
		t.Parallel()

		// 9223372036854775807 + 10
		bigValue := new(big.Int).Add(big.NewInt(math.MaxInt64), big.NewInt(10))
		value := decimal.NewFromBigInt(bigValue, 10)

		expected := decimal.New(922337203685477581, 11)
		actual := decimalpb.RoundByInt64(value)

		// input:  9223372036854775817 * 10^10
		// output: 922337203685477581  * 10^11
		assert.Equal(t, expected, actual)
	})
	t.Run("big positive value, negative exp", func(t *testing.T) {
		t.Parallel()

		// 9223372036854775807 + 4
		bigValue := new(big.Int).Add(big.NewInt(math.MaxInt64), big.NewInt(4))
		value := decimal.NewFromBigInt(bigValue, -16)

		expected := decimal.New(922337203685477581, -15)
		actual := decimalpb.RoundByInt64(value)

		// input:  9223372036854775811 * 10^-16 == 922.3372036854775811
		// output: 922337203685477581  * 10^-15 == 922.337203685477581
		assert.Equal(t, expected, actual)
	})
	t.Run("big positive value, zero exp", func(t *testing.T) {
		t.Parallel()

		// 9223372036854775807 + 4
		bigValue := new(big.Int).Add(big.NewInt(math.MaxInt64), big.NewInt(4))
		value := decimal.NewFromBigInt(bigValue, 0)

		expected := decimal.New(922337203685477581, 1)
		actual := decimalpb.RoundByInt64(value)

		// input:  9223372036854775811 * 10^0 == 9223372036854775811
		// output: 922337203685477581  * 10^1 == 9223372036854775810
		assert.Equal(t, expected, actual)
	})
	t.Run("big positive value (1000 * max int), positive exp", func(t *testing.T) {
		t.Parallel()

		// 9223372036854775807 * 1001 == 9232595408891630582807
		bigValue := new(big.Int).Mul(big.NewInt(math.MaxInt64), big.NewInt(1001))
		value := decimal.NewFromBigInt(bigValue, 10)

		expected := decimal.New(923259540889163058, 14)
		actual := decimalpb.RoundByInt64(value)

		// input:  9232595408891630582807 * 10^10
		// output: 923259540889163058     * 10^14
		assert.Equal(t, expected, actual)
	})
	t.Run("big positive value (1000 * max int), negative exp", func(t *testing.T) {
		t.Parallel()

		// 9223372036854775807 * 1001 == 9232595408891630582807
		bigValue := new(big.Int).Mul(big.NewInt(math.MaxInt64), big.NewInt(1001))
		value := decimal.NewFromBigInt(bigValue, -22)

		expected := decimal.New(923259540889163058, -18)
		actual := decimalpb.RoundByInt64(value)

		// input:  9232595408891630582807 * 10^-22 == 0.9232595408891630582807
		// output: 923259540889163058     * 10^-18 == 0.923259540889163058
		assert.Equal(t, expected, actual)
	})
	t.Run("max int64 value, zero exp", func(t *testing.T) {
		t.Parallel()

		// 9223372036854775807
		bigValue := new(big.Int).Mul(big.NewInt(math.MaxInt64), big.NewInt(1))
		value := decimal.NewFromBigInt(bigValue, 0)

		expected := decimal.New(9223372036854775807, 0)
		actual := decimalpb.RoundByInt64(value)

		// input:  9223372036854775807
		// output: 9223372036854775817
		assert.Equal(t, expected, actual)
	})

	// negatives

	t.Run("big negative value, positive exp", func(t *testing.T) {
		t.Parallel()

		// -9223372036854775808 - 4
		bigValue := new(big.Int).Sub(big.NewInt(math.MinInt64), big.NewInt(4))
		value := decimal.NewFromBigInt(bigValue, 10)

		expected := decimal.New(-922337203685477581, 11)
		actual := decimalpb.RoundByInt64(value)

		// input:  -9223372036854775812 * 10^10
		// output: -922337203685477581  * 10^11
		assert.Equal(t, expected, actual)
	})
	t.Run("big negative value, negative exp", func(t *testing.T) {
		t.Parallel()

		// -9223372036854775808 - 4
		bigValue := new(big.Int).Sub(big.NewInt(math.MinInt64), big.NewInt(4))
		value := decimal.NewFromBigInt(bigValue, -16)

		expected := decimal.New(-922337203685477581, -15)
		actual := decimalpb.RoundByInt64(value)

		// input:  -9223372036854775812 * 10^-16 == -922.3372036854775811
		// output: -922337203685477581  * 10^-15 == -922.337203685477581
		assert.Equal(t, expected, actual)
	})
	t.Run("big negative value, zero exp", func(t *testing.T) {
		t.Parallel()

		// -9223372036854775808 - 4
		bigValue := new(big.Int).Sub(big.NewInt(math.MinInt64), big.NewInt(4))
		value := decimal.NewFromBigInt(bigValue, 0)

		expected := decimal.New(-922337203685477581, 1)
		actual := decimalpb.RoundByInt64(value)

		// input:  -9223372036854775812 * 10^0
		// output: -922337203685477581  * 10^1
		assert.Equal(t, expected, actual)
	})
	t.Run("big negative value (1000 * max int), positive exp", func(t *testing.T) {
		t.Parallel()

		// -9223372036854775808 * 1001 == -9232595408891630583808
		bigValue := new(big.Int).Mul(big.NewInt(math.MinInt64), big.NewInt(1001))
		value := decimal.NewFromBigInt(bigValue, 10)

		expected := decimal.New(-923259540889163058, 14)
		actual := decimalpb.RoundByInt64(value)

		// input:  -9232595408891630583808 * 10^10
		// output: -923259540889163058     * 10^14
		assert.Equal(t, expected, actual)
	})
	t.Run("big negative value (1000 * max int), negative exp", func(t *testing.T) {
		t.Parallel()

		// -9223372036854775808 * 1001 == -9232595408891630583808
		bigValue := new(big.Int).Mul(big.NewInt(math.MinInt64), big.NewInt(1001))
		value := decimal.NewFromBigInt(bigValue, -22)

		expected := decimal.New(-923259540889163058, -18)
		actual := decimalpb.RoundByInt64(value)

		// input:  -9232595408891630583808 * 10^-22 == -0.9232595408891630583808
		// output: -923259540889163058     * 10^-18 == -0.923259540889163058
		assert.Equal(t, expected, actual)
	})
	t.Run("min int64 value, zero exp", func(t *testing.T) {
		t.Parallel()

		// -9223372036854775808
		bigValue := new(big.Int).Mul(big.NewInt(math.MinInt64), big.NewInt(1))
		value := decimal.NewFromBigInt(bigValue, 0)

		expected := decimal.New(-9223372036854775808, 0)
		actual := decimalpb.RoundByInt64(value)

		// input:  -9223372036854775808
		// output: -9223372036854775808
		assert.Equal(t, expected, actual)
	})

	// extra cases

	t.Run("zero value", func(t *testing.T) {
		t.Parallel()
		value := decimal.Zero

		expected := decimal.Zero
		actual := decimalpb.RoundByInt64(value)

		assert.Equal(t, expected, actual)
	})
	t.Run("max exponent", func(t *testing.T) {
		t.Parallel()

		// 9223372036854775807 + 10
		bigValue := new(big.Int).Add(big.NewInt(math.MaxInt64), big.NewInt(10))
		value := decimal.NewFromBigInt(bigValue, math.MaxInt32)

		expected := decimal.NewFromBigInt(bigValue, math.MaxInt32)
		actual := decimalpb.RoundByInt64(value)

		// input:  9223372036854775817 * 10^MaxInt32
		// output: 9223372036854775817 * 10^MaxInt32
		// exponent overflows, coefficient can't be rounded to int64 presentation
		assert.Equal(t, expected, actual)
	})
}
