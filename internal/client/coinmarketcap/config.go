package coinmarketcap

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	secretKeyEnv = "COINMARKETCAP_SECRET"
)

type Config struct {
	CoinmarketcapHost      string `json:"coinmarketcap_host"`
	CoinmarketcapSecretKey string `json:"coinmarketcap_secret_key"`
}

func readConfig() (Config, error) {
	configPath := "config.json"
	f, err := os.Open(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("failed to open file %s: %v", configPath, err)
	}
	defer f.Close()

	var c Config
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return Config{}, fmt.Errorf("failed to decode config file: %v", err)
	}

	if len(c.CoinmarketcapSecretKey) == 0 {
		secretKey, ok := os.LookupEnv(secretKeyEnv)
		if !ok {
			return Config{}, fmt.Errorf("%s env is not found", secretKeyEnv)
		}
		c.CoinmarketcapSecretKey = secretKey
	}

	return c, nil
}
