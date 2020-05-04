// Package gerte provides an api for the GERT system.
// More info: https://github.com/GlobalEmpire/GERT
package gerte

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	// CommandState is used by relays to indicate result of a command to a gateway and by gateways to request the state.
	CommandState GertCommand = 0
	// CommandRegister claims an address for this gateway using a key.
	CommandRegister GertCommand = 1
	// CommandData transmits data from a GERTi address to a GERTc address.
	CommandData GertCommand = 2
	// CommandClose gracefully closes a connection to a relay.
	CommandClose GertCommand = 3

	// StateFailure is the initial gateway state.
	// Should be changed upon negotiation.
	// Also a response to failed commands with an error
	StateFailure GertStatus = 0
	// StateConnected indicates that the Gateway is connected to a relay, no other action has been taken
	StateConnected GertStatus = 1
	// StateAssigned indicates that the Gateway has successfully claimed an address.
	// Used as a response to the REGISTER command.
	StateAssigned GertStatus = 2
	// StateClosed indicates that the Gateway has closed the connection.
	// Used as a response to the CLOSE command.
	StateClosed GertStatus = 3
	// StateSent indicates that the Data has been successfully sent.
	// Used as a response to the DATA command.
	// This is not a guarantee for data that has to be sent to another peer, although it's unlikely to be incorrect.
	StateSent GertStatus = 4

	// ErrorVersion indicates that an incompatible version was used during negotiation.
	ErrorVersion GertError = 0
	// ErrorBadKey indicates that the Key did not match that used for the requested address.
	// Requested address may not exist
	ErrorBadKey GertError = 1
	// ErrorAlreadyRegistered indicates that Registration has already been performed successfully
	ErrorAlreadyRegistered GertError = 2
	// ErrorNotRegistered indicates that the Gateway hasn't been registered yet.
	// The gateway cannot send data before claiming an address
	ErrorNotRegistered GertError = 3
	// ErrorNoRoute indicates that Data failed to send because the remote gateway couldn't be found
	ErrorNoRoute GertError = 4
	// ErrorAddressTaken indicates that the Address request has already been claimed
	ErrorAddressTaken GertError = 5
)

type (
	// GertStatus indicates the Status from a Status Command
	GertStatus byte

	// GertCommand indicates the Command of a Request
	GertCommand byte

	// GertError is the Error Code in a "Failed" Status
	GertError byte
)

type (
	// GertAddress is a 3 byte address used as a GERTe/i Address
	GertAddress struct {
		Upper int
		Lower int
	}

	// GERTc is a 6 byte GERTc Address
	GERTc struct {
		GERTe GertAddress
		GERTi GertAddress
	}
)

// Api is used to perform GERTe API Operations
type Api struct {
	socket     net.Conn
	listener   net.Listener
	Registered bool
	Address    GertAddress
	Version    Version
}

type (
	// Command is the Parsed Data returned from Api.Parse
	Command struct {
		Command GertCommand
		Packet  Packet
		Status  Status
	}

	// Packet is the Parsed Data returned from Api.Parse for a Data Command
	Packet struct {
		Source GERTc
		Target GERTc
		Data   []byte
	}

	// Status is the Parsed Status returned from Api.Parse for a Status Command
	Status struct {
		Status  GertStatus
		Size    byte
		Error   GertError
		Version Version
	}

	// Version is the Version used by the Connected Status
	Version struct {
		Major byte
		Minor byte
		Patch byte
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
				Major: data[1],
				Minor: data[2],
				Patch: data[3],
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
		return fmt.Errorf("socket already open")
	}
	api.socket = c
	cmd, err := api.Parse()
	if err != nil {
		return fmt.Errorf("error on parse response: %+v", err)
	}
	if cmd.Command == CommandState {
		switch cmd.Status.Status {
		case StateConnected:
			api.Version = cmd.Status.Version
			return nil
		case StateFailure:
			return cmd.Status.parseError()
		case StateSent:
			return fmt.Errorf("invalid response: state \"sent\"")
		case StateClosed:
			err := api.socket.Close()
			if err != nil {
				return fmt.Errorf("error while closing socket: %+v", err)
			}
			api.socket = nil
		case StateAssigned:
			return fmt.Errorf("invalid response: state \"assigned\"")
		}

	}
	return fmt.Errorf("invalid response: %+v", cmd)
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
