package GERTe

import (
	"bufio"
	"net"
	"sync"
	"testing"
)

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
	err := api.Startup(client)
	if err != nil {
		t.Errorf("client errored on startup: %+v", err)
	}

	err = api.Shutdown()
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
	err := api.Startup(client)
	if err != nil {
		t.Errorf("client errored on startup: %+v", err)
	}
	addr, _ := AddressFromString("0000.1999:0123.0456")
	_, err = api.Register(addr, "test")
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
	err := api.Startup(client)
	if err != nil {
		t.Errorf("client errored on startup: %+v", err)
	}
	addrE, _ := AddressFromString("0000.1999")
	addrI, _ := AddressFromString("0123.0456")
	gertC := GERTc{
		GERTe: addrE,
		GERTi: addrI,
	}
	_, err = api.Transmit(gertC, addrI, []byte("hello world!"))
	if err != nil {
		t.Errorf("client errored on transmit: %+v", err)
	}

	err = api.socket.Close()
	if err != nil {
		t.Errorf("client errored on close socket: %+v", err)
	}
	wg.Wait()
}
