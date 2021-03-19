package bueno

import (
	"bufio"
	"fmt"
	handler "github.com/ramzis/bueno/internal/pkg/connection"
	"log"
	"net"
	"strings"
	"time"
)

func (b *Bueno) HandleConnection(conn net.Conn) {

	err := handler.PerformHandshake(conn, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	pong := make(chan struct{})
	go b.KeepAlive(conn, pong)

	msgChan := make(chan string)
	go handler.ReadConn(conn, msgChan)

	defer b.RemoveConn(conn)
	b.RegisterConn(conn)

	// Wait for cmd or failed pong
	var s string
	var ok bool
	for {
		select {
			case _, ok := <-pong:
				if !ok {
					log.Println("Pong channel closed but shouldn't be!")
					return
				}
				close(pong)
				log.Println("Server failed to receive PONG")
				return
			case s, ok = <-msgChan:
				log.Println("Server reading...")
				if !ok {
					log.Println("Closed from error")
					return
				}
		}

		cmd, err := handler.Decode(s)
		if err != nil {
			log.Println(err)
			continue
		}

		switch cmd[0] {
		case "HI":
			log.Println("Unexpected HI after handshake")
		case "PONG":
			pong <- struct{}{}
		case "MSG":
			if len(cmd) > 1 {
				msg := strings.Join(cmd[1:], " ")
				b.TellEveryone(fmt.Sprintf("MSG %s %s", conn.RemoteAddr().String(), msg))
			}
		default:
			log.Println(s, "is unhandled")
		}
	}
}

func (b *Bueno) KeepAlive(conn net.Conn, pong chan struct{}) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	w := bufio.NewWriter(conn)

	writePing := func(delay time.Duration) {
		ticker.Stop()
		time.Sleep(delay)
		log.Println("Sending PING")
		w.WriteString("PING")
		w.WriteByte(0x0)
		w.Flush()
		ticker.Reset(time.Second * 5)
	}

	writePing(0)

	for {
		select {
			case <- ticker.C:
				pong <- struct{}{}
				return
			case <- pong:
				log.Println("Received PONG")
				writePing(time.Second * 5)
		}
	}
}
