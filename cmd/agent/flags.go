package main

import (
	"flag"
	"metal/config"
)

var (
	startAddress   string
	reportInterval int
	pollInterval   int
)

var cfg = config.GetConfigAgent()

func parseFlags() {

	flag.StringVar(&startAddress, "a", "localhost:8080", "host and port which data will be sent to")
	flag.IntVar(&reportInterval, "r", 10, "set frequency of sending data to server in seconds")
	flag.IntVar(&pollInterval, "p", 2, "set frequency of getting runtime metrics")
	
	if cfg.Address != "" {
		startAddress = cfg.Address
	}
	if cfg.ReportInterval != 0 {
		reportInterval = cfg.ReportInterval
	}
	if cfg.PollInterval != 0 {
		pollInterval = cfg.PollInterval
	}

	flag.Parse()
}
