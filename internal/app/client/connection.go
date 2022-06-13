package client

import (
	handler "github.com/ramzis/bueno/internal/pkg/connection"
	"log"
	"net"
	"strings"
)

func HandleConnection(conn net.Conn) {
	c := handler.HandleConnection(conn, false)

	input := make(chan string)
	go handler.ReadInput(input)

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
