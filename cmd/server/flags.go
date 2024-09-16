package main

import (
	"errors"
	"flag"
	"metal/config"
	"strconv"
	"strings"
)

var cfg = config.GetConfigServer()
var (
	startAddress    Address
	storeInterval   int
	fileStoragePath string
	connectionURL   string
	restore         bool
)

// Можно доработать для большей разбивки на детали(протокол и т.д)
type Address struct {
	addr string
}

func (a Address) String() string {
	return a.addr
}

// Проверить что адрес передан в правильном формате
func (a *Address) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}
	_, err := strconv.Atoi(hp[1])
	if err != nil {
		return err
	}
	a.addr = s
	return nil
}

func parseFlags() {
	startAddress = Address{
		addr: "localhost:8080",
	}
	flag.Var(&startAddress, "a", "host and port which server will run on")
	flag.IntVar(&storeInterval, "interval", 300, "interval of saving data to disk, in seconds")
	flag.StringVar(&fileStoragePath, "p", "./save.json", "path where file will be stored")

	// ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	`localhost`, `postgres`, `yjdfz21f`, `metal`)

	flag.StringVar(&connectionURL, "d", "", "PostgreSQL connection url")
	flag.BoolVar(&restore, "r", true, "whether to restore data from save file or not")
	if cfg.Address != "" {
		startAddress.addr = cfg.Address
	}
	if cfg.ConnectionURL != "" {
		connectionURL = cfg.ConnectionURL
	}
	if cfg.StoreInterval != nil {
		storeInterval = *cfg.StoreInterval
	}
	if cfg.FileStoragePath != "" {
		fileStoragePath = cfg.FileStoragePath
	}
	if cfg.Restore != nil {
		restore = *cfg.Restore
	}
	flag.Parse()

}
