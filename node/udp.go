package node

import (
	"github.com/eaigner/igi/queue"
	"github.com/eaigner/igi/storage"
	"net"

	"github.com/eaigner/igi/hash"
	"github.com/eaigner/igi/trinary"
)

type UDP struct {
	host         string
	done         chan bool
	minWeightMag int
	logger       Logger
	store        storage.Store
	conn         *net.UDPConn
	txCache      *Cache
	receiveQueue *queue.WeightQueue
	replyQueue   *queue.WeightQueue
	closed       bool
}

func NewUDP(host string, minWeightMag int, logger Logger, store storage.Store) *UDP {
	return &UDP{
		host:         host,
		minWeightMag: minWeightMag,
		done:         make(chan bool, 1),
		logger:       logger,
		store:        store,
		txCache:      NewCache(1024),
		receiveQueue: queue.NewWeightQueue(1024),
		replyQueue:   queue.NewWeightQueue(1024),
		closed:       false,
	}
}

func (udp *UDP) Listen() error {
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

	go udp.replyLoop()
	go udp.receiveLoop()
	go udp.read(conn)

	return nil
}

func (udp *UDP) Close() {
	udp.conn.Close()
	udp.closed = true
	<-udp.done
}

func (udp *UDP) read(conn *net.UDPConn) {
	var buf [1024 * 10]byte
	for {
		n, addr, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			udp.logger.Printf("error reading UDP packet: %v", err)
			break
		} else {
			udp.handleMessage(buf[:n], addr)
		}
	}
	udp.logger.Printf("udp server closed")
	udp.done <- true
}

func (udp *UDP) replyLoop() {
	for !udp.closed {
		var _ = udp.replyQueue.Pop().(*replyItem)
		// TODO: do something with item, implement "Node.replyToRequest"
	}
}

func (udp *UDP) receiveLoop() {
	for !udp.closed {
		item := udp.receiveQueue.Pop().(*receiveItem)

		// TODO: do something with item, implement "Node.processReceivedData"
		if err := item.msg.Store(udp.store); err != nil {
			// TODO: handle error
			udp.logger.Printf("message not stored: %v", err)
		} else {
			udp.logger.Printf("message stored %x", item.msg.digest)
			// TODO: was stored, broadcast
		}
	}
}

func (udp *UDP) handleMessage(b []byte, neighbor *net.UDPAddr) {
	udp.logger.Printf("message from UDP neighbor: %v", neighbor)

	msg, err := ParseUdpBytes(b)
	if err != nil {
		udp.logger.Printf("error parsing message: %v", err)
		return // drop
	}

	// Check if we have seen this transaction lately.
	// []uint8 is not a valid cache key, so we use the hex digest.
	key := msg.TxDigestHex()

	_, cached := udp.txCache.Get(key)

	if !cached {
		if err := msg.Validate(udp.minWeightMag); err != nil {
			udp.logger.Printf("invalid message: %v", err)
			return // drop
		}
		udp.txCache.Add(key, msg.TxHash())
		udp.receiveQueue.Push(&receiveItem{msg, neighbor}, hash.WeightMagnitude(msg.TxHash()))
	}

	// Check if the trailer hash is the same as the current message transaction hash.
	// If it's the same, request a random tip by sending the zero hash.
	requestedHash := msg.TrailerHash()

	if trinary.Equals(msg.TxHash(), requestedHash) {
		requestedHash = make([]int8, len(requestedHash))
	}

	udp.replyQueue.Push(&replyItem{requestedHash, neighbor}, hash.WeightMagnitude(requestedHash))
}

type receiveItem struct {
	msg      *Message
	neighbor *net.UDPAddr
}

type replyItem struct {
	requestedHash []int8
	neighbor      *net.UDPAddr
}
