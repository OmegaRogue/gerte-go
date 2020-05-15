package main

import (
	"log"
	"net"
	"os"

	"github.com/OmegaRogue/gerte-go"
)

var address string
var targetAddress string
var key string
var serverAddress string

func init() {
	address = os.Getenv("ADDR")
	targetAddress = os.Getenv("TARGET_ADDR")
	key = os.Getenv("KEY")
	serverAddress = os.Getenv("SERVER_ADDR")
}

func main() {

	api := gerte.NewApi(gerte.Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
	})

	addr, err := gerte.AddressFromString(address)
	if err != nil {
		log.Fatalf("error on parse address string: %+v", err)
	}
	target, err := gerte.AddressFromString(targetAddress)
	if err != nil {
		log.Fatalf("error on parse target address string: %+v", err)
	}

	// b := string(addr.ToBytes()) + "aaaaaaaaaaaaaaaaaaaa"
	// ioutil.WriteFile("test/resolutions.geds", []byte(b), os.ModePerm)

	con, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatalf("error on tcp dial: %+v", err)
	}
	err = api.Startup(con)
	if err != nil {
		log.Fatalf("error on startup: %+v", err)
	}

	register, err := api.Register(addr, key)
	if err != nil {
		log.Fatalf("error on register: %+v", err)
	}
	log.Printf("registered: %v", register)

	source := gerte.GERTc{
		GERTe: gerte.GertAddress{},
		GERTi: gerte.GertAddress{},
	}
	targ := gerte.GERTc{
		GERTe: target,
		GERTi: gerte.GertAddress{
			Upper: 1,
			Lower: 1,
		},
	}
	pkt := gerte.Packet{
		Source: source,
		Target: targ,
		Data:   []byte("test"),
	}
	transmit, err := api.Transmit(pkt)
	if err != nil {
		log.Fatalf("error on transmit: %+v", err)
	}
	log.Printf("transmitted: %v", transmit)

	err = api.Shutdown()
	if err != nil {
		log.Fatalf("error on shutdown: %+v", err)
	}
}
