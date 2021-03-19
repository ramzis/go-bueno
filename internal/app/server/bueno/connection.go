package bueno

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func (b *Bueno) HandleConnection(conn net.Conn) {
	w := bufio.NewWriter(conn)

	pong := make(chan struct{})

	msgChan := make(chan string)
	go b.ReadConn(conn, msgChan)

	// Write greeting
	w.WriteString("HI")
	w.WriteByte(0x0)
	w.Flush()

	// Wait for HI response or timeout
	responded := false
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

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

		// Trim 0x0
		s = s[:len(s)-1]

		log.Println("Server received", s)

		cmd := strings.Split(s, " ")
		if len(cmd) < 1 {
			log.Println("Invalid cmd received", cmd)
			continue
		}

		switch {
		case cmd[0] == "HI":
			if responded {
				continue
			}
			responded = true
			_ = conn.SetReadDeadline(time.Time{})
			b.RegisterConn(conn)
			go b.KeepAlive(conn, pong)
			//goland:noinspection ALL
			defer b.RemoveConn(conn)
		case cmd[0] == "PONG":
			pong <- struct{}{}
		case cmd[0] == "MSG":
			if len(cmd) > 1 {
				msg := strings.Join(cmd[1:], " ")
				b.TellEveryone(fmt.Sprintf("MSG %s %s", conn.RemoteAddr().String(), msg))
			}
		default:
			log.Println(s, "is unhandled")
		}
	}
}

func (b *Bueno) ReadConn(conn net.Conn, msg chan string) {
	r := bufio.NewReader(conn)

	for {
		s, err := r.ReadString(byte(0x0))
		if err != nil {
			log.Println("Server received error", s, err.Error())
			close(msg)
			return
		}
		msg <- s
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
