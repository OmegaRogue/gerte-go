package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/OmegaRogue/gerte-go"
)

var Address string
var Key string
var ServerAddress string

func init() {
	Address = os.Getenv("ADDR")
	Key = os.Getenv("KEY")
	ServerAddress = os.Getenv("SERVER_ADDR")
}

func main() {

	api := gerte.NewApi(gerte.Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
	})

	addr := gerte.GertAddress{
		Upper: 1123,
		Lower: 1456,
	}
	b := string(addr.ToBytes()) + "aaaaaaaaaaaaaaaaaaaa"
	ioutil.WriteFile("test/resolutions.geds", []byte(b), os.ModePerm)

	con, err := net.Dial("tcp", "localhost:43780")
	if err != nil {
		log.Fatalf("error on tcp dial: %+v", err)
	}
	err = api.Startup(con)
	if err != nil {
		log.Fatalf("error on startup: %+v", err)
	}

	register, err := api.Register(addr, "aaaaaaaaaaaaaaaaaaaa")
	if err != nil {
		log.Fatalf("error on register: %+v", err)
	}
	log.Printf("registered: %v", register)

	targ := gerte.GERTc{
		GERTe: addr,
		GERTi: gerte.GertAddress{
			Upper: 1,
			Lower: 1,
		},
	}
	pkt := gerte.Packet{
		Source: targ,
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

func versionToBytes(ver gerte.Version) []byte {
	return []byte{ver.Major, ver.Minor, ver.Patch}
}

func toBytes(addr gerte.GertAddress) []byte {
	var b strings.Builder
	b.WriteByte(byte(addr.Upper >> 4))
	b.WriteByte(byte(((addr.Upper & 0x0F) << 4) | (addr.Lower >> 8)))
	b.WriteByte(byte(addr.Lower & 0xFF))
	return []byte(b.String())
}
