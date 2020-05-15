package gerte

import (
	"fmt"
)

// Packet is the Parsed Data returned from Api.Parse for a Data Command
type Packet struct {
	Source GERTc
	Target GERTc
	Data   []byte
}

// PacketFromBytes parses bytes to a GERT Packet
func PacketFromBytes(data []byte) (Packet, error) {
	source := GertCFromBytes(data[:6])
	target := GertCFromBytes(data[6:12])

	return Packet{
		Source: source,
		Target: target,
		Data:   data[13:],
	}, nil
}

// ToBytes converts a Packet to bytes for sending
func (pkt Packet) ToBytes() ([]byte, error) {
	if len(pkt.Data) > 255 {
		return nil, fmt.Errorf("data cannot exceed 255 bytes")
	}
	addressPart := append(pkt.Target.ToBytes(), pkt.Source.GERTi.ToBytes()...)
	dataPart := append([]byte{byte(len(pkt.Data))}, pkt.Data...)
	return append(addressPart, dataPart...), nil
}

// String prints a GERT Packet to a Human-readable string
func (pkt Packet) String() string {
	return fmt.Sprintf("%v %v %v: %v", pkt.Source, pkt.Target, len(pkt.Data), string(pkt.Data))
}

// GoString prints a GERT Packet to a Human-readable string surrounded with brackets
func (pkt Packet) GoString() string {
	return fmt.Sprintf("[%v][%v][%v][%v]", pkt.Source, pkt.Target, len(pkt.Data), string(pkt.Data))
}
