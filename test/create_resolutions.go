package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/OmegaRogue/gerte-go"
)

var Address1 string
var Address2 string
var Key1 string
var Key2 string

func init() {
	Address1 = os.Getenv("ADDR1")
	Address2 = os.Getenv("ADDR2")
	Key1 = os.Getenv("KEY1")
	Key2 = os.Getenv("KEY2")
}

func main() {
	addr1, err := gerte.AddressFromString(Address1)
	if err != nil {
		log.Fatalf("error on parse address1 string: %+v", err)
	}
	addr2, err := gerte.AddressFromString(Address2)
	if err != nil {
		log.Fatalf("error on parse address2 string: %+v", err)
	}
	var b strings.Builder
	b.Write(addr1.ToBytes())
	b.WriteString(Key1)
	b.Write(addr2.ToBytes())
	b.WriteString(Key2)
	err = ioutil.WriteFile("test/resolutions.geds", []byte(b.String()), os.ModePerm)
	if err != nil {
		log.Fatalf("error on write resolutions: %+v", err)
	}

}
