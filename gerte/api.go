package gerte

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type GertStatus byte
type GertCommand byte
type GertError byte

const (
	CommandState    GertCommand = 0
	CommandRegister GertCommand = 1
	CommandData     GertCommand = 2
	CommandClose    GertCommand = 3

	StateFailure   GertStatus = 0
	StateConnected GertStatus = 1
	StateAssigned  GertStatus = 2
	StateClosed    GertStatus = 3
	StateSent      GertStatus = 4

	ErrorVersion           GertError = 0
	ErrorBadKey            GertError = 1
	ErrorAlreadyRegistered GertError = 2
	ErrorNotRegistered     GertError = 3
	ErrorNoRoute           GertError = 4
	ErrorAddressTaken      GertError = 5
)

type GertAddress struct {
	Upper int
	Lower int
}
type GERTc struct {
	GERTe GertAddress
	GERTi GertAddress
}

type ApiConfig struct {
}
type Command struct {
	Command GertCommand
	Packet  Packet
	Status  Status
}

type Api struct {
	socket     net.Conn
	listener   net.Listener
	Registered bool
	Address    GertAddress
}

type Packet struct {
	Source GERTc
	Target GERTc
	Data   []byte
}
type Status struct {
	Status  GertStatus
	Size    int
	Error   GertError
	Version Version
}
type Version struct {
	Major int
	Minor int
	Patch int
}

func (status Status) ParseError() error {
	switch status.Error {
	case ErrorVersion:
		return fmt.Errorf("incompatible version during negotiation: %v", status.Version)
	case ErrorBadKey:
		return fmt.Errorf("key did not match that used for the requested address. Requested address may not exist")
	case ErrorAlreadyRegistered:
		return fmt.Errorf("registration has already been performed successfully")
	case ErrorNotRegistered:
		return fmt.Errorf("gateway cannot send data before claiming an address")
	case ErrorNoRoute:
		return fmt.Errorf("data failed to send because remote gateway could not be found")
	case ErrorAddressTaken:
		return fmt.Errorf("address request has already been claimed")
	}
	return fmt.Errorf("no valid error")
}

func FromString(addr string) (GertAddress, error) {
	parts := strings.Split(addr, ".")

	upper, err := strconv.ParseInt(parts[0], 10, 0)
	if err != nil {
		return GertAddress{}, fmt.Errorf("error on parse upper String: %+v", err)
	}
	lower, err := strconv.ParseInt(parts[1], 10, 0)
	if err != nil {
		return GertAddress{}, fmt.Errorf("error on parse lower String: %+v", err)
	}
	return GertAddress{
		Upper: int(upper),
		Lower: int(lower),
	}, nil
}

func (addr GertAddress) Conv() []byte {
	var b strings.Builder
	b.WriteByte(byte(addr.Upper >> 4))
	b.WriteByte(byte(((addr.Upper & 0x0F) << 4) | (addr.Lower >> 8)))
	b.WriteByte(byte(addr.Lower & 0xFF))
	return []byte(b.String())
}

func Parse(data []byte) GertAddress {
	return GertAddress{
		Upper: (int(data[0]) << 4) | (int(data[1]) >> 4),
		Lower: ((int(data[1]) & 0x0F) << 8) | int(data[2]),
	}
}

func ParseFull(data []byte) GERTc {
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

func (ver Version) PrintVersion() string {
	return fmt.Sprintf("%v.%v.%v", ver.Major, ver.Minor, ver.Patch)
}
func (addr GertAddress) PrintAddress() string {
	return fmt.Sprintf("%v.%v", addr.Upper, addr.Lower)
}
func (addr GERTc) PrintGERTc() string {
	return fmt.Sprintf("%v.%v:%v.%v", addr.GERTe.Upper, addr.GERTe.Lower, addr.GERTi.Upper, addr.GERTi.Lower)
}
func makePacket(data []byte) (Packet, error) {
	if len(data) < 13 {
		return Packet{}, fmt.Errorf("data too short: %v<13", len(data))
	}

	source := ParseFull(data[:6])
	target := ParseFull(data[6:12])
	length := int(data[12])

	if len(data) < 13+length {
		return Packet{}, fmt.Errorf("data too short: %v<13+%v", len(data), length)
	}

	return Packet{
		Source: source,
		Target: target,
		Data:   data[13:],
	}, nil
}

func makeStatus(data []byte) (Status, error) {
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
				Major: int(data[1]),
				Minor: int(data[2]),
				Patch: int(data[3]),
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

func (api *Api) Startup(c net.Conn) error {
	if api.socket != nil {
		return nil
	}

	api.socket = c
	return nil
}

func (api *Api) Register(addr GertAddress, key string) (bool, error) {
	var b strings.Builder
	b.WriteByte(byte(CommandRegister))
	rawAddr := addr.Conv()
	b.Write(rawAddr)
	b.WriteString(key)
	_, err := api.socket.Write([]byte(b.String()))
	if err != nil {
		return false, fmt.Errorf("error on write: %+v", err)
	}

	cmd, err := api.Parse()
	if err != nil {
		return false, fmt.Errorf("error parsing response: %+v", err)
	}
	if cmd.Command == CommandState {
		switch cmd.Status.Status {
		case StateFailure:
			return false, cmd.Status.ParseError()
		case StateAssigned:
			api.Address = addr
			return true, nil
		}
	}
	return false, fmt.Errorf("no valid response")
}

func (api *Api) Transmit(target GERTc, source GertAddress, data []byte) (bool, error) {

	if api.socket == nil {
		return false, fmt.Errorf("not connected")
	}
	if len(data) > 255 {
		return false, fmt.Errorf("data cannot exceed 255 bytes")
	}

	var b strings.Builder

	b.WriteByte(byte(CommandData))
	b.Write(target.GERTe.Conv())
	b.Write(target.GERTi.Conv())
	b.Write(source.Conv())
	b.WriteByte(byte(len(data)))
	b.Write(data)

	_, err := api.socket.Write([]byte(b.String()))
	if err != nil {
		return false, fmt.Errorf("error on write: %+v", err)
	}
	cmd, err := api.Parse()
	if err != nil {
		return false, fmt.Errorf("error on parse response: %+v", err)
	}
	if cmd.Command == CommandState {
		switch cmd.Status.Status {
		case StateFailure:
			return false, cmd.Status.ParseError()
		case StateSent:
			return true, nil
		case StateAssigned:
			return false, fmt.Errorf("invalid status: Assigned")
		case StateClosed:
			return false, fmt.Errorf("invalid status: Closed")
		case StateConnected:
			return false, fmt.Errorf("invalid status: Connected")
		}

	}
	return false, fmt.Errorf("no valid response")
}

func (api *Api) Shutdown() error {
	if api.socket != nil {
		_, err := api.socket.Write([]byte{byte(CommandClose)})
		if err != nil {
			return fmt.Errorf("error on write close command: %+v", err)
		}
		cmd, err := api.Parse()
		if err != nil {
			return fmt.Errorf("error on parsing response: %+v", err)
		}
		if cmd.Command == CommandState && cmd.Status.Status == StateClosed {
			err = api.socket.Close()
			if err != nil {
				return fmt.Errorf("error on close socker: %+v", err)
			}
			api.socket = nil
			return nil
		}
		return fmt.Errorf("no valid response received")
	}
	return fmt.Errorf("socket already closed")
}

func (api *Api) Parse() (Command, error) {
	data := make([]byte, 1024)
	//_, err := bufio.NewReader(api.socket).Read(data)
	_, err := api.socket.Read(data)
	if err != nil {
		return Command{}, fmt.Errorf("error on read data: %+v", err)
	}
	switch data[0] {
	case byte(CommandState):
		state, err := makeStatus(data[1:])
		if err != nil {
			return Command{}, fmt.Errorf("error while parsing status data: %+v", err)
		}
		return Command{
			Command: CommandState,
			Packet:  Packet{},
			Status:  state,
		}, nil
	case byte(CommandRegister):
		return Command{
			Command: CommandRegister,
			Packet:  Packet{},
			Status:  Status{},
		}, fmt.Errorf("geds returned command register")
	case byte(CommandData):
		packet, err := makePacket(data[1:])
		if err != nil {
			return Command{}, fmt.Errorf("error while parsing packet data: %+v", err)
		}
		return Command{
			Command: CommandData,
			Packet:  packet,
			Status:  Status{},
		}, nil
	case byte(CommandClose):
		err := api.socket.Close()
		if err != nil {
			return Command{}, fmt.Errorf("error while closing socket: %+v", err)
		}
		api.socket = nil
		return Command{
			Command: CommandClose,
			Packet:  Packet{},
			Status:  Status{},
		}, nil
	}
	return Command{}, fmt.Errorf("no valid command: %v", data[0])
}

func PrettyPrint(data []byte) (string, error) {
	output := ""
	switch data[0] {
	case byte(CommandState):
		state, err := makeStatus(data[1:])
		if err != nil {
			return "", fmt.Errorf("error while parsing status data: %+v", err)
		}
		output += "[STATE]"
		switch state.Status {
		case StateFailure:
			output += "[FAILURE]"
			switch state.Error {
			case ErrorVersion:
				output += "[VERSION]"
				break
			case ErrorBadKey:
				output += "[BAD_KEY]"
				break
			case ErrorAlreadyRegistered:
				output += "[ALREADY_REGISTERED]"
				break
			case ErrorNotRegistered:
				output += "[NOT_REGISTERED]"
				break
			case ErrorNoRoute:
				output += "[NO_ROUTE]"
				break
			case ErrorAddressTaken:
				output += "[ADDRESS_TAKEN]"
				break
			}
			break
		case StateConnected:
			output += "[CONNECTED]"
			output += "[" + state.Version.PrintVersion() + "]"
			break
		case StateAssigned:
			output += "[ASSIGNED]"
			break
		case StateClosed:
			output += "[CLOSED]"
			break
		case StateSent:
			output += "[SENT]"
			break
		}
		break
	case byte(CommandRegister):
		output += "[REGISTER]"
		addr := Parse(data[1:4])
		output += "[" + addr.PrintAddress() + "]"
		key := string(data[4:24])
		output += "[" + key + "]"
		break
	case byte(CommandData):
		source := ParseFull(data[1:7])
		target := Parse(data[7:10])
		length := data[10]
		dat := data[11 : 11+length]
		output += "[DATA]"
		output += fmt.Sprintf("[%v][%v][%v][%v]", source.PrintGERTc(), target.PrintAddress(), length, string(dat))
		break
	case byte(CommandClose):
		output += "[CLOSE]"
		break
	default:
		return "", fmt.Errorf("no valid command: %v", data[0])
	}
	return output, nil
}
