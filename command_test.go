package gerte

import "testing"

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
