package client

import (
	"TradingServer/internal/api/grpc/client/pb"
	"TradingServer/internal/core/client/backtest"
	"context"
	"errors"

	"github.com/shopspring/decimal"
)

func DecimalToProtoDecimal(d decimal.Decimal) *pb.Decimal {
	return &pb.Decimal{
		Coefficent: d.CoefficientInt64(),
		Exponent:   d.Exponent(),
	}
}

type _clientService struct {
	pb.UnimplementedClientServiceServer

	backtestClientManager *backtest.BacktestClientManager
}

func (t *_clientService) RemoveClient(ctx context.Context, req *pb.RemoveClientRequest) (*pb.RemoveClientResponse, error) {
	return &pb.RemoveClientResponse{}, errors.New("not implemented")
}

func (t *_clientService) StreamOpenCandles(req *pb.StreamOpenCandlesReq, stream pb.ClientService_StreamOpenCandlesServer) error {
	client, err := t.backtestClientManager.GetClient(req.GetClientId())
	if err != nil {
		return err
	}

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return nil
		case candle, ok := <-client.OpenCandles():
			if !ok {
				return errors.New("clients open candles channel closed")
			}

			ret := &pb.StreamOpenCandlesRet{
				OpenCandles: map[uint32]*pb.StreamOpenCandlesRet_OpenCandle{
					candle.SymbolId: &pb.StreamOpenCandlesRet_OpenCandle{
						Price:  DecimalToProtoDecimal(candle.Candle.Price),
						High:   DecimalToProtoDecimal(candle.Candle.High),
						Low:    DecimalToProtoDecimal(candle.Candle.Low),
						Open:   DecimalToProtoDecimal(candle.Candle.Open),
						Volume: DecimalToProtoDecimal(candle.Candle.Volume),
					},
				},
			}

			err := stream.Send(ret)
			if err != nil {
				return err
			}
		}
	}
}

func (t *_clientService) StreamCloseCandles(req *pb.StreamClosedCandlesReq, stream pb.ClientService_StreamClosedCandlesServer) error {
	client, err := t.backtestClientManager.GetClient(req.GetClientId())
	if err != nil {
		return err
	}

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return nil
		case candle, ok := <-client.ClosedCandles():
			if !ok {
				return errors.New("clients closed candles channel closed")
			}

			ret := &pb.StreamClosedCandlesRet{
				ClosedCandles: map[uint32]*pb.StreamClosedCandlesRet_ClosedCandle{
					candle.SymbolId: &pb.StreamClosedCandlesRet_ClosedCandle{
						High:   DecimalToProtoDecimal(candle.Candle.High),
						Low:    DecimalToProtoDecimal(candle.Candle.Low),
						Open:   DecimalToProtoDecimal(candle.Candle.Open),
						Close:  DecimalToProtoDecimal(candle.Candle.Close),
						Volume: DecimalToProtoDecimal(candle.Candle.Volume),
					},
				},
			}

			err := stream.Send(ret)
			if err != nil {
				return err
			}
		}
	}
}

func (t *_clientService) GetHistoricalCandles(ctx context.Context, req *pb.GetHistoricalCandlesReq) (*pb.GetHistoricalCandlesRet, error) {
	return &pb.GetHistoricalCandlesRet{}, errors.New("not implemented")
}

type ClientServiceOpts struct {
	BacktestClientManager *backtest.BacktestClientManager
}

func NewClientService(opts ClientServiceOpts) *_clientService {
	return &_clientService{
		backtestClientManager: opts.BacktestClientManager,
	}
}
