package node

type Node struct {
	conf   Conf
	logger Logger
	udp    *UDP
}

func New(conf Conf, logger Logger) *Node {
	return &Node{
		conf:   conf,
		logger: logger,
		udp:    NewUDP(conf.UdpHost, conf.MinWeightMagnitude, logger),
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
