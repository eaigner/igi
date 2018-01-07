package main

import (
	"flag"

	"github.com/eaigner/igi/node"
)

var conf node.Conf

func init() {
	flag.IntVar(&conf.Port, "p", 14600, "server port")
	flag.IntVar(&conf.UdpPort, "t", 14600, "udp receiver port")
	flag.IntVar(&conf.TcpPort, "u", 15600, "tcp receiver port")
	flag.BoolVar(&conf.Debug, "debug", false, "turn on debug mode")
	flag.BoolVar(&conf.Testnet, "testnet", false, "use testnet")
	flag.Var(&conf.Neighbors, "n", "single neighbor node URL, flag can be used multiple times")
}

func main() {
	flag.Parse()
}
