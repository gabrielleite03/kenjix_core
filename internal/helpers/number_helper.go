package helpers

import "github.com/shopspring/decimal"

func DecimalFromFloat(v float64) decimal.Decimal {
	return decimal.NewFromFloat(v)
}
