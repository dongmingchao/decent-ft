package main

import (
	"flag"
	"github.com/dongmingchao/decent-ft/src/caretaker"
	"github.com/dongmingchao/decent-ft/src/courier"
	"log"
	"net"
	"sync"
)

var remoteAddr string

func init() {
	flag.StringVar(&remoteAddr, "h", "", "host")
}

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(1)
	go caretaker.WatchDir(wg)
	if remoteAddr != "" {
		println("send to: ", remoteAddr)
		raddr, err := net.ResolveUDPAddr("udp", remoteAddr)
		if err != nil {
			log.Fatal(err)
		}
		courier.Send(raddr)
	} else {
		wg.Add(1)
		go courier.Start(wg)
	}
	wg.Wait()
}
