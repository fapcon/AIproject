package domain

import (
	"github.com/shopspring/decimal"
)

type Balance struct {
	Symbol    Symbol
	Available decimal.Decimal
}
