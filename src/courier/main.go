package courier

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/dongmingchao/decent-ft/src/caretaker"
	resourcePool "github.com/dongmingchao/decent-ft/src/resource-pool"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
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
	fmt.Println("获取Index信息完毕")
	if err != nil {
		log.Fatal(err)
	}
}

func recv(conn *net.UDPConn, askCode Op) {
	reader := bufio.NewReader(conn)
	length,err := reader.ReadByte()
	if err != nil {
		log.Fatal(err)
	}
	bakCode, err := reader.ReadByte()
	if err != nil {
		log.Fatal(err)
	}
	println("content len: ", length)
	println("back code: ", bakCode)
	switch askCode {
	case AskIndex:
		if bakCode == byte(Done) {
			buf := make([]byte, length - 1)
			reader.Read(buf)
			stash := resourcePool.GTree{}
			stash.Read(bytes.NewReader(buf))
			fmt.Println(stash)
		}
	}

	//println("remote addr: ", raddr.String())
	//println("msg length: ", n)
	//fmt.Println("recv: ", string(b[:bakLen]))
}

func Start(wg *sync.WaitGroup) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	var conn net.PacketConn
	var err error
	go func() {
		conn, err = net.ListenPacket("udp", ":0")
		//laddr := freePort()
		if err != nil {
			fmt.Println("ListenUDP err:", err)
			return
		}
		println("local addr: ", conn.LocalAddr().String())
		for {
			ask := make([]byte, 1)
			n, raddr, err := conn.ReadFrom(ask)
			if err != nil {
				// connect closed
				log.Printf("%v", err)
				break
			}
			if n == 1 {
				switch Op(ask[0]) {
				case AskIndex:
					var buf bytes.Buffer
					buf.WriteByte(byte(Done))
					if err != nil {
						log.Fatal(err)
					}
					caretaker.GlobalStash.Write(&buf)
					var fin bytes.Buffer
					fin.WriteByte(byte(len(buf.Bytes())))
					fin.ReadFrom(&buf)
					n, err = conn.WriteTo(fin.Bytes(), raddr)
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
	defer wg.Done()
}
