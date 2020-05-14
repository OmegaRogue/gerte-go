package main

import (
	"log"
	"net"
	"os"

	"github.com/OmegaRogue/gerte-go"
)

var Address string
var Key string
var ServerAddress string

func init() {
	Address = os.Getenv("TARGET_ADDR")
	Key = os.Getenv("KEY")
	ServerAddress = os.Getenv("SERVER_ADDR")
}

func main() {

	api := gerte.NewApi(gerte.Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
	})

	addr, err := gerte.AddressFromString(Address)
	if err != nil {
		log.Fatalf("error on parse address string: %+v", err)
	}

	// b := string(addr.ToBytes()) + "aaaaaaaaaaaaaaaaaaaa"
	// ioutil.WriteFile("test/resolutions.geds", []byte(b), os.ModePerm)

	con, err := net.Dial("tcp", ServerAddress)
	if err != nil {
		log.Fatalf("error on tcp dial: %+v", err)
	}
	err = api.Startup(con)
	if err != nil {
		log.Fatalf("error on startup: %+v", err)
	}

	register, err := api.Register(addr, Key)
	if err != nil {
		log.Fatalf("error on register: %+v", err)
	}
	log.Printf("registered: %v", register)

	cmd, err := api.Parse()
	if err != nil {
		log.Fatalf("error on transmit: %+v", err)
	}
	p, err := cmd.PrintCommand()
	if err != nil {
		log.Fatalf("error on print command: %+v", err)
	}
	log.Printf("received: %v", p)

	err = api.Shutdown()
	if err != nil {
		log.Fatalf("error on shutdown: %+v", err)
	}
}
