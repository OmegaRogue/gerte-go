package gerte

import (
	"fmt"
	"strings"
)

// TODO Make Doc
func (ver Version) ToBytes() []byte {
	return []byte{ver.Major, ver.Minor}
}

// TODO Make Doc
func VersionFromBytes(b []byte) Version {
	return Version{
		Major: b[0],
		Minor: b[1],
	}
}

// TODO Make Doc
func (addr GertAddress) ToBytes() []byte {
	var b strings.Builder
	b.WriteByte(byte(addr.Upper >> 4))
	b.WriteByte(byte(((addr.Upper & 0x0F) << 4) | (addr.Lower >> 8)))
	b.WriteByte(byte(addr.Lower & 0xFF))
	return []byte(b.String())
}

// TODO Make Doc
func AddressFromBytes(data []byte) GertAddress {
	return GertAddress{
		Upper: (int(data[0]) << 4) | (int(data[1]) >> 4),
		Lower: ((int(data[1]) & 0x0F) << 8) | int(data[2]),
	}
}

// TODO Make Doc
func GertCFromBytes(data []byte) GERTc {
	return GERTc{
		GERTe: GertAddress{
			Upper: (int(data[0]) << 4) | (int(data[1]) >> 4),
			Lower: ((int(data[1]) & 0x0F) << 8) | int(data[2]),
		},
		GERTi: GertAddress{
			Upper: (int(data[3]) << 4) | (int(data[4]) >> 4),
			Lower: ((int(data[4]) & 0x0F) << 8) | int(data[5]),
		},
	}

}

// TODO Make Doc
func (addr GERTc) ToBytes() []byte {
	return append(addr.GERTe.ToBytes(), addr.GERTi.ToBytes()...)
}

// TODO Make Doc
func PacketFromBytes(data []byte) (Packet, error) {
	source := GertCFromBytes(data[:6])
	target := GertCFromBytes(data[6:12])

	return Packet{
		Source: source,
		Target: target,
		Data:   data[13:],
	}, nil
}

// TODO Make Doc
func (pkt Packet) ToBytes() ([]byte, error) {
	if len(pkt.Data) > 255 {
		return nil, fmt.Errorf("data cannot exceed 255 bytes")
	}
	addressPart := append(pkt.Target.ToBytes(), pkt.Source.GERTi.ToBytes()...)
	dataPart := append([]byte{byte(len(pkt.Data))}, pkt.Data...)
	return append(addressPart, dataPart...), nil
}

// TODO Make Doc
func StatusFromBytes(data []byte) (Status, error) {

	switch data[0] {
	case byte(StateFailure):
		return Status{
			Status: StateFailure,
			Size:   2,
			Error:  GertError(data[1]),
		}, nil
	case byte(StateConnected):
		if len(data) < 3 {
			return Status{}, fmt.Errorf("data too short: %v<3", len(data))
		}
		return Status{
			Status: StateConnected,
			Size:   4,
			Version: Version{
				Major: data[1],
				Minor: data[2],
			},
		}, nil
	case byte(StateAssigned):
		return Status{
			Status: StateAssigned,
			Size:   1,
		}, nil
	case byte(StateClosed):
		return Status{
			Status: StateClosed,
			Size:   1,
		}, nil
	case byte(StateSent):
		return Status{
			Status: StateSent,
			Size:   1,
		}, nil
	}
	return Status{}, fmt.Errorf("state didn't match any known state: %v", data[0])
}

// TODO Make Doc
func CommandFromBytes(data []byte) (Command, error) {
	switch data[0] {
	case byte(CommandState):
		state, err := StatusFromBytes(data[1:])
		if err != nil {
			return Command{}, fmt.Errorf("error while parsing status data: %w", err)
		}
		return Command{
			Command: CommandState,
			Status:  state,
		}, nil
	case byte(CommandRegister):
		return Command{
			Command: CommandRegister,
		}, nil
	case byte(CommandData):
		packet, err := PacketFromBytes(data[1:])
		if err != nil {
			return Command{}, fmt.Errorf("error while parsing packet data: %w", err)
		}
		return Command{
			Command: CommandData,
			Packet:  packet,
		}, nil
	case byte(CommandClose):
		return Command{
			Command: CommandClose,
		}, nil
	}
	return Command{}, fmt.Errorf("error while parsing command data: invalid command")
}
