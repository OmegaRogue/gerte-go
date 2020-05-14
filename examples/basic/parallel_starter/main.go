package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

var Cmd1 string
var Cmd2 string

func init() {
	Cmd1 = os.Getenv("CMD2")
	Cmd2 = os.Getenv("CMD1")
}

func main() {
	err1 := make(chan error, 1)
	err2 := make(chan error, 1)
	out1 := make(chan []byte, 1)
	out2 := make(chan []byte, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go Run(out1, err1, &wg, Cmd1)
	wg.Add(1)
	go Run(out2, err2, &wg, Cmd2)
	wg.Wait()

	log.Printf("Output 1: %s", <-out1)
	log.Printf("Output 2: %s", <-out2)

	if err := <-err1; err != nil {
		log.Fatalf("Cmd1: %+v", err)
	}
	if err := <-err2; err != nil {
		log.Fatalf("Cmd2: %+v", err)
	}

}

func Run(out chan []byte, err chan error, wg *sync.WaitGroup, cmd string) {

	fmt.Println(os.Environ())
	defer wg.Done()
	cm := exec.Command("go", "run", "examples/basic/requester/main.go")
	out1, err1 := cm.CombinedOutput()
	out <- out1
	err <- err1

}
