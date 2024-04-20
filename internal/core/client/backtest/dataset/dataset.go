package dataset

import (
	"TradingServer/internal/pkg"
	"time"
)

type Dataset struct {
	dataset map[uint32]struct {
		symbol pkg.Symbol
		data   pkg.TimeSeries[pkg.Candle]
	}
}

func (t *Dataset) SymbolCandlesBetween(symbol_id uint32, start time.Time, end time.Time) (pkg.TimeSeries[pkg.Candle], error) {
	symbol, ok := t.dataset[symbol_id]
	if !ok {
		return pkg.TimeSeries[pkg.Candle]{}, SymbolNotFoundError{}
	}

	return symbol.data.GetDataPointsWithin(start, end), nil
}

func (t *Dataset) Symbols() map[uint32]pkg.Symbol {
	symbols := map[uint32]pkg.Symbol{}
	for symbol_id, symbol := range t.dataset {
		symbols[symbol_id] = symbol.symbol
	}
	return symbols
}

func (t *Dataset) Starttime() time.Time {
	min := time.Now()
	for _, symbol := range t.dataset {
		if symbol.data.Timestamps[0].Before(min) {
			min = symbol.data.Timestamps[0]
		}
	}
	return min
}

func (t *Dataset) Endtime() time.Time {
	max := time.Now()
	for _, symbol := range t.dataset {
		if symbol.data.Timestamps[len(symbol.data.Timestamps)-1].After(max) {
			max = symbol.data.Timestamps[len(symbol.data.Timestamps)-1]
		}
	}
	return max
}

type DatasetOpts struct {
}

func NewDataset(opts DatasetOpts) *Dataset {
	return &Dataset{}
}
