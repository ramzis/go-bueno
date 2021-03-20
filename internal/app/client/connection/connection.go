package connection

import (
	handler "github.com/ramzis/bueno/internal/pkg/connection"
	"log"
	"net"
)

func HandleConnection(conn net.Conn) {
	r, w, e := handler.HandleConnection(conn, false)

	for {
		select {
		case <- r:
			w <- "Test"
		case <- e:
			log.Println(e)
			return
		}
	}
}