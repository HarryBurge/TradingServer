package backtest

import (
	"TradingServer/internal/core/client/backtest/dataset"
	"TradingServer/internal/pkg"
	"time"

	"github.com/shopspring/decimal"
)

// type BacktestClient interface {
// 	client.Client
// 	StepToTime(time time.Time) error
// }

// TODO h: Rename to Backtest Engine
type BacktestClient struct {
	open_candles chan struct {
		SymbolId uint32
		Candle   pkg.OpenCandle
	}

	closed_candles chan struct {
		SymbolId uint32
		Candle   pkg.Candle
	}

	dataset *dataset.Dataset

	current_time time.Time
	cash         decimal.Decimal
	orders       map[uint32][]pkg.Order
	trades       map[uint32][]pkg.Trade
}

func (t *BacktestClient) OpenCandles() chan struct {
	SymbolId uint32
	Candle   pkg.OpenCandle
} {
	return t.open_candles
}

func (t *BacktestClient) ClosedCandles() chan struct {
	SymbolId uint32
	Candle   pkg.Candle
} {
	return t.closed_candles
}

func (t *BacktestClient) SymbolCandlesPrevious(symbol_id uint32, duration time.Duration) (pkg.TimeSeries[pkg.Candle], error) {
	return t.dataset.SymbolCandlesBetween(symbol_id, t.current_time.Add(-duration), t.current_time)
}

func (t *BacktestClient) SymbolCandlesBetween(symbol_id uint32, starttime time.Time, endtime time.Time) (pkg.TimeSeries[pkg.Candle], error) {
	if endtime.After(t.current_time) {
		return pkg.TimeSeries[pkg.Candle]{}, IncorrectTimeError{}
	}

	return t.dataset.SymbolCandlesBetween(symbol_id, starttime, endtime)
}

func (t *BacktestClient) positionSize(symbol_id uint32) decimal.Decimal {
	trades, _ := t.trades[symbol_id]

	position_size := decimal.Zero
	for _, trade := range trades {
		position_size.Add(trade.Size)
	}

	return position_size
}

func (t *BacktestClient) StepToTime(time time.Time) error {
	if time.Before(t.current_time) {
		return IncorrectTimeError{}
	}

	// Get all symbol data between the current time and the given time
	symbol_timerange_candles := make(map[uint32]pkg.TimeSeries[pkg.Candle])
	for symbol_id, _ := range t.dataset.Symbols() {
		symbol_timerange_candles[symbol_id], _ = t.dataset.SymbolCandlesBetween(symbol_id, t.current_time, time)
	}

	// Intialise a map of iterators for each symbol
	symbol_index_iterators := make(map[uint32]int)
	for symbol_id, _ := range symbol_timerange_candles {
		symbol_index_iterators[symbol_id] = 0
	}

	// Check whether all iterators are at the end
	at_end := func() bool {
		for symbol_id, idx := range symbol_index_iterators {
			if idx < len(symbol_timerange_candles[symbol_id].Datapoints) {
				return false
			}
		}
		return true
	}

	// Iterate through finding min time candle, and emitting it
	for !at_end() {
		min_time_symbol_id := uint32(0)
		min_time := symbol_timerange_candles[0].Timestamps[symbol_index_iterators[0]]
		for symbol_id, timerange_candles := range symbol_timerange_candles {
			if timerange_candles.Timestamps[symbol_index_iterators[symbol_id]].Before(min_time) {
				min_time = timerange_candles.Timestamps[symbol_index_iterators[symbol_id]]
				min_time_symbol_id = symbol_id
			}
		}

		// Emit closed candle
		t.closed_candles <- struct {
			SymbolId uint32
			Candle   pkg.Candle
		}{
			SymbolId: min_time_symbol_id,
			Candle:   symbol_timerange_candles[min_time_symbol_id].Datapoints[symbol_index_iterators[min_time_symbol_id]],
		}
		symbol_index_iterators[min_time_symbol_id]++

		// Fill orders on new candle
		orders, ok := t.orders[min_time_symbol_id]
		if ok {
			for _, order := range orders {
				can_be_filled := false

				price := symbol_timerange_candles[min_time_symbol_id].Datapoints[symbol_index_iterators[min_time_symbol_id]].Open

				// Check order can be filled
				if order.Size.IsPositive() && t.cash.GreaterThanOrEqual(price) {
					can_be_filled = true
				} else if order.Size.IsNegative() && t.positionSize(min_time_symbol_id).GreaterThanOrEqual(order.Size) {
					can_be_filled = true
				}

				if can_be_filled {
					t.trades[min_time_symbol_id] = append(t.trades[min_time_symbol_id], pkg.Trade{
						Size:  order.Size,
						Price: price,
					})
					t.cash = t.cash.Add(order.Size.Mul(price).Neg())
				}
			}
			delete(t.orders, min_time_symbol_id)
		}
	}

	// As all candles have been emmitted, set the current time to the given time
	t.current_time = time
	return nil
}

func (t *BacktestClient) PlaceOrder(symbol_id uint32, order pkg.Order) error {
	t.orders[symbol_id] = append(t.orders[symbol_id], order)
	return nil
}

type BacktestClientOpts struct {
	Dataset *dataset.Dataset
}

func NewBacktestClient(opts BacktestClientOpts) *BacktestClient {
	return &BacktestClient{
		dataset:      opts.Dataset,
		current_time: opts.Dataset.Starttime(),
		closed_candles: make(chan struct {
			SymbolId uint32
			Candle   pkg.Candle
		}),
		open_candles: make(chan struct {
			SymbolId uint32
			Candle   pkg.OpenCandle
		}),
	}
}
