package pkg

import (
	"github.com/shopspring/decimal"
)

type Candle struct {
	High   decimal.Decimal
	Low    decimal.Decimal
	Open   decimal.Decimal
	Close  decimal.Decimal
	Volume decimal.Decimal
}

type OpenCandle struct {
	High   decimal.Decimal
	Low    decimal.Decimal
	Open   decimal.Decimal
	Price  decimal.Decimal
	Volume decimal.Decimal
}
