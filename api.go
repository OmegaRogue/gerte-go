// Package gerte provides an api for the GERT system.
// More info: https://github.com/GlobalEmpire/GERT
package gerte

import (
	"fmt"
	"net"
	"strings"
)

// Api is used to perform GERTe API Operations
type Api struct {
	socket     net.Conn
	listener   net.Listener
	Registered bool
	Address    GertAddress
	Version    Version
}

// NewApi is the constructor for Api, it assigns the Version
func NewApi(ver Version) *Api {
	api := new(Api)
	api.Version = ver

	return api
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
	_, err := api.socket.Write([]byte{api.Version.Major, api.Version.Minor})
	if err != nil {
		return fmt.Errorf("error on send version: %w", err)
	}
	cmd, err := api.Parse()
	if err != nil {
		return fmt.Errorf("error on parse response: %w", err)
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
				return fmt.Errorf("error while closing socket: %w", err)
			}
			api.socket = nil
		case StateAssigned:
			return fmt.Errorf("invalid response: state \"assigned\"")
		}

	}
	return fmt.Errorf("invalid response: %v", cmd)
}

// Register registers the GERTe client on the GERTe address with the associated 20 byte key.
// It returns a bool whether the registration was successful and any encountered errors.
// All gateways must register themselves with a valid GERTe address and key before sending data.
// The API does not track if it is registered nor what address it has registered, it instead relies on the relay it is connected to.
func (api *Api) Register(addr GertAddress, key string) (bool, error) {
	var b strings.Builder
	b.WriteByte(byte(CommandRegister))
	b.Write(addr.ToBytes())
	b.WriteString(key)
	_, err := api.socket.Write([]byte(b.String()))
	if err != nil {
		return false, fmt.Errorf("error on write: %w", err)
	}
	cmd, err := api.Parse()
	if err != nil {
		return false, fmt.Errorf("error parsing response: %w", err)
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
	return false, fmt.Errorf("no valid response: %v", cmd)
}

// Transmit sends data to the target Address.
// It returns a bool whether the operation was successful and any encountered errors.
// The official API only allows transmissions from GERTi to GERTi via GERTe.
// his means that a GERTi address must be provided for each endpoint in a message.
func (api *Api) Transmit(pkt Packet) (bool, error) {

	if api.socket == nil {
		return false, fmt.Errorf("not connected")
	}

	var b strings.Builder

	b.WriteByte(byte(CommandData))
	data, err := pkt.ToBytes()
	if err != nil {
		return false, fmt.Errorf("error on marshal packet: %w", err)
	}
	b.Write(data)

	_, err = api.socket.Write([]byte(b.String()))
	if err != nil {
		return false, fmt.Errorf("error on write: %w", err)
	}
	cmd, err := api.Parse()
	if err != nil {
		return false, fmt.Errorf("error on parse response: %w", err)
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
	return false, fmt.Errorf("no valid response: %v", cmd)
}

// Shutdown Gracefully closes the GERTe Socket.
// It returns any errors encountered.
// The official API prefers using a safe shutdown procedure, although the GEDS servers should be more than stable enough to survive any number of unclean shutdowns.
func (api *Api) Shutdown() error {
	if api.socket != nil {
		_, err := api.socket.Write([]byte{byte(CommandClose)})
		if err != nil {
			return fmt.Errorf("error on write close command: %w", err)
		}
		cmd, err := api.Parse()
		if err != nil {
			return fmt.Errorf("error on parsing response: %w", err)
		}
		if cmd.Command == CommandState && cmd.Status.Status == StateClosed {
			err = api.socket.Close()
			if err != nil {
				return fmt.Errorf("error on close socket: %w", err)
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
	n, err := api.socket.Read(data)
	if err != nil {
		return Command{}, fmt.Errorf("error on read data (%v bytes): %w", n, err)
	}
	cmd, err := CommandFromBytes(data)
	if err != nil {
		return Command{}, fmt.Errorf("error parsing command: %w", err)
	}

	switch cmd.Command {
	case CommandRegister:
		return cmd, fmt.Errorf("geds returned command register")
	case CommandClose:
		err := api.socket.Close()
		if err != nil {
			return Command{}, fmt.Errorf("error while closing socket: %w", err)
		}
		api.socket = nil
		break
	}

	return cmd, nil
}
