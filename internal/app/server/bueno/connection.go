package bueno

import (
	"fmt"
	handler "github.com/ramzis/bueno/internal/pkg/connection"
	"log"
	"net"
)

func (b *Bueno) HandleConnection(conn net.Conn) {
	r, w, e := handler.HandleConnection(conn, true)

	id := conn.RemoteAddr().String()

	defer b.RemoveConn(id)
	b.RegisterConn(id, Connection{r, w, e})

	defer b.TellEveryoneBut(id, fmt.Sprintf("%s %s has disconnected", "Server", id))
	b.TellEveryoneBut(id, fmt.Sprintf("%s %s has connected", "Server", id))

	for {
		select {
		case cmd := <-r:
			log.Println("Server got", cmd)
			b.TellEveryoneBut(id, fmt.Sprintf("%s %s", id, cmd))
		case err := <-e:
			log.Println(err)
			return
		}
	}
}
