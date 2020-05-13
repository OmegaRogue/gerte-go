package gerte

import (
	"strings"
	"testing"
)

func TestVersionFromToBytes(t *testing.T) {
	version := Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
	}
	ver := version.ToBytes()
	version2 := VersionFromBytes(ver)

	if version != version2 {
		t.Error("versions don't match")
	}
}

func TestAddressFromToBytes(t *testing.T) {
	address := GertAddress{
		Upper: 123,
		Lower: 456,
	}
	addr := address.ToBytes()
	address2 := AddressFromBytes(addr)

	if address != address2 {
		t.Error("addresses don't match")
	}
}

func TestGertCFromToBytes(t *testing.T) {
	address := GERTc{
		GERTe: GertAddress{
			Upper: 0,
			Lower: 0,
		},
		GERTi: GertAddress{
			Upper: 123,
			Lower: 456,
		},
	}

	addr := address.ToBytes()
	address2 := GertCFromBytes(addr)

	if address != address2 {
		t.Error("addresses don't match")
	}
}

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

func TestCommandFromBytes(t *testing.T) {
	cmd := []byte{byte(CommandState), byte(StateConnected)}
	command := Command{
		Command: CommandState,
		Packet:  Packet{},
		Status: Status{
			Status: StateConnected,
			Size:   4,
			Version: Version{
				Major: 1,
				Minor: 1,
				Patch: 0,
			},
		},
	}
	verByte := Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
	}.ToBytes()
	cmd = append(cmd, verByte[0], verByte[1])
	command2, err := CommandFromBytes(cmd)
	if err != nil {
		t.Errorf("error on unmarshal command: %+v", err)
	}
	if command.Command != command2.Command ||
		string(command.Packet.Data) != string(command2.Packet.Data) ||
		command.Packet.Target != command2.Packet.Target ||
		command.Packet.Source != command2.Packet.Source ||
		command.Status != command2.Status {
		t.Errorf("commands don't match:\n%+v\n%+v", command, command2)
	}

}

func TestStatusFromBytes(t *testing.T) {
	st := []byte{byte(StateConnected)}
	state := Status{
		Status: StateConnected,
		Size:   4,
		Version: Version{
			Major: 1,
			Minor: 1,
			Patch: 0,
		},
	}
	verByte := Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
	}.ToBytes()
	st = append(st, verByte[0], verByte[1])
	state2, err := StatusFromBytes(st)
	if err != nil {
		t.Errorf("error on unmarshal command: %+v", err)
	}
	if state != state2 {
		t.Errorf("commands don't match:\n%+v\n%+v", state, state2)
	}
}
