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
		case cmd := <-r:
			log.Println("Server got", cmd)
			continue
			w <- "Test"
			//defer b.RemoveConn(conn)
			//b.RegisterConn(conn)
			//b.TellEveryone(fmt.Sprintf("MSG %s %s", conn.RemoteAddr().String(), msg))
		case err := <-e:
			log.Println(err)
			return
		}
	}
}
