package client

import (
	"github.com/ramzis/bueno/internal/pkg/connection"
	"log"
	"net"
	"strings"
)

func HandleConnection(conn net.Conn) {
	c := connection.HandleConnection(conn, false)

	input := make(chan string)
	go connection.ReadInput(input)

	for {
		select {
		case line := <-input:
			c.W <- line
		case cmd := <-c.R:
			cmds := strings.Split(cmd, " ")
			if len(cmds) > 1 {
				from := cmds[0]
				msg := strings.Join(cmds[1:], " ")
				log.Printf("[%s]: %s", from, msg)
			}
		case err := <-c.E:
			log.Println(err)
			return
		}
	}
}
