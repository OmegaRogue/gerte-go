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
	t.Run("Successful Connect", StartupSuccessful)
	t.Run("Unsuccessful Connect", StartupUnsuccessful)
	t.Run("Unsuccessful Connect with Sent Status", StartupSent)
	t.Run("Unsuccessful Connect with Assigned Status", StartupAssigned)
	t.Run("Unsuccessful Connect with Invalid Command Response", StartupInvalidCmd)
}
func StartupSuccessful(t *testing.T) {
	server, client := net.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		dat := make([]byte, 1024)
		_, err := bufio.NewReader(server).Read(dat)
		if err != nil {
			t.Errorf("server errored on read: %+v", err)
		}

		t.Logf("server received: %v", versionFromBytes(dat).printVersion())

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
	api.Version = Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
	}

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
func StartupUnsuccessful(t *testing.T) {
	server, client := net.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		dat := make([]byte, 1024)
		_, err := bufio.NewReader(server).Read(dat)
		if err != nil {
			t.Errorf("server errored on read: %+v", err)
		}

		t.Logf("server received: %v", versionFromBytes(dat).printVersion())

		cmd := []byte{byte(CommandState), byte(StateFailure), byte(ErrorVersion)}
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
	api.Version = Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
	}

	err := api.Startup(client)
	if err != nil {
		if err.Error() == fmt.Sprintf("incompatible version during negotiation: %v", Version{0, 0, 0}) {
			t.Logf("client errored successfully on startup: %+v", err)
		} else {
			t.Errorf("client errored on startup: %+v", err)
		}

	}
	err = api.socket.Close()
	if err != nil {
		t.Errorf("client errored on close socket: %+v", err)
	}
	wg.Wait()
}
func StartupSent(t *testing.T) {
	server, client := net.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		dat := make([]byte, 1024)
		_, err := bufio.NewReader(server).Read(dat)
		if err != nil {
			t.Errorf("server errored on read: %+v", err)
		}

		t.Logf("server received: %v", versionFromBytes(dat).printVersion())

		cmd := []byte{byte(CommandState), byte(StateSent)}
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
	api.Version = Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
	}

	err := api.Startup(client)
	if err != nil {
		if err.Error() == "invalid response: state \"sent\"" {
			t.Logf("client errored successfully on startup: %+v", err)
		} else {
			t.Errorf("client errored on startup: %+v", err)
		}

	}
	err = api.socket.Close()
	if err != nil {
		t.Errorf("client errored on close socket: %+v", err)
	}
	wg.Wait()
}
func StartupAssigned(t *testing.T) {
	server, client := net.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		dat := make([]byte, 1024)
		_, err := bufio.NewReader(server).Read(dat)
		if err != nil {
			t.Errorf("server errored on read: %+v", err)
		}

		t.Logf("server received: %v", versionFromBytes(dat).printVersion())

		cmd := []byte{byte(CommandState), byte(StateAssigned)}
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
	api.Version = Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
	}

	err := api.Startup(client)
	if err != nil {
		if err.Error() == "invalid response: state \"assigned\"" {
			t.Logf("client errored successfully on startup: %+v", err)
		} else {
			t.Errorf("client errored on startup: %+v", err)
		}

	}
	err = api.socket.Close()
	if err != nil {
		t.Errorf("client errored on close socket: %+v", err)
	}
	wg.Wait()
}
func StartupInvalidCmd(t *testing.T) {
	server, client := net.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		dat := make([]byte, 1024)
		_, err := bufio.NewReader(server).Read(dat)
		if err != nil {
			t.Errorf("server errored on read: %+v", err)
		}

		t.Logf("server received: %v", versionFromBytes(dat).printVersion())

		cmd := []byte(string([]byte{byte(CommandRegister)}) +
			string(GertAddress{Upper: 0, Lower: 0}.toBytes()) +
			"testtesttesttesttesttest")
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
	api.Version = Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
	}

	err := api.Startup(client)
	if err != nil {
		if err.Error() == "error on parse response: geds returned command register" {
			t.Logf("client errored successfully on startup: %+v", err)
		} else {
			t.Errorf("client errored on startup: %+v", err)
		}

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
	t.Run("Successful", RegisterSuccessful)
	t.Run("BAD_KEY", RegisterBadKey)
	t.Run("ADDRESS_TAKEN", RegisterAddressTaken)
	t.Run("ALREADY_REGISTERED", RegisterAlreadyRegistered)

}
func RegisterSuccessful(t *testing.T) {
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
func RegisterBadKey(t *testing.T) {
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
		cmd := []byte{byte(CommandState), byte(StateFailure), byte(ErrorBadKey)}
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
		if err.Error() == "key did not match that used for the requested address. Requested address may not exist" {
			t.Logf("client errored successfully on register: %+v", err)
		} else {
			t.Errorf("client errored on register: %+v", err)
		}

	}
	err = api.socket.Close()
	if err != nil {
		t.Errorf("client errored on close socket: %+v", err)
	}
	wg.Wait()
}
func RegisterAddressTaken(t *testing.T) {
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
		cmd := []byte{byte(CommandState), byte(StateFailure), byte(ErrorAddressTaken)}
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
		if err.Error() == "address request has already been claimed" {
			t.Logf("client errored successfully on register: %+v", err)
		} else {
			t.Errorf("client errored on register: %+v", err)
		}

	}
	err = api.socket.Close()
	if err != nil {
		t.Errorf("client errored on close socket: %+v", err)
	}
	wg.Wait()
}
func RegisterAlreadyRegistered(t *testing.T) {
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
		cmd := []byte{byte(CommandState), byte(StateFailure), byte(ErrorAlreadyRegistered)}
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
		if err.Error() == "registration has already been performed successfully" {
			t.Logf("client errored successfully on register: %+v", err)
		} else {
			t.Errorf("client errored on register: %+v", err)
		}

	}
	err = api.socket.Close()
	if err != nil {
		t.Errorf("client errored on close socket: %+v", err)
	}
	wg.Wait()
}

func TestApi_Transmit(t *testing.T) {
	t.Run("Successful", TransmitSuccessful)
	t.Run("NOT_REGISTERED", TransmitNotRegistered)
	t.Run("NO_ROUTE", TransmitNoRoute)
}
func TransmitSuccessful(t *testing.T) {
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
func TransmitNotRegistered(t *testing.T) {
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
		cmd := []byte{byte(CommandState), byte(StateFailure), byte(ErrorNotRegistered)}
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
		if err.Error() == "gateway cannot send data before claiming an address" {
			t.Logf("client errored successfully on transmit: %+v", err)
		} else {
			t.Errorf("client errored on transmit: %+v", err)
		}
	}
	err = api.socket.Close()
	if err != nil {
		t.Errorf("client errored on close socket: %+v", err)
	}
	wg.Wait()
}
func TransmitNoRoute(t *testing.T) {
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
		cmd := []byte{byte(CommandState), byte(StateFailure), byte(ErrorNoRoute)}
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
		if err.Error() == "data failed to send because remote gateway could not be found" {
			t.Logf("client errored successfully on transmit: %+v", err)
		} else {
			t.Errorf("client errored on transmit: %+v", err)
		}
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
	p, err := cmd.PrintCommand()
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
		output += "[" + addr.PrintAddress() + "]"
		key := string(data[4:24])
		output += "[" + key + "]"
		break
	case byte(CommandData):
		source := gertCFromBytes(data[1:7])
		target := addressFromBytes(data[7:10])
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

func (ver Version) versionToBytes() []byte {
	return []byte{ver.Major, ver.Minor, ver.Patch}
}
func versionFromBytes(b []byte) Version {
	return Version{
		Major: b[0],
		Minor: b[1],
	}
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
