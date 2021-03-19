package connection

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

// HandleConnection ...
func HandleConnection(conn net.Conn) {
	w := bufio.NewWriter(conn)

	ping := make(chan struct{})

	msgChan := make(chan string)
	go ReadConn(conn, msgChan)

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

		// Trim 0x0
		s = s[:len(s)-1]

		log.Println("Client received", s)
		
		cmd := strings.Split(s, " ")
		if len(cmd) < 1 {
			log.Println("Invalid cmd received", cmd)
			continue
		}

		switch {
		case cmd[0] == "HI":
			w.WriteString("HI")
			w.WriteByte(0x0)
			w.Flush()
			go KeepAlive(conn, ping)
			go SendMessage(w)
		case cmd[0] == "PING":
			ping <- struct{}{}
		case cmd[0] == "MSG":
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

func ReadConn(conn net.Conn, msg chan string) {
	r := bufio.NewReader(conn)

	for {
		s, err := r.ReadString(byte(0x0))
		if err != nil {
			log.Println("Client received error", s, err.Error())
			close(msg)
			return
		}
		msg <- s
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
