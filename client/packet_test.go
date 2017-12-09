package client

import (
	"net"
	"reflect"
	"testing"

	"github.com/u-root/dhcp6"
)

var mac = net.HardwareAddr([]byte{0xb8, 0xae, 0xed, 0x7a, 0x10, 0x66})

func TestNewSolicitOptions(t *testing.T) {
	options, err := newSolicitOptions(mac)
	if err != nil {
		t.Fatalf("error in newSolicitOptions: %v\n", err)
	}
	expected := dhcp6.Options(map[dhcp6.OptionCode][][]byte{
		dhcp6.OptionIANA:        [][]byte{[]byte{0x72, 0x6f, 0x6f, 0x74, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		dhcp6.OptionRapidCommit: [][]byte{nil},
		dhcp6.OptionElapsedTime: [][]byte{[]byte{0x00, 0x00}},
		dhcp6.OptionORO:         [][]byte{[]byte{0x00, byte(dhcp6.OptionDNSServers), 0x00, byte(dhcp6.OptionDomainList), 0x00, byte(dhcp6.OptionBootFileURL), 0x00, byte(dhcp6.OptionBootFileParam)}},
		dhcp6.OptionClientID:    [][]byte{[]byte{0x00, 0x03, 0x00, 0x06, 0xb8, 0xae, 0xed, 0x7a, 0x10, 0x66}},
	})

	optionsIANA, err := options.IANA()
	if err != nil {
		t.Fatalf("getting IANA from options: got %v, want nil", err)
	}

	expectedIANA, err := expected.IANA()
	if err != nil {
		t.Fatalf("getting IANA from expected: got %v, want nil", err)
	}
	if !reflect.DeepEqual(optionsIANA, expectedIANA) {
		t.Fatalf(
			"incorrect newSolicitOptions: IANAs do not match, get %v, but should be %v instead\n",
			optionsIANA, expectedIANA,
		)
	}

	if err := options.RapidCommit(); err != nil {
		t.Fatalf("getting RapidCommit from options: got %v, want nil", err)
	}
	if err := expected.RapidCommit(); err != nil {
		t.Fatalf("getting RapidCommit from expected: got %v, want nil", err)
	}

	optionsOR, err := options.OptionRequest()
	if err != nil {
		t.Fatalf("getting OptionRequest from options: got %v, want nil", err)
	}
	expectedOR, err := expected.OptionRequest()
	if err != nil {
		t.Fatalf("getting OptionRequest from expected: got %v, want nil", err)
	}
	if !reflect.DeepEqual(optionsOR, expectedOR) {
		t.Fatalf(
			"incorrect newSolicitOptions: Option request do not match, get %v, but should be %v instead\n",
			optionsOR, expectedOR,
		)
	}

	optionsElapsedTime, err := options.ElapsedTime()
	if err != nil {
		t.Fatalf("getting ElapsedTime from options: got %v, want nil", err)
	}
	expectedElapsedTime, err := expected.ElapsedTime()
	if err != nil {
		t.Fatalf("getting ElapsedTime from expected: got %v, want nil", err)
	}
	if !reflect.DeepEqual(optionsElapsedTime, expectedElapsedTime) {
		t.Fatalf(
			"incorrect newSolicitOptions: Elapsed time do not match, get %v, but should be %v instead\n",
			optionsElapsedTime, expectedElapsedTime,
		)
	}

	optionsClientID, err := options.ClientID()
	if err != nil {
		t.Fatalf("getting ClientID from options: got %v, want nil", err)
	}
	expectedClientID, err := expected.ClientID()
	if err != nil {
		t.Fatalf("getting ClientID from expected: got %v, want nil", err)
	}
	if !reflect.DeepEqual(optionsClientID, expectedClientID) {
		t.Fatalf(
			"incorrect newSolicitOptions: Client IDs do not match, get %v, but should be %v instead\n",
			optionsClientID, expectedClientID,
		)
	}

	if !reflect.DeepEqual(expected, options) {
		t.Fatalf("incorrect newSolicitOptions: extra unnecessary options\n%v\n%v\n", expected, options)
	}
}

func TestNewSolicitPacket(t *testing.T) {
	p, err := newSolicitPacket(mac)
	if err != nil {
		t.Fatalf("error in newSolicitPacket: %v\n", err)
	}

	options, err := newSolicitOptions(mac)
	expected := &dhcp6.Packet{
		MessageType:   dhcp6.MessageTypeSolicit,
		TransactionID: [3]byte{0x00, 0x01, 0x02},
		Options:       options,
	}
	if !reflect.DeepEqual(p, expected) {
		t.Fatalf("incorrect newSolicitPacket: get %v but should be %v\n", p, expected)
	}
}
