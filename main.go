package main

import (
	"flag"
	"github.com/dongmingchao/decent-ft/src/courier"
	"log"
	"net"
)

var remoteAddr string

func init() {
	flag.StringVar(&remoteAddr, "h", "", "host")
}

func main() {
	flag.Parse()
	//go caretaker.WatchDir()
	if remoteAddr != "" {
		println("send to: ", remoteAddr)
		raddr, err := net.ResolveUDPAddr("udp", remoteAddr)
		if err != nil {
			log.Fatal(err)
		}
		courier.Send(raddr)
	} else {
		courier.Start()
	}
}
