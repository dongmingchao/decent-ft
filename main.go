package main

import (
	"flag"
	"fmt"
	"github.com/dongmingchao/decent-ft/src/caretaker"
	"github.com/dongmingchao/decent-ft/src/courier"
	"log"
	"math"
	"net"
	"sync"
)

var remoteAddr string

func init() {
	flag.StringVar(&remoteAddr, "h", "", "host")
}

func main() {
	flag.Parse()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go caretaker.WatchDir(&wg)
	if remoteAddr != "" {
		println("send to: ", remoteAddr)
		raddr, err := net.ResolveUDPAddr("udp", remoteAddr)
		if err != nil {
			log.Fatal(err)
		}
		courier.Send(raddr)
	} else {
		wg.Add(1)
		go courier.Start(&wg)
	}
	wg.Wait()
}

func intMaxLength() {
	const MaxUint = ^uint(0)
	const MinUint = 0
	const MaxInt = int(MaxUint >> 1)
	const MinInt = -MaxInt - 1
	fmt.Println(MaxInt)
	fmt.Println(math.MaxInt64)
}
