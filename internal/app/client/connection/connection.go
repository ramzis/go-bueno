package connection

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

// HandleConnection ...
func HandleConnection(conn net.Conn) {
	// conn.SetDeadline(time.Now().Add(time.Second * 5))
	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	// Wait for cmd
	for {
		s, err := r.ReadString(byte(0x0))
		if err != nil {
			log.Println("Client received error", s, err.Error())
			break
		}

		// Trim 0x0
		s = s[:len(s)-1]

		// log.Println("Client received", s)

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
			go SendMessage(w)
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
