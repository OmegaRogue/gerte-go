// Package GERTe provides an api for the GERT system
// more info: https://github.com/GlobalEmpire/GERT
package gerte

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	// Used by relays to indicate result of a command to a gateway and by gateways to request the state.
	CommandState GertCommand = 0
	// Claims an address for this gateway using a key.
	CommandRegister GertCommand = 1
	// Transmits data from a GERTi address to a GERTc address.
	CommandData GertCommand = 2
	// Gracefully closes a connection to a relay.
	CommandClose GertCommand = 3

	// Initial gateway state.
	// Should be changed upon negotiation.
	// Also a response to failed commands with an error
	StateFailure GertStatus = 0
	// Gateway is connected to a relay, no other action has been taken
	StateConnected GertStatus = 1
	// Gateway has successfully claimed an address.
	// Used as a response to the REGISTER command.
	StateAssigned GertStatus = 2
	// Gateway has closed the connection.
	// Used as a response to the CLOSE command.
	StateClosed GertStatus = 3
	// Data has been successfully sent.
	// Used as a response to the DATA command.
	// This is not a guarantee for data that has to be sent to another peer, although it's unlikely to be incorrect.
	StateSent GertStatus = 4

	// Incompatible version during negotiation.
	ErrorVersion GertError = 0
	// Key did not match that used for the requested address.
	// Requested address may not exist
	ErrorBadKey GertError = 1
	// Registration has already been performed successfully
	ErrorAlreadyRegistered GertError = 2
	// Gateway cannot send data before claiming an address
	ErrorNotRegistered GertError = 3
	// Data failed to send because remote gateway could not be found
	ErrorNoRoute GertError = 4
	// Address request has already been claimed
	ErrorAddressTaken GertError = 5
)

type (
	GertStatus byte

	GertCommand byte

	GertError byte
)

// GERT addresses consist of 3 or 6 bytes depending on the usage.
// GERTe/i is 3 bytes while GERTc is 6 bytes.
// GEDS never parses GERTi addresses, however it will parse and enforce GERTe addresses.
type (
	GertAddress struct {
		Upper int
		Lower int
	}

	GERTc struct {
		GERTe GertAddress
		GERTi GertAddress
	}
)

type Api struct {
	socket     net.Conn
	listener   net.Listener
	Registered bool
	Address    GertAddress
}

type (
	Command struct {
		Command GertCommand
		Packet  Packet
		Status  Status
	}

	Packet struct {
		Source GERTc
		Target GERTc
		Data   []byte
	}

	Status struct {
		Status  GertStatus
		Size    int
		Error   GertError
		Version Version
	}

	Version struct {
		Major int
		Minor int
		Patch int
	}
)

func (status Status) parseError() error {
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

// AddressFromString converts a string with an address in the format "XXXX.YYYY" into the corresponding GertAddress.
// It returns the GertAddress and any encountered errors.
func AddressFromString(addr string) (GertAddress, error) {
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

func (addr GertAddress) toBytes() []byte {
	var b strings.Builder
	b.WriteByte(byte(addr.Upper >> 4))
	b.WriteByte(byte(((addr.Upper & 0x0F) << 4) | (addr.Lower >> 8)))
	b.WriteByte(byte(addr.Lower & 0xFF))
	return []byte(b.String())
}

func addressFromBytes(data []byte) GertAddress {
	return GertAddress{
		Upper: (int(data[0]) << 4) | (int(data[1]) >> 4),
		Lower: ((int(data[1]) & 0x0F) << 8) | int(data[2]),
	}
}

func gertCFromBytes(data []byte) GERTc {
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

func (ver Version) printVersion() string {
	return fmt.Sprintf("%v.%v.%v", ver.Major, ver.Minor, ver.Patch)
}
func (addr GertAddress) printAddress() string {
	return fmt.Sprintf("%v.%v", addr.Upper, addr.Lower)
}
func (addr GERTc) printGERTc() string {
	return fmt.Sprintf("%v.%v:%v.%v", addr.GERTe.Upper, addr.GERTe.Lower, addr.GERTi.Upper, addr.GERTi.Lower)
}

func makePacket(data []byte) (Packet, error) {
	if len(data) < 13 {
		return Packet{}, fmt.Errorf("data too short: %v<13", len(data))
	}

	source := gertCFromBytes(data[:6])
	target := gertCFromBytes(data[6:12])
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

// Startup initializes the API.
// It returns any encountered errors.
// Initializing the API is incredibly simple.
// The API should be initialized for every program when it decides to use it (although it is potentially already initialized, the API will ensure it's safe to initialize.)
func (api *Api) Startup(c net.Conn) error {
	if api.socket != nil {
		return nil
	}

	api.socket = c
	return nil
}

// Register registers the GERTe client on the GERTe address with the associated 20 byte key.
// It returns a bool whether the registration was successful and any encountered errors.
// All gateways must register themselves with a valid GERTe address and key before sending data.
// The API does not track if it is registered nor what address it has registered, it instead relies on the relay it is connected to.
func (api *Api) Register(addr GertAddress, key string) (bool, error) {
	var b strings.Builder
	b.WriteByte(byte(CommandRegister))
	rawAddr := addr.toBytes()
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
			return false, cmd.Status.parseError()
		case StateAssigned:
			api.Address = addr
			return true, nil
		}
	}
	return false, fmt.Errorf("no valid response")
}

// Transmit sends data to the target Address.
// It returns a bool whether the operation was successful and any encountered errors.
// The official API only allows transmissions from GERTi to GERTi via GERTe.
// his means that a GERTi address must be provided for each endpoint in a message.
func (api *Api) Transmit(target GERTc, source GertAddress, data []byte) (bool, error) {

	if api.socket == nil {
		return false, fmt.Errorf("not connected")
	}
	if len(data) > 255 {
		return false, fmt.Errorf("data cannot exceed 255 bytes")
	}

	var b strings.Builder

	b.WriteByte(byte(CommandData))
	b.Write(target.GERTe.toBytes())
	b.Write(target.GERTi.toBytes())
	b.Write(source.toBytes())
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
			return false, cmd.Status.parseError()
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

// Shutdown Gracefully closes the GERTe Socket.
// It returns any errors encountered.
// The official API prefers using a safe shutdown procedure, although the GEDS servers should be more than stable enough to survive any number of unclean shutdowns.
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
				return fmt.Errorf("error on close socket: %+v", err)
			}
			api.socket = nil
			return nil
		}
		return fmt.Errorf("no valid response received")
	}
	return fmt.Errorf("socket already closed")
}

// Parse reads data from the GERTe socket and parses it.
// It returns the received Command and any errors encountered.
// The official API only checks the connection for data when requested.
// This includes connection closures from the relay.
// If the connection is closed, the API will call the error function instead of returning anything.
func (api *Api) Parse() (Command, error) {
	data := make([]byte, 1024)
	// _, err := bufio.NewReader(api.socket).Read(data)
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

// PrettyPrint returns a GERTe Message in a Human readable string.
// Mainly for testing purposes
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
			output += "[" + state.Version.printVersion() + "]"
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
		addr := addressFromBytes(data[1:4])
		output += "[" + addr.printAddress() + "]"
		key := string(data[4:24])
		output += "[" + key + "]"
		break
	case byte(CommandData):
		source := gertCFromBytes(data[1:7])
		target := addressFromBytes(data[7:10])
		length := data[10]
		dat := data[11 : 11+length]
		output += "[DATA]"
		output += fmt.Sprintf("[%v][%v][%v][%v]", source.printGERTc(), target.printAddress(), length, string(dat))
		break
	case byte(CommandClose):
		output += "[CLOSE]"
		break
	default:
		return "", fmt.Errorf("no valid command: %v", data[0])
	}
	return output, nil
}
