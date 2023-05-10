package config

import "flag"

var (
	Addr    string
	ResAddr string
)

func ParseFlags() {
	flag.StringVar(&Addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&ResAddr, "b", "http://localhost:8080", "")

	flag.Parse()
}
