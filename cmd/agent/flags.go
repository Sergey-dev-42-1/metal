package main

import (
	"flag"
)

var (
	startAddress   string
	reportInterval int
	pollInterval   int
)

func parseFlags() {
	flag.StringVar(&startAddress, "a", "localhost:8080", "host and port which data will be sent to")
	flag.IntVar(&reportInterval, "r", 10, "set frequency of sending data to server in seconds")
	flag.IntVar(&pollInterval, "p", 2, "set frequency of getting runtime metrics")
	flag.Parse()
}
