package dataset

import (
	"TradingServer/internal/api/grpc/dataset/pb"
	"TradingServer/internal/core/client/backtest/dataset"
	"context"
	"errors"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
)

type _datasetService struct {
	pb.UnimplementedDatasetServiceServer

	backtestDatasetManager *dataset.DatasetManager
	alpacaMarketDataClient *marketdata.Client
}

// mustEmbedUnimplementedDatasetServiceServer implements pb.DatasetServiceServer.
func (t *_datasetService) mustEmbedUnimplementedDatasetServiceServer() {
	panic("unimplemented")
}

func (t *_datasetService) CreateDataset(ctx context.Context, req *pb.CreateDatasetRequest) (*pb.CreateDatasetResponse, error) {

	switch req.Dataset.(type) {
	case *pb.CreateDatasetRequest_Alpaca:
		{
			data, err := dataset.AlpacaDatasetCreate(req.GetAlpaca(), t.alpacaMarketDataClient)
			if err != nil {
				return &pb.CreateDatasetResponse{}, err
			}

			dataset_id := t.backtestDatasetManager.AddDataset(data)

			return &pb.CreateDatasetResponse{
				DatasetId: dataset_id,
			}, nil
		}
	case *pb.CreateDatasetRequest_Csv:
		{
			return &pb.CreateDatasetResponse{}, errors.New("not implemented")
		}
	default:
		{
			return &pb.CreateDatasetResponse{}, errors.New("not implemented")
		}
	}
}

func (t *_datasetService) RemoveDataset(ctx context.Context, req *pb.RemoveDatasetRequest) (*pb.RemoveDatasetResponse, error) {
	return &pb.RemoveDatasetResponse{}, errors.New("not implemented")
}

func (t *_datasetService) GetStartTime(ctx context.Context, req *pb.GetStartTimeRequest) (*pb.GetStartTimeResponse, error) {
	return &pb.GetStartTimeResponse{}, errors.New("not implemented")
}

func (t *_datasetService) GetEndTime(ctx context.Context, req *pb.GetEndTimeRequest) (*pb.GetEndTimeResponse, error) {
	return &pb.GetEndTimeResponse{}, errors.New("not implemented")
}

func (t *_datasetService) GetSymbols(ctx context.Context, req *pb.GetSymbolsRequest) (*pb.GetSymbolsResponse, error) {
	return &pb.GetSymbolsResponse{}, errors.New("not implemented")
}

type DatasetServiceOpts struct {
	BacktestDatasetManager *dataset.DatasetManager
	AlpacaMarketDataClient *marketdata.Client
}

func NewDatasetService(opts DatasetServiceOpts) *_datasetService {
	return &_datasetService{
		backtestDatasetManager: opts.BacktestDatasetManager,
		alpacaMarketDataClient: opts.AlpacaMarketDataClient,
	}
}
