package di

import (
	"TradingServer/config"
	"TradingServer/internal/core/client/backtest"
	"TradingServer/internal/core/client/backtest/dataset"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/sarulabs/di"
)

const (
	BacktestClientManager  string = "backtest_client_manager"
	BacktestDatasetManager string = "backtest_dataset_manager"
	AlpacaMarketDataClient string = "alpaca_market_data_client"
)

func New(config config.Config) di.Container {
	builder, _ := di.NewBuilder()

	builder.Add([]di.Def{
		{
			Name:  BacktestClientManager,
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				return backtest.NewClientManager(backtest.ClientManagerOpts{}), nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  BacktestDatasetManager,
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				return dataset.NewDatasetManager(dataset.DatasetManagerOpts{}), nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  AlpacaMarketDataClient,
			Scope: di.App,
			Build: func(ctn di.Container) (interface{}, error) {
				return marketdata.NewClient(marketdata.ClientOpts{
					APIKey:    config.AlpacaConfig.APIKey,
					APISecret: config.AlpacaConfig.APISecret,
				}), nil
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
	}...)

	return builder.Build()
}
