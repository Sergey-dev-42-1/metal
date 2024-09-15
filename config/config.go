package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
)

type ConfigAgent struct {
	Address        string `env:"ADDRESS"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

type ConfigServer struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   *int   `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         *bool  `env:"RESTORE"`
}

func GetConfigAgent() *ConfigAgent {
	var cfg ConfigAgent
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg.Address, cfg.PollInterval, cfg.ReportInterval)
	return &cfg
}

func GetConfigServer() *ConfigServer {
	var cfg ConfigServer
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg.Address, cfg.Restore, cfg.StoreInterval, cfg.FileStoragePath)
	return &cfg
}
