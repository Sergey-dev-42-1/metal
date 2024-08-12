package main

import (
	"errors"
	"flag"
	"strconv"
	"strings"
)

var startAddress Address

// Можно доработать для больше разбивки на детали(протокол и т.д)
type Address struct {
	addr string
}

func (a Address) String() string {
	return a.addr
}

// Проверить что адрес передан в правильно формате
func (a *Address) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("Need address in a form host:port")
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
	flag.Parse()
}
