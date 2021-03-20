package connection

import (
	handler "github.com/ramzis/bueno/internal/pkg/connection"
	"log"
	"net"
	"strings"
)

func HandleConnection(conn net.Conn) {
	r, w, e := handler.HandleConnection(conn, false)

	input := make(chan string)
	go handler.ReadInput(input)

	for {
		select {
		case line := <-input:
			w <- line
		case cmd := <-r:
			cmds := strings.Split(cmd, " ")
			if len(cmds) > 1 {
				from := cmds[0]
				msg := strings.Join(cmds[1:], " ")
				log.Printf("[%s]: %s", from, msg)
			}
		case err := <-e:
			log.Println(err)
			return
		}
	}
}
