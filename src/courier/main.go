package courier

import (
	"bytes"
	"fmt"
	"github.com/dongmingchao/decent-ft/src/caretaker"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func freePort() *net.UDPAddr {
	raddr, _ := net.ResolveUDPAddr("udp", "8.8.8.8:22")
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = conn.Close()
	}()
	laddr, _ := net.ResolveUDPAddr("udp", conn.LocalAddr().String())
	return laddr
}
type Op byte

const (
	Done Op = 1 << iota
	AskIndex
	Get
	Remove
	Update
	Patch
)

func Send(raddr *net.UDPAddr) {
	conn, err := net.DialUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: 0}, raddr)
	if err != nil {
		log.Fatal(err)
	}
	//go recv(conn, buf)
	_, err = conn.Write([]byte{byte(AskIndex)})
	if err != nil {
		log.Fatal(err)
	}
	recv(conn, AskIndex)
	err = conn.Close()
	fmt.Println("链接关闭")
	if err != nil {
		log.Fatal(err)
	}
}

func recv(conn *net.UDPConn, askCode Op) {
	bakLen := 0
	var bakCode byte
	bak := make([]byte, 1024)
	n, err := conn.Read(bak)
	bakLen += n
	fmt.Println(string(bak[:bakLen]))
	if err != nil {
		log.Fatal(err)
	}
	if bakLen == 1 {
		bakCode = bak[0]
	}
	switch askCode {
	case AskIndex:
		if bakCode == byte(Done) {
			fmt.Println("read index")
			for {
				bak = make([]byte, 1024)
				n, err := conn.Read(bak)
				bakLen += n
				if err != nil {
					if err != io.EOF {
						log.Printf("Read error: %s", err)
					}
					break
				}
			}
		}
	}

	//println("remote addr: ", raddr.String())
	//println("msg length: ", n)
	//fmt.Println("recv: ", string(b[:bakLen]))
}

func Start() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	var conn *net.UDPConn
	var err error
	go func() {
		laddr := freePort()
		println("local addr: ", laddr.String())
		conn, err = net.ListenUDP("udp", laddr)
		if err != nil {
			fmt.Println("ListenUDP err:", err)
			return
		}
		for {
			ask := make([]byte, 1)
			n, raddr, err := conn.ReadFromUDP(ask)
			if err != nil {
				log.Fatal(err)
			}
			if n == 1 {
				switch Op(ask[0]) {
				case AskIndex:
					var buf bytes.Buffer
					buf.WriteByte(byte(Done))
					index, err := ioutil.ReadFile(caretaker.StashIndexFile)
					if err != nil {
						log.Fatal(err)
					}
					buf.Write(index)
					n, err = conn.WriteToUDP(buf.Bytes(), raddr)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
			println("remote addr: ", raddr.String())
			println("msg length: ", n)
			println("feed back finish: ", raddr.String())
		}
	}()
	fmt.Println("[Courier] Start")
	<-sigs
	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[Courier] Stop")
}

