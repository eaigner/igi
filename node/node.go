package node

import "log"

type Node struct {
	conf   Conf
	logger *log.Logger
	udp    *Udp
}

func New(conf Conf, logger *log.Logger) *Node {
	return &Node{
		conf:   conf,
		logger: logger,
		udp:    NewUdp(conf.UdpHost, logger),
	}
}

func (node *Node) Serve() error {
	if err := node.udp.Serve(); err != nil {
		return err
	}
	return nil
}

func (node *Node) Shutdown() error {
	if err := node.udp.Shutdown(); err != nil {
		return err
	}

	<-node.udp.done

	return nil
}
