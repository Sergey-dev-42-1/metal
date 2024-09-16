package main

import (
	"flag"
	"metal/config"
)

var (
	startAddress   string
	reportInterval int
	pollInterval   int
	batching   bool
)

var cfg = config.GetConfigAgent()

func parseFlags() {

	flag.StringVar(&startAddress, "a", "localhost:8080", "host and port which data will be sent to")
	flag.IntVar(&reportInterval, "r", 10, "set frequency of sending data to server in seconds")
	flag.IntVar(&pollInterval, "p", 2, "set frequency of getting runtime metrics")
	flag.BoolVar(&batching, "b", false, "turn on / off batching, off by default")

	if cfg.Address != "" {
		startAddress = cfg.Address
	}
	if cfg.ReportInterval != 0 {
		reportInterval = cfg.ReportInterval
	}
	if cfg.PollInterval != 0 {
		pollInterval = cfg.PollInterval
	}
	if cfg.Batching != nil {
		batching = *cfg.Batching
	}

	flag.Parse()
}
