package dataset

import (
	"TradingServer/internal/api/grpc/dataset/pb"
	"TradingServer/internal/pkg"
	"time"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/shopspring/decimal"
)

func _grpcTimeFrameToAlpacaTimeFrameUnit(time_frame_unit pb.CreateDatasetRequest_AlpacaDataset_SymbolQuery_TimeFrameUnit) marketdata.TimeFrameUnit {
	switch time_frame_unit {
	case pb.CreateDatasetRequest_AlpacaDataset_SymbolQuery_HOUR:
		return marketdata.Hour
	case pb.CreateDatasetRequest_AlpacaDataset_SymbolQuery_DAY:
		return marketdata.Day
	case pb.CreateDatasetRequest_AlpacaDataset_SymbolQuery_WEEK:
		return marketdata.Week
	case pb.CreateDatasetRequest_AlpacaDataset_SymbolQuery_MONTH:
		return marketdata.Month
	default:
		return marketdata.Min
	}
}

func AlpacaDatasetCreate(msg *pb.CreateDatasetRequest_AlpacaDataset, client *marketdata.Client) (*Dataset, error) {

	dataset := make(map[uint32]struct {
		symbol pkg.Symbol
		data   pkg.TimeSeries[pkg.Candle]
	})

	for symbol_id, symbol_query := range msg.GetSymbols() {
		switch symbol_query.SymbolType {
		case pb.CreateDatasetRequest_AlpacaDataset_SymbolQuery_CRYPTO:
			{
				// Collection
				bars, err := client.GetCryptoBars(
					symbol_query.GetSymbolName(),
					marketdata.GetCryptoBarsRequest{
						TimeFrame: marketdata.NewTimeFrame(1, _grpcTimeFrameToAlpacaTimeFrameUnit(symbol_query.GetTimeFrameUnit())),
						Start:     symbol_query.GetStartTime().AsTime(),
						End:       symbol_query.GetEndTime().AsTime(),
					},
				)
				if err != nil {
					return &Dataset{}, err
				}

				// Conversion
				candles := make([]pkg.Candle, len(bars))
				timestamps := make([]time.Time, len(bars))

				for i, bar := range bars {
					candles[i] = pkg.Candle{
						Open:   decimal.NewFromFloat(bar.Open),
						High:   decimal.NewFromFloat(bar.High),
						Low:    decimal.NewFromFloat(bar.Low),
						Close:  decimal.NewFromFloat(bar.Close),
						Volume: decimal.NewFromFloat(bar.Volume),
					}
					timestamps[i] = bar.Timestamp
				}

				// Store
				dataset[uint32(symbol_id)] = struct {
					symbol pkg.Symbol
					data   pkg.TimeSeries[pkg.Candle]
				}{
					symbol: pkg.Symbol{
						Name: symbol_query.GetSymbolName(),
					},
					data: pkg.TimeSeries[pkg.Candle]{
						Timestamps: timestamps,
						Datapoints: candles,
					},
				}
			}
		case pb.CreateDatasetRequest_AlpacaDataset_SymbolQuery_STOCK:
			{
				// Collection
				bars, err := client.GetBars(
					symbol_query.GetSymbolName(),
					marketdata.GetBarsRequest{
						TimeFrame: marketdata.NewTimeFrame(1, _grpcTimeFrameToAlpacaTimeFrameUnit(symbol_query.GetTimeFrameUnit())),
						Start:     symbol_query.GetStartTime().AsTime(),
						End:       symbol_query.GetEndTime().AsTime(),
					},
				)
				if err != nil {
					return &Dataset{}, err
				}

				// Conversion
				candles := make([]pkg.Candle, len(bars))
				timestamps := make([]time.Time, len(bars))

				for i, bar := range bars {
					candles[i] = pkg.Candle{
						Open:   decimal.NewFromFloat(bar.Open),
						High:   decimal.NewFromFloat(bar.High),
						Low:    decimal.NewFromFloat(bar.Low),
						Close:  decimal.NewFromFloat(bar.Close),
						Volume: decimal.NewFromUint64(bar.Volume),
					}
					timestamps[i] = bar.Timestamp
				}

				// Store
				dataset[uint32(symbol_id)] = struct {
					symbol pkg.Symbol
					data   pkg.TimeSeries[pkg.Candle]
				}{
					symbol: pkg.Symbol{
						Name: symbol_query.GetSymbolName(),
					},
					data: pkg.TimeSeries[pkg.Candle]{
						Timestamps: timestamps,
						Datapoints: candles,
					},
				}
			}
		}
	}

	return &Dataset{
		dataset: dataset,
	}, nil
}
