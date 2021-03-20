package bueno

import (
	"fmt"
	handler "github.com/ramzis/bueno/internal/pkg/connection"
	"log"
	"net"
)

func (b *Bueno) HandleConnection(conn net.Conn) {
	r, w, e := handler.HandleConnection(conn, true)

	defer b.RemoveConn(conn.RemoteAddr().String())
	b.RegisterConn(conn.RemoteAddr().String(), Connection{r, w, e})

	for {
		select {
		case cmd := <-r:
			log.Println("Server got", cmd)
			b.TellEveryone(fmt.Sprintf("%s %s", conn.RemoteAddr().String(), cmd))
		case err := <-e:
			log.Println(err)
			return
		}
	}
}
