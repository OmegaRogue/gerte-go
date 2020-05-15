package gerte

import (
	"strings"
	"testing"
)

func TestPacketFromToBytes(t *testing.T) {
	packet := Packet{
		Source: GERTc{
			GERTe: GertAddress{},
			GERTi: GertAddress{
				Upper: 123,
				Lower: 456,
			},
		},
		Target: GERTc{
			GERTe: GertAddress{},
			GERTi: GertAddress{
				Upper: 123,
				Lower: 456,
			},
		},
		Data: []byte("test"),
	}
	pkt, err := packet.ToBytes()
	if err != nil {
		t.Errorf("error on marshal packet: %+v", err)
	}

	var b strings.Builder
	b.Write(pkt[:6])
	b.Write([]byte{0, 0, 0})
	b.Write(pkt[6:])

	pkt = []byte(b.String())

	packet2, err := PacketFromBytes(pkt)
	if err != nil {
		t.Errorf("error on unmarshal packet: %+v", err)
	}

	if string(packet.Data) != string(packet2.Data) ||
		packet.Target != packet2.Target ||
		packet.Source != packet2.Source {
		t.Errorf("packets don't match:\n%+v\n%+v", packet, packet2)
	}
}
