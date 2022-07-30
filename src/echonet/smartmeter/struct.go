package smartmeter

import "github.com/shopspring/decimal"

type NowDenryuu struct {
	Rphase decimal.Decimal
	Tphase decimal.Decimal
	Total  decimal.Decimal
}
