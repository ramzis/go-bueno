package connection

import (
	"bufio"
	"fmt"
	handler "github.com/ramzis/bueno/internal/pkg/connection"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// HandleConnection ...
func HandleConnection(conn net.Conn) {

	err := handler.PerformHandshake(conn, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	w := bufio.NewWriter(conn)
	go SendMessage(w)

	ping := make(chan struct{})
	go KeepAlive(conn, ping)

	msgChan := make(chan string)
	go handler.ReadConn(conn, msgChan)

	// Wait for cmd or failed ping
	var s string
	var ok bool
	for {
		select {
		case _, ok := <-ping:
			if !ok {
				log.Println("Ping channel closed but shouldn't be!")
				return
			}
			close(ping)
			log.Println("Client failed to receive PING")
			return
		case s, ok = <-msgChan:
			log.Println("Client reading...")
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
		case "PING":
			ping <- struct{}{}
		case "MSG":
			if len(cmd) > 2 {
				from := cmd[1]
				msg := strings.Join(cmd[2:], " ")
				log.Printf("[%s]: %s", from, msg)
			}
		default:
			log.Println(s, "is unhandled")
		}
	}
}

func SendMessage(w *bufio.Writer) {
	r := bufio.NewReader(os.Stdin)
	for {
		msg, _ := r.ReadString('\n')
		// convert CRLF to LF
		msg = strings.Replace(msg, "\n", "", -1)
		w.WriteString("MSG ")
		w.WriteString(msg)
		w.WriteByte(0x0)
		w.Flush()
	}
}

func KeepAlive(conn net.Conn, ping chan struct{}) {
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()
	w := bufio.NewWriter(conn)

	writePong := func() {
		ticker.Stop()
		log.Println("Sending PONG")
		w.WriteString("PONG")
		w.WriteByte(0x0)
		w.Flush()
		ticker.Reset(time.Second * 6)
	}

	for {
		select {
		case <- ticker.C:
			ping <- struct{}{}
			return
		case <- ping:
			log.Println("Received PING")
			writePong()
		}
	}
}
