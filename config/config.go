package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	AlpacaConfig AlpacaConfig `json:"alpaca"`
}

func ParseConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var config Config
	_ = json.Unmarshal(data, &config)

	return config, err
}

func WriteDefaultConfig(path string) error {
	data, err := json.MarshalIndent(
		Config{
			AlpacaConfig: AlpacaConfig{
				APIKey: "YOUR_API_KEY",
				APISecret: "YOUR_API_SECRET",
				APIURL: "https://paper-api.alpaca.markets",
			},
		},
		"",
		"\t",
	)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(data)

	return nil
}