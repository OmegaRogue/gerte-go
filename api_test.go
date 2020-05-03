package gerte

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestAddressFromString(t *testing.T) {
	addr, err := AddressFromString("0123.0456")
	if err != nil {
		t.Errorf("error on parse address string: %+v", err)
	}
	addrT := GertAddress{
		Upper: 123,
		Lower: 456,
	}
	if addr.Upper != addrT.Upper || addr.Lower != addrT.Lower {
		t.Error("addresses don't match")
	}
}

func TestApi_Startup(t *testing.T) {
	server, client := net.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		cmd := []byte{byte(CommandState), byte(StateConnected)}
		verByte := Version{
			Major: 1,
			Minor: 1,
			Patch: 0,
		}.versionToBytes()
		cmd = append(cmd, verByte[0], verByte[1], verByte[2])
		p, err := prettyPrint(cmd)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server sent: %v", p)
		_, err = server.Write(cmd)
		if err != nil {
			t.Errorf("server errored on write: %+v", err)
		}

		err = server.Close()
		if err != nil {
			t.Errorf("server errored on close: %+v", err)
		}
		wg.Done()
	}()

	var api Api
	err := api.Startup(client)
	if err != nil {
		t.Errorf("client errored on startup: %+v", err)
	}
	err = api.socket.Close()
	if err != nil {
		t.Errorf("client errored on close socket: %+v", err)
	}
	wg.Wait()
}

func TestApi_Shutdown(t *testing.T) {
	server, client := net.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dat := make([]byte, 1024)
		_, err := bufio.NewReader(server).Read(dat)
		if err != nil {
			t.Errorf("server errored on read: %+v", err)
		}
		p, err := prettyPrint(dat)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server received: %v", p)
		cmd := []byte{byte(CommandState), byte(StateClosed)}
		p, err = prettyPrint(cmd)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server sent: %v", p)
		_, err = server.Write(cmd)
		if err != nil {
			t.Errorf("server errored on write: %+v", err)
		}
		// Do some stuff
		err = server.Close()
		if err != nil {
			t.Errorf("server errored on close: %+v", err)
		}
		wg.Done()
	}()

	var api Api
	api.socket = client

	err := api.Shutdown()
	if err != nil {
		t.Errorf("client errored on shutdown: %+v", err)
	}
	wg.Wait()
}

func TestApi_Register(t *testing.T) {
	server, client := net.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dat := make([]byte, 1024)
		_, err := bufio.NewReader(server).Read(dat)
		if err != nil {
			t.Errorf("server errored on read: %+v", err)
		}
		p, err := prettyPrint(dat)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server received: %v", p)
		cmd := []byte{byte(CommandState), byte(StateAssigned)}
		p, err = prettyPrint(cmd)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server sent: %v", p)
		_, err = server.Write(cmd)
		if err != nil {
			t.Errorf("server errored on write: %+v", err)
		}

		err = server.Close()
		if err != nil {
			t.Errorf("server errored on close: %+v", err)
		}
		wg.Done()
	}()

	var api Api
	api.socket = client
	addr, _ := AddressFromString("0000.1999")
	_, err := api.Register(addr, "test")
	if err != nil {
		t.Errorf("client errored on register: %+v", err)
	}
	err = api.socket.Close()
	if err != nil {
		t.Errorf("client errored on close socket: %+v", err)
	}
	wg.Wait()
}

func TestApi_Transmit(t *testing.T) {
	server, client := net.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		dat := make([]byte, 1024)
		_, err := bufio.NewReader(server).Read(dat)
		if err != nil {
			t.Errorf("server errored on read: %+v", err)
		}
		p, err := prettyPrint(dat)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server received: %v", p)
		cmd := []byte{byte(CommandState), byte(StateSent)}
		p, err = prettyPrint(cmd)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server sent: %v", p)
		_, err = server.Write(cmd)
		if err != nil {
			t.Errorf("server errored on write: %+v", err)
		}

		err = server.Close()
		if err != nil {
			t.Errorf("server errored on close: %+v", err)
		}
		wg.Done()
	}()

	var api Api
	api.socket = client
	addrE, _ := AddressFromString("0000.1999")
	addrI, _ := AddressFromString("0123.0456")
	gertC := GERTc{
		GERTe: addrE,
		GERTi: addrI,
	}
	_, err := api.Transmit(gertC, addrI, []byte("hello world!"))
	if err != nil {
		t.Errorf("client errored on transmit: %+v", err)
	}

	err = api.socket.Close()
	if err != nil {
		t.Errorf("client errored on close socket: %+v", err)
	}
	wg.Wait()
}

func TestApi_Parse(t *testing.T) {
	server, client := net.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		cmd := []byte{byte(CommandState), byte(StateConnected)}
		verByte := Version{
			Major: 1,
			Minor: 1,
			Patch: 0,
		}.versionToBytes()
		cmd = append(cmd, verByte[0], verByte[1], verByte[2])
		p, err := prettyPrint(cmd)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server sent: %v", p)
		_, err = server.Write(cmd)
		if err != nil {
			t.Errorf("server errored on write: %+v", err)
		}

		err = server.Close()
		if err != nil {
			t.Errorf("server errored on close: %+v", err)
		}
		wg.Done()
	}()

	var api Api
	api.socket = client
	cmd, err := api.Parse()
	if err != nil {
		t.Errorf("client errored on parse response: %+v", err)
	}
	p, err := cmd.printCommand()
	if err != nil {
		t.Errorf("client errored on print command: %+v", err)
	}
	t.Logf("client received: %v", p)
	err = api.socket.Close()
	if err != nil {
		t.Errorf("client errored on close socket: %+v", err)
	}
	wg.Wait()
}

func (command Command) printCommand() (string, error) {
	output := ""
	switch command.Command {
	case CommandState:
		output += "[STATE]"
		switch command.Status.Status {
		case StateFailure:
			output += "[FAILURE]"
			switch command.Status.Error {
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
			output += "[" + command.Status.Version.printVersion() + "]"
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
	case CommandClose:
		output += "[CLOSE]"
		break
	default:
		return "", fmt.Errorf("no valid command: %v", command.Command)
	}
	return output, nil
}

// ...
func prettyPrint(data []byte) (string, error) {
	output := ""
	switch data[0] {
	case byte(CommandState):
		state, err := makeStatus(data[1:])
		if err != nil {
			return "", fmt.Errorf("error while parsing status data: %+v", err)
		}
		p, err := state.printStatus()
		if err != nil {
			return "", fmt.Errorf("error while printing status: %+v", err)
		}
		output += "[STATE]" + p
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

func (ver Version) printVersion() string {
	return fmt.Sprintf("%v.%v.%v", ver.Major, ver.Minor, ver.Patch)
}
func (addr GertAddress) printAddress() string {
	return fmt.Sprintf("%v.%v", addr.Upper, addr.Lower)
}
func (addr GERTc) printGERTc() string {
	return fmt.Sprintf("%v.%v:%v.%v", addr.GERTe.Upper, addr.GERTe.Lower, addr.GERTi.Upper, addr.GERTi.Lower)
}

func (ver Version) versionToBytes() []byte {
	return []byte{ver.Major, ver.Minor, ver.Patch}
}

func addressFromBytes(data []byte) GertAddress {
	return GertAddress{
		Upper: (int(data[0]) << 4) | (int(data[1]) >> 4),
		Lower: ((int(data[1]) & 0x0F) << 8) | int(data[2]),
	}
}

func (status Status) printFailure() (string, error) {
	switch status.Error {
	case ErrorVersion:
		return "[VERSION]", nil
	case ErrorBadKey:
		return "[BAD_KEY]", nil
	case ErrorAlreadyRegistered:
		return "[ALREADY_REGISTERED]", nil
	case ErrorNotRegistered:
		return "[NOT_REGISTERED]", nil
	case ErrorNoRoute:
		return "[NO_ROUTE]", nil
	case ErrorAddressTaken:
		return "[ADDRESS_TAKEN]", nil
	}
	return "", fmt.Errorf("invalid Failure")
}

func (status Status) printStatus() (string, error) {
	switch status.Status {
	case StateFailure:
		p, err := status.printFailure()
		if err != nil {
			return "", fmt.Errorf("error while parsing Failure: %+v", err)
		}
		return "[FAILURE]" + p, nil
	case StateConnected:
		return "[CONNECTED][" + status.Version.printVersion() + "]", nil
	case StateAssigned:
		return "[ASSIGNED]", nil
	case StateClosed:
		return "[CLOSED]", nil
	case StateSent:
		return "[SENT]", nil
	}
	return "", fmt.Errorf("invalid status")
}
