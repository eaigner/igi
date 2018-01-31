package node

import "github.com/eaigner/igi/storage"

type Node struct {
	conf   Conf
	logger Logger
	store  storage.Store
	udp    *UDP
}

func New(conf Conf, store storage.Store, logger Logger) *Node {
	return &Node{
		conf:   conf,
		logger: logger,
		store:  store,
		udp:    NewUDP(conf.UdpHost, conf.MinWeightMagnitude, logger, store),
	}
}

func (node *Node) Serve() error {
	if err := node.udp.Listen(); err != nil {
		return err
	}
	return nil
}

func (node *Node) Shutdown() error {
	node.udp.Close()
	return node.store.Close()
}
