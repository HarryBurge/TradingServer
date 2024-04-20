package client

import (
	// "TradingServer/internal/pkg"
	// "time"
)

// type Client interface {
// 	OpenCandles() chan struct {
// 		SymbolId uint32
// 		Candle   pkg.OpenCandle
// 	}
// 	ClosedCandles() chan struct {
// 		SymbolId uint32
// 		Candle   pkg.Candle
// 	}

// 	SymbolCandlesBetween(symbol_id uint32, starttime time.Time, endtime time.Time) (pkg.TimeSeries[pkg.Candle], error)
// 	SymbolCandlesPrevious(symbol_id uint32, duration time.Duration) (pkg.TimeSeries[pkg.Candle], error)
// }