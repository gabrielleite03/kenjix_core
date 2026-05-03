package mercadolivre

import (
	"errors"

	"github.com/shopspring/decimal"
)

type Input struct {
	Cost            decimal.Decimal // custo produto
	Operational     decimal.Decimal // custo operacional base
	Weight          decimal.Decimal // kg
	MarginPercent   decimal.Decimal // ex: 0.20
	MLFeePercent    decimal.Decimal // ex: 0.175
	UseFreeShipping bool            // >= 79
	AdsPercent      decimal.Decimal // opcional (ex: 0.05)
}

// Simulação simples de frete (ajuste conforme sua operação)
func estimateShipping(weight decimal.Decimal) decimal.Decimal {
	if weight.LessThan(decimal.NewFromFloat(0.5)) {
		return decimal.NewFromFloat(18)
	}
	if weight.LessThan(decimal.NewFromFloat(1)) {
		return decimal.NewFromFloat(22)
	}
	if weight.LessThan(decimal.NewFromFloat(5)) {
		return decimal.NewFromFloat(35)
	}
	return decimal.NewFromFloat(60)
}

// Simulação custo operacional para produtos baratos
func estimateOperationalBelow79(weight decimal.Decimal) decimal.Decimal {
	if weight.LessThan(decimal.NewFromFloat(0.5)) {
		return decimal.NewFromFloat(6)
	}
	if weight.LessThan(decimal.NewFromFloat(1)) {
		return decimal.NewFromFloat(8)
	}
	return decimal.NewFromFloat(12)
}

func CalculatePrice(input Input) (decimal.Decimal, error) {
	one := decimal.NewFromInt(1)

	// Ajuste operacional para produtos baratos
	operational := input.Operational

	if !input.UseFreeShipping {
		operational = operational.Add(estimateOperationalBelow79(input.Weight))
	}

	// Frete (se aplicável)
	shipping := decimal.Zero
	if input.UseFreeShipping {
		shipping = estimateShipping(input.Weight)
	}

	// Base de custo
	baseCost := input.Cost.
		Add(operational).
		Add(shipping)

	// Denominador
	totalPercent := input.MLFeePercent.
		Add(input.MarginPercent).
		Add(input.AdsPercent)

	denominator := one.Sub(totalPercent)

	if denominator.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, errors.New("taxas + margem inválidas")
	}

	price := baseCost.Div(denominator)

	return price.Round(2), nil
}
