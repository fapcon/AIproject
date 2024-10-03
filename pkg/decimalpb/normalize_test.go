package decimalpb_test

import (
	"math"
	"math/rand"
	"testing"
	"time"

	"studentgit.kata.academy/quant/torque/pkg/decimalpb"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	type testCase struct {
		name    string
		in, out decimal.Decimal
	}

	testCases := []testCase{
		{name: "1", in: decimal.New(0, 0), out: decimal.New(0, 0)},
		{name: "2", in: decimal.New(0, 10), out: decimal.New(0, 0)},
		{name: "3", in: decimal.New(0, -10), out: decimal.New(0, 0)},
		{name: "4", in: decimal.New(10, 0), out: decimal.New(1, 1)},
		{name: "5", in: decimal.New(100, 0), out: decimal.New(1, 2)},
		{name: "6", in: decimal.New(1000, 0), out: decimal.New(1, 3)},
		{name: "7", in: decimal.New(10, 1), out: decimal.New(1, 2)},
		{name: "8", in: decimal.New(100, 1), out: decimal.New(1, 3)},
		{name: "9", in: decimal.New(1000, 1), out: decimal.New(1, 4)},
		{name: "10", in: decimal.New(10, -1), out: decimal.New(1, 0)},
		{name: "11", in: decimal.New(100, -1), out: decimal.New(1, 1)},
		{name: "12", in: decimal.New(1000, -1), out: decimal.New(1, 2)},
		{name: "13", in: decimal.New(-10, 0), out: decimal.New(-1, 1)},
		{name: "14", in: decimal.New(-100, 0), out: decimal.New(-1, 2)},
		{name: "15", in: decimal.New(-1000, 0), out: decimal.New(-1, 3)},
		{name: "16", in: decimal.New(-10, 1), out: decimal.New(-1, 2)},
		{name: "17", in: decimal.New(-100, 1), out: decimal.New(-1, 3)},
		{name: "18", in: decimal.New(-1000, 1), out: decimal.New(-1, 4)},
		{name: "19", in: decimal.New(-10, -1), out: decimal.New(-1, 0)},
		{name: "20", in: decimal.New(-100, -1), out: decimal.New(-1, 1)},
		{name: "21", in: decimal.New(-1000, -1), out: decimal.New(-1, 2)},
		{name: "22", in: decimal.New(1234, 2), out: decimal.New(1234, 2)},
		{name: "23", in: decimal.New(1234, -2), out: decimal.New(1234, -2)},
		{name: "24", in: decimal.New(1234, math.MaxInt32), out: decimal.New(1234, math.MaxInt32)},
		{name: "25", in: decimal.New(12340, math.MaxInt32), out: decimal.New(12340, math.MaxInt32)},
		{name: "26", in: decimal.New(123400, math.MaxInt32-1), out: decimal.New(12340, math.MaxInt32)},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// assert.True(t, tc.out.Equal(decimalpb.Normalize(tc.in)))
			assert.Equal(t, tc.out, decimalpb.Normalize(tc.in))
		})
	}
}

var xxx decimal.Decimal //nolint:gochecknoglobals // test

// 225.0 ns/op.
func BenchmarkNormalize(b *testing.B) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	var x decimal.Decimal
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		d := decimal.New(r.Int63(), r.Int31())
		b.StartTimer()
		x = decimalpb.Normalize(d)
	}
	xxx = x
}
