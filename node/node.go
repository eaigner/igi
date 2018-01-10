package node

import "log"

type Node struct {
	conf   Conf
	logger *log.Logger
	udp    *UDPNeighbor
}

func New(conf Conf, logger *log.Logger) *Node {
	return &Node{
		conf:   conf,
		logger: logger,
		udp:    NewUDPNeighbor(conf.UdpHost, conf.MinWeightMagnitude, logger),
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
	return nil
}
