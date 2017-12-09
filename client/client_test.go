package client

import (
	"net"
	"reflect"
	"testing"

	"github.com/u-root/dhcp6"
)

const (
	v6addr = "fe80::baae:edff:fe79:6191"
)

func TestSolicitAndAdvertise(t *testing.T) {
	p, _ := newSolicitPacket(mac)
	pb, _ := p.MarshalBinary()

	r := &testMessage{
		addr: &net.UDPAddr{
			IP: net.ParseIP(v6addr),
		},
	}
	r.b.Write(pb)

	reply, err := serve(r)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if reply == nil {
		t.Fatalf("The reply is nil?")
	}

	if reply.MessageType != dhcp6.MessageTypeAdvertise {
		t.Fatalf("Reply does not have a correct typeShould be MessageTypeAdvertise but %s instead\n", reply.MessageType)
	}
	if !reflect.DeepEqual(reply.TransactionID, [3]byte{0x00, 0x01, 0x02}) {
		t.Fatalf("Reply txID does not match")
	}

	iana, err := reply.Options.IANA()
	if err != nil {
		t.Fatalf("Reply does not contain valid IANA: %v", err)
	}

	iaaddr, err := iana[0].Options.IAAddr()
	if err != nil {
		t.Fatalf("Reply does not contain valid IAAddr: %v", err)
	}
	t.Logf("Get assigned ipv6 addr from server: %+v", iaaddr[0])
}
