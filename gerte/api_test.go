package gerte

import (
	"bufio"
	"net"
	"sync"
	"testing"
)

const (
	Server   = "server"
	Client   = "client"
	Sent     = "sent"
	Received = "received"
	Error    = "error"
)

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
		dat = make([]byte, 1024)
		_, err = bufio.NewReader(server).Read(dat)
		if err != nil {
			t.Errorf("server errored on read: %+v", err)
		}
		p, err = PrettyPrint(dat)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server received: %v", p)
		cmd = []byte{byte(CommandState), byte(StateClosed)}
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
	addr, _ := FromString("0000.1999:0123.0456")
	_, err = api.Register(addr, "test.txt")
	if err != nil {
		t.Errorf("client errored on register: %+v", err)
	}
	// Do some stuff
	err = api.Shutdown()
	if err != nil {
		t.Errorf("client errored on shutdown: %+v", err)
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
		dat = make([]byte, 1024)
		_, err = bufio.NewReader(server).Read(dat)
		if err != nil {
			t.Errorf("server errored on read: %+v", err)
		}
		p, err = PrettyPrint(dat)
		if err != nil {
			t.Errorf("server errored on pretty print: %+v", err)
		}
		t.Logf("server received: %v", p)
		// Do some stuff
		wg.Done()
	}()

	var api Api
	err := api.Startup(client)
	if err != nil {
		t.Errorf("client errored on startup: %+v", err)
	}
	addrE, _ := FromString("0000.1999")
	addrI, _ := FromString("0123.0456")
	gertC := GERTc{
		GERTe: addrE,
		GERTi: addrI,
	}
	_, err = api.Transmit(gertC, addrI, []byte("hello world!"))
	if err != nil {
		t.Errorf("client errored on transmit: %+v", err)
	}

	// Do some stuff
	err = api.Shutdown()
	if err != nil {
		t.Errorf("client errored on shutdown: %+v", err)
	}
	wg.Wait()
}
