package bueno

import (
	handler "github.com/ramzis/bueno/internal/pkg/connection"
	"log"
	"net"
)

func (b *Bueno) HandleConnection(conn net.Conn) {
	r, w, e := handler.HandleConnection(conn, true)

	for {
		select {
			case <- r:
				w <- "Test"
				//defer b.RemoveConn(conn)
				//b.RegisterConn(conn)
			case <- e:
				log.Println(e)
				return
		}
	}
}