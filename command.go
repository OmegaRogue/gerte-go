package gerte

import "fmt"

type (
	// GertCommand indicates the Command of a Request
	GertCommand byte

	// Command is the Parsed Data returned from Api.Parse
	Command struct {
		Command GertCommand
		Packet  Packet
		Status  Status
	}
)

const (
	// CommandState is used by relays to indicate result of a command to a gateway and by gateways to request the state.
	CommandState GertCommand = iota
	// CommandRegister claims an address for this gateway using a key.
	CommandRegister
	// CommandData transmits data from a GERTi address to a GERTc address.
	CommandData
	// CommandClose gracefully closes a connection to a relay.
	CommandClose
)

// CommandFromBytes parses bytes to a GERT Command
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

// GoString prints a GERT Command to a Human-readable string
func (cmd Command) String() string {
	switch cmd.Command {
	case CommandState:
		return fmt.Sprintf("%v %v", cmd.Command, cmd.Status)
	case CommandClose:
		return fmt.Sprintf("%v", cmd.Command)
	case CommandData:
		return fmt.Sprintf("%v %v", cmd.Command, cmd.Packet)
	}
	return "nil"
}

// GoString prints a GERT Command to a Human-readable string surrounded with brackets
func (cmd Command) GoString() string {
	switch cmd.Command {
	case CommandState:
		return fmt.Sprintf("%#v%#v", cmd.Command, cmd.Status)
	case CommandClose:
		return fmt.Sprintf("%#v", cmd.Command)
	case CommandData:
		return fmt.Sprintf("%#v%#v", cmd.Command, cmd.Packet)
	}
	return "[nil]"
}

// GoString prints a GERT Command to a Human-readable string
func (cmd GertCommand) String() string {
	switch cmd {
	case CommandState:
		return "STATE"
	case CommandClose:
		return "CLOSE"
	case CommandData:
		return "DATA"
	case CommandRegister:
		return "REGISTER"
	}

	return "nil"
}

// GoString prints a GERT Command to a Human-readable string surrounded with brackets
func (cmd GertCommand) GoString() string {
	return fmt.Sprintf("[%v]", cmd)
}
