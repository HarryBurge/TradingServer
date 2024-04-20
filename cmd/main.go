package main

import (
	"TradingServer/config"
	grpc_dataset "TradingServer/internal/api/grpc/dataset"
	grpc_dataset_pb "TradingServer/internal/api/grpc/dataset/pb"
	"TradingServer/internal/core/client/backtest"
	"TradingServer/internal/core/client/backtest/dataset"

	grpc_client "TradingServer/internal/api/grpc/client"
	grpc_client_pb "TradingServer/internal/api/grpc/client/pb"

	"TradingServer/internal/di"
	"errors"
	"flag"
	"net"
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"google.golang.org/grpc"
)

var (
	config_path             = flag.String("config_path", "", "path to config file")
	generate_default_config = flag.Bool("generate_default_config", false, "generate a default config, at the config path specified")
)

// ----------------------------------------------------------------
// marketclient := marketdata.NewClient(marketdata.ClientOpts{
// })

// bars, err := marketclient.GetBars("META", marketdata.GetBarsRequest{
// 	TimeFrame: marketdata.OneDay,
// 	Start:     time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC),
// 	End:       time.Date(2022, 6, 22, 0, 0, 0, 0, time.UTC),
// 	AsOf:      "2022-06-10", // Leaving it empty yields the same results
// })
// if err != nil {
// 	panic(err)
// }
// fmt.Println("META bars:")
// for _, bar := range bars {
// 	fmt.Printf("%+v\n", bar)
// }
// ----------------------------------------------------------------

func flagsParseValidation() error {
	flag.Parse()

	if *config_path == "" {
		return errors.New("config_path must be specified")
	}

	if *generate_default_config {
		if _, err := os.Stat(*config_path); err == nil {
			return errors.New("the path in config_path already exists")
		}
	} else {
		if _, err := os.Stat(*config_path); err != nil {
			return errors.New("the path in config_path doesn't exist")
		}
	}

	return nil
}

func main() {

	// Validate flags parsed from the command line
	err := flagsParseValidation()
	if err != nil {
		panic(err)
	}

	// Generate default configuration
	if *generate_default_config {
		err := config.WriteDefaultConfig(*config_path)
		if err != nil {
			panic(err)
		}
		return
	}

	// Read configuration
	server_config, err := config.ParseConfig(*config_path)
	if err != nil {
		panic(err)
	}

	// Generate di
	ctn := di.New(server_config)

	// GRPC registration
	lis, _ := net.Listen("tcp", "localhost:50051")
	var opts []grpc.ServerOption
	grpc_server := grpc.NewServer(opts...)

	grpc_dataset_pb.RegisterDatasetServiceServer(
		grpc_server,
		grpc_dataset.NewDatasetService(grpc_dataset.DatasetServiceOpts{
			BacktestDatasetManager: ctn.Get(di.BacktestDatasetManager).(*dataset.DatasetManager),
			AlpacaMarketDataClient: ctn.Get(di.AlpacaMarketDataClient).(*marketdata.Client),
		}),
	)

	grpc_client_pb.RegisterClientServiceServer(
		grpc_server,
		grpc_client.NewClientService(grpc_client.ClientServiceOpts{
			BacktestClientManager: ctn.Get(di.BacktestClientManager).(*backtest.BacktestClientManager),
		}),
	)

	grpc_client_pb.RegisterBacktestClientServiceServer(
		grpc_server,
		grpc_client.NewBacktestClientService(grpc_client.BacktestClientServiceOpts{
			BacktestClientManager: ctn.Get(di.BacktestClientManager).(*backtest.BacktestClientManager),
			DatasetManager:        ctn.Get(di.BacktestDatasetManager).(*dataset.DatasetManager),
		}),
	)

	// Service GRPC Requests
	if err := grpc_server.Serve(lis); err != nil {
		panic(err)
	}
}
