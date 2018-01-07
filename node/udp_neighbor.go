package node

import (
	"log"
	"net"

	"github.com/eaigner/igi/trinary"
)

type UDPNeighbor struct {
	host   string
	done   chan bool
	logger *log.Logger
	conn   *net.UDPConn
}

func NewUDPNeighbor(host string, logger *log.Logger) *UDPNeighbor {
	return &UDPNeighbor{
		host:   host,
		done:   make(chan bool, 1),
		logger: logger,
	}
}

func (udp *UDPNeighbor) Listen() error {
	addr, err := net.ResolveUDPAddr("udp", udp.host)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}

	udp.logger.Printf("listening on udp://%v", addr)
	udp.conn = conn

	go udp.read(conn)

	return nil
}

func (udp *UDPNeighbor) Close() {
	udp.conn.Close()
	<-udp.done
}

func (udp *UDPNeighbor) read(conn *net.UDPConn) {
	var buf [1024 * 10]byte
	var t trinary.Trits

	for {
		n, err := conn.Read(buf[:])
		if err != nil {
			udp.logger.Printf("error reading UDP packet: %v", err)
			break
		} else {
			// TODO: remove log and handle trytes
			if trinary.BytesToTrits(buf[:n], &t) > 0 {
				udp.logger.Printf("read trytes UDP packet: len=%v, %s", t.Len(), t.Trytes())
			} else {
				udp.logger.Printf("read raw UDP packet: len=%v, %x", n, buf[:n])
			}
		}
	}
	udp.logger.Printf("udp server closed")
	udp.done <- true
}
