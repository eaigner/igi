package node

import (
	"log"
	"net"
)

type Udp struct {
	host   string
	logger *log.Logger
	conn   *net.UDPConn
	done   chan bool
}

func NewUdp(host string, logger *log.Logger) *Udp {
	return &Udp{host: host, logger: logger, done: make(chan bool, 1)}
}

func (udp *Udp) Serve() error {
	addr, err := net.ResolveUDPAddr("udp", udp.host)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	udp.conn = conn
	udp.logger.Printf("listening on udp://%v", addr)

	go run(udp, conn)

	return nil
}

func (udp *Udp) Shutdown() error {
	if udp.conn != nil {
		if err := udp.conn.Close(); err != nil {
			return err
		}
		udp.conn = nil
	}
	return nil
}

func run(udp *Udp, conn *net.UDPConn) {
	var buf [2048]byte
	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			udp.logger.Printf("error reading UDP packet: %v", err)
			break
		} else {
			udp.logger.Printf("read UDP packet: len=%v, %v", n, buf[:n])
		}
	}
	udp.logger.Printf("udp server closed")
	udp.done <- true
}
