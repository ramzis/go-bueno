package bueno

import (
	"fmt"
	handler "github.com/ramzis/bueno/internal/pkg/connection"
	"log"
	"net"
)

func (b *Bueno) HandleConnection(conn net.Conn) {
	c := handler.HandleConnection(conn, true)

	id := conn.RemoteAddr().String()

	defer b.RemoveConn(id)
	b.RegisterConn(id, c)

	defer b.TellEveryoneBut(id, fmt.Sprintf("%s %s has disconnected", "Server", id))
	b.TellEveryoneBut(id, fmt.Sprintf("%s %s has connected", "Server", id))

	for {
		select {
		case cmd := <-c.R:
			log.Println("Server got", cmd)
			b.TellEveryoneBut(id, fmt.Sprintf("%s %s", id, cmd))
		case err := <-c.E:
			log.Println(err)
			return
		}
	}
}
