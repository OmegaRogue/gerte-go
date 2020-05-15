package gerte

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"testing"
)

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

		t.Logf("server received: %v", VersionFromBytes(dat))

		cmd := []byte{byte(CommandState), byte(StateConnected)}
		verByte := Version{
			Major: 1,
			Minor: 1,
			Patch: 0,
		}.ToBytes()
		cmd = append(cmd, verByte[0], verByte[1])
		p, err := PrettyPrint(cmd)
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

		t.Logf("server received: %v", VersionFromBytes(dat))

		cmd := []byte{byte(CommandState), byte(StateFailure), byte(ErrorVersion)}
		p, err := PrettyPrint(cmd)
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

		t.Logf("server received: %v", VersionFromBytes(dat))

		cmd := []byte{byte(CommandState), byte(StateSent)}
		p, err := PrettyPrint(cmd)
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

		t.Logf("server received: %v", VersionFromBytes(dat))

		cmd := []byte{byte(CommandState), byte(StateAssigned)}
		p, err := PrettyPrint(cmd)
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

		t.Logf("server received: %v", VersionFromBytes(dat))

		cmd := []byte(string([]byte{byte(CommandRegister)}) +
			string(GertAddress{Upper: 0, Lower: 0}.ToBytes()) +
			"testtesttesttesttesttest")
		p, err := PrettyPrint(cmd)
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
		p, err := PrettyPrint(dat)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server received: %v", p)
		cmd := []byte{byte(CommandState), byte(StateClosed)}
		p, err = PrettyPrint(cmd)
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
		p, err := PrettyPrint(dat)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server received: %v", p)
		cmd := []byte{byte(CommandState), byte(StateAssigned)}
		p, err = PrettyPrint(cmd)
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
		p, err := PrettyPrint(dat)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server received: %v", p)
		cmd := []byte{byte(CommandState), byte(StateFailure), byte(ErrorBadKey)}
		p, err = PrettyPrint(cmd)
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
		p, err := PrettyPrint(dat)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server received: %v", p)
		cmd := []byte{byte(CommandState), byte(StateFailure), byte(ErrorAddressTaken)}
		p, err = PrettyPrint(cmd)
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
		p, err := PrettyPrint(dat)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server received: %v", p)
		cmd := []byte{byte(CommandState), byte(StateFailure), byte(ErrorAlreadyRegistered)}
		p, err = PrettyPrint(cmd)
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
		p, err := PrettyPrint(dat)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server received: %v", p)
		cmd := []byte{byte(CommandState), byte(StateSent)}
		p, err = PrettyPrint(cmd)
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
	pkt := Packet{
		Source: gertC,
		Target: gertC,
		Data:   []byte("hello world!"),
	}
	_, err := api.Transmit(pkt)
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
		p, err := PrettyPrint(dat)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server received: %v", p)
		cmd := []byte{byte(CommandState), byte(StateFailure), byte(ErrorNotRegistered)}
		p, err = PrettyPrint(cmd)
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
	pkt := Packet{
		Source: gertC,
		Target: gertC,
		Data:   []byte("hello world!"),
	}
	_, err := api.Transmit(pkt)

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
		p, err := PrettyPrint(dat)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server received: %v", p)
		cmd := []byte{byte(CommandState), byte(StateFailure), byte(ErrorNoRoute)}
		p, err = PrettyPrint(cmd)
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
	pkt := Packet{
		Source: gertC,
		Target: gertC,
		Data:   []byte("hello world!"),
	}
	_, err := api.Transmit(pkt)
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
		}.ToBytes()
		cmd = append(cmd, verByte[0], verByte[1])
		p, err := PrettyPrint(cmd)
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

	t.Logf("client received: %#v", cmd)
	err = api.socket.Close()
	if err != nil {
		t.Errorf("client errored on close socket: %+v", err)
	}
	wg.Wait()
}
