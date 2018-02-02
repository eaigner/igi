package main

import (
	"flag"
	"github.com/eaigner/igi/storage"
	"log"
	"os"
	"os/signal"
	"syscall"

	gonode "github.com/eaigner/igi/node"
)

var conf gonode.Conf

func init() {
	flag.StringVar(&conf.HttpHost, "p", ":15100", "http server address")
	flag.StringVar(&conf.UdpHost, "u", ":15200", "udp socket address")
	flag.StringVar(&conf.TcpHost, "t", ":15300", "tcp socket address")
	flag.StringVar(&conf.DbPath, "db", "tangle.db", "tangle database path")
	flag.BoolVar(&conf.Debug, "debug", false, "turn on debug mode")
	flag.BoolVar(&conf.Testnet, "testnet", false, "use testnet")
	flag.Var(&conf.Neighbors, "n", "single neighbor node URL, flag can be used multiple times")
	flag.IntVar(&conf.MinWeightMagnitude, "w", 14, "min weight magnitude")
	flag.Parse()
}

func main() {

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	logger := gonode.NewNullLogger()

	if conf.Debug {
		logger = log.New(os.Stdout, "igi: ", 0)
	}

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logger.Println(sig)
		done <- true
	}()

	db, err := storage.NewBoltStore(conf.DbPath)

	if err != nil {
		panic(err)
	}

	node := gonode.New(conf, db, logger)

	logger.Println("starting node...")

	if err := node.Serve(); err != nil {
		logger.Printf("error starting node: %v", err)
		return
	}

	logger.Println("node started")

	<-done

	if err := node.Shutdown(); err != nil {
		logger.Printf("error shutting down node: %v", err)
	}

	logger.Println("node stopped")
}
