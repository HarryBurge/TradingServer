package client

import (
	"TradingServer/internal/api/grpc/client/pb"
	"TradingServer/internal/core/client/backtest"
	"TradingServer/internal/core/client/backtest/dataset"
	"context"
)

type _backtestClientService struct {
	pb.UnimplementedBacktestClientServiceServer

	backtestClientManager *backtest.BacktestClientManager
	datasetManager        *dataset.DatasetManager
}

func (t *_backtestClientService) CreateClient(ctx context.Context, req *pb.CreateBacktestClientRequest) (*pb.CreateBacktestClientResponse, error) {

	dataset, err := t.datasetManager.GetDataset(req.GetDatasetId())
	if err != nil {
		return &pb.CreateBacktestClientResponse{}, err
	}

	client_id := t.backtestClientManager.AddClient(
		backtest.NewBacktestClient(backtest.BacktestClientOpts{
			Dataset: dataset,
		}),
	)

	return &pb.CreateBacktestClientResponse{
		ClientId: client_id,
	}, nil
}

func (t *_backtestClientService) StepToTime(ctx context.Context, req *pb.StepToTimeReq) (*pb.StepToTimeRet, error) {
	client, err := t.backtestClientManager.GetClient(req.GetClientId())
	if err != nil {
		return &pb.StepToTimeRet{}, err
	}

	err = client.StepToTime(req.GetTime().AsTime())
	if err != nil {
		return &pb.StepToTimeRet{}, err
	}

	return &pb.StepToTimeRet{}, nil
}

type BacktestClientServiceOpts struct {
	BacktestClientManager *backtest.BacktestClientManager
	DatasetManager        *dataset.DatasetManager
}

func NewBacktestClientService(opts BacktestClientServiceOpts) *_backtestClientService {
	return &_backtestClientService{
		backtestClientManager: opts.BacktestClientManager,
		datasetManager:        opts.DatasetManager,
	}
}
