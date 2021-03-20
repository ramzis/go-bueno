package bueno

import (
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
		log.Println(err)
		return
	}

	pong := make(chan struct{})
	rxDelay := time.Duration(0)
	txDelay := time.Second * 5
	netDelay := time.Second * 1
	go handler.KeepAlive(conn, pong, rxDelay, txDelay, netDelay)
	pong <- struct{}{}

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
		case "KA":
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
