package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/OmegaRogue/gerte-go"
)

var address1 string
var address2 string
var key1 string
var key2 string

func init() {
	address1 = os.Getenv("ADDR1")
	address2 = os.Getenv("ADDR2")
	key1 = os.Getenv("KEY1")
	key2 = os.Getenv("KEY2")
}

func main() {
	addr1, err := gerte.AddressFromString(address1)
	if err != nil {
		log.Fatalf("error on parse address1 string: %+v", err)
	}
	addr2, err := gerte.AddressFromString(address2)
	if err != nil {
		log.Fatalf("error on parse address2 string: %+v", err)
	}
	var b strings.Builder
	b.Write(addr1.ToBytes())
	b.WriteString(key1)
	b.Write(addr2.ToBytes())
	b.WriteString(key2)
	err = ioutil.WriteFile("test/resolutions.geds", []byte(b.String()), os.ModePerm)
	if err != nil {
		log.Fatalf("error on write resolutions: %+v", err)
	}

}
