package opts

import (
	"io"
	"net"

	"github.com/mdlayher/dhcp6"
	"github.com/mdlayher/dhcp6/internal/buffer"
)

// RelayMessage represents a raw RelayMessage generated by DHCPv6 relay agent, using RFC 3315,
// Section 7.
type RelayMessage struct {
	// RELAY-FORW or RELAY-REPL only
	MessageType dhcp6.MessageType

	// Number of relay agents that have relayed this
	// message.
	HopCount uint8

	// A global or site-local address that will be used by
	// the server to identify the link on which the client
	// is located.
	LinkAddress net.IP

	// The address of the client or relay agent from which
	// the message to be relayed was received.
	PeerAddress net.IP

	// Options specifies a map of DHCP options.  Its methods can be used to
	// retrieve data from an incoming RelayMessage, or send data with an outgoing
	// RelayMessage.
	// MUST include a "Relay Message option" (see
	// section 22.10); MAY include other options added by
	// the relay agent.
	Options dhcp6.Options
}

// MarshalBinary allocates a byte slice containing the data
// from a RelayMessage.
func (rm *RelayMessage) MarshalBinary() ([]byte, error) {
	// 1 byte: message type
	// 1 byte: hop-count
	// 16 bytes: link-address
	// 16 bytes: peer-address
	// N bytes: options slice byte count
	b := buffer.New(nil)

	b.Write8(uint8(rm.MessageType))
	b.Write8(rm.HopCount)
	copy(b.WriteN(net.IPv6len), rm.LinkAddress)
	copy(b.WriteN(net.IPv6len), rm.PeerAddress)
	rm.Options.Marshal(b)

	return b.Data(), nil
}

// UnmarshalBinary unmarshals a raw byte slice into a RelayMessage.
//
// If the byte slice does not contain enough data to form a valid RelayMessage,
// ErrInvalidPacket is returned.
func (rm *RelayMessage) UnmarshalBinary(p []byte) error {
	b := buffer.New(p)
	// RelayMessage must contain at least message type, hop-count, link-address and peer-address
	if b.Len() < 34 {
		return io.ErrUnexpectedEOF
	}

	rm.MessageType = dhcp6.MessageType(b.Read8())
	rm.HopCount = b.Read8()

	rm.LinkAddress = make(net.IP, net.IPv6len)
	copy(rm.LinkAddress, b.Consume(net.IPv6len))

	rm.PeerAddress = make(net.IP, net.IPv6len)
	copy(rm.PeerAddress, b.Consume(net.IPv6len))

	if err := (&rm.Options).Unmarshal(b); err != nil {
		// Invalid options means an invalid RelayMessage
		return dhcp6.ErrInvalidPacket
	}
	return nil
}
