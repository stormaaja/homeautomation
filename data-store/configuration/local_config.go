package configuration

import (
	"log"
	"stormaaja/go-ha/data-store/tools"
)

type MinerLocalConfig struct {
	MinerId    string
	MacAddress string
	WakeOnLan  bool
}

type LocalConfig struct {
	MaxSpotPriceForMining float64
	Miners                []MinerLocalConfig
	Hash                  []byte
}

func CreateLocalConfig() LocalConfig {
	config := LocalConfig{
		MaxSpotPriceForMining: 0.0,
		Hash:                  []byte{},
	}
	err := tools.ReadJsonFile("local_config.json", &config)
	if err != nil {
		log.Printf("failed to read local config file: %v\n", err)
	}
	config.Hash, err = tools.CalculateHash(config)
	if err != nil {
		log.Printf("failed to calculate hash for local config: %v\n", err)
	}
	return config
}

func (lc *LocalConfig) ReloadIfNeeded() bool {
	newConfig := CreateLocalConfig()
	if string(newConfig.Hash) != string(lc.Hash) {
		log.Println("Reloading local config")
		lc.MaxSpotPriceForMining = newConfig.MaxSpotPriceForMining
		lc.Hash = newConfig.Hash
		return true
	}
	return false
}

func (lc *LocalConfig) GetMinerConfig(minerId string) *MinerLocalConfig {
	for _, minerConfig := range lc.Miners {
		if minerConfig.MinerId == minerId {
			return &minerConfig
		}
	}
	return nil
}
