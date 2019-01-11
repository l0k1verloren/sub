package sub

import (
	"net"
)

// BaseInterface is the core functions required for a Base
type BaseInterface interface {
	SetupListener()
}

// BaseCfg is the configuration for a Base
type BaseCfg struct {
	Listener   string
	Password   []byte
	BufferSize int
}

// Base is the common structure between a worker and a node
type Base struct {
	cfg      BaseCfg
	listener *net.UDPConn
	packets  chan *Packet
	kill     chan bool
	Messages chan *Message
}

// A Node is a server with some number of subscribers
type Node struct {
	Base
	subscribers []*net.UDPAddr
}

// A Worker is a node that subscribes to a Node's messages
type Worker struct {
	Base
	node *net.UDPAddr
}

// Packet is the structure of individual encoded packets of the message. These are made from a 9/3 Reed Solomon code and 9 are sent in distinct packets and only 3 are required to guarantee retransmit-free delivery.
// If the CRC is incorrect, the message reconstructor will omit this block from the reconstruction process.
type Packet struct {
	bytes  []byte       // raw FEC encoded bytes of packet
	check  uint32       // CRC32 checksum of bytes to quickly identify corrupt packets
	sender *net.UDPAddr // address packet was received from
}

// A Bundle is a collection of the received packets received from the same sender with up to 9 pieces.
type Bundle struct {
	packets []Packet
}

// Message is the data reconstructed from a complete Bundle, containing data in messagepack format
type Message struct {
	sender    string
	recipient string
	timestamp uint64
	numBytes  uint16
	bytes     []byte // messagepack payload
}

// Subscription is the message sent by a worker node to request updates from the node
type Subscription struct {
	address string
	pubKey  []byte
}

// Confirmation is the reply message for a subscription request
type Confirmation struct {
	subscriber string // confirming address of subscriber
	pubKey     []byte // public key of server for message verification
}

const (
	uNet              = "udp4"
	defaultBufferSize = 16384
	maxMessageSize    = 3072
)
