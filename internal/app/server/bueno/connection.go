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
	// conn.SetDeadline(time.Now().Add(time.Second * 3))
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	// Write greeting
	w.WriteString("HI")
	w.WriteByte(0x0)
	w.Flush()

	// Wait for HI response or timeout
	responded := false
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))

	// Wait for cmd
	for {
		s, err := r.ReadString(byte(0x0))
		if err != nil {
			log.Println("Server received error", s, err.Error())
			break
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
			conn.SetReadDeadline(time.Time{})
			b.RegisterConn(conn)
			defer b.RemoveConn(conn)
		case cmd[0] == "MSG":
			if len(cmd) > 1 {
				msg := strings.Join(cmd[1:], " ")
				b.TellEveryone(fmt.Sprintf("MSG %s %s", conn.RemoteAddr().String(), msg))
			}
		}
	}
}
