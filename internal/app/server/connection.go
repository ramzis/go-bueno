package server

import (
	"fmt"
	handler "github.com/ramzis/bueno/internal/pkg/connection"
	"log"
	"net"
)

func (b *server) HandleConnection(conn net.Conn) {
	c := handler.HandleConnection(conn, true)

	id := conn.RemoteAddr().String()

	defer b.RemoveConn(id)
	b.RegisterConn(id, c)

	entityID := b.lobby.Join()
	defer b.lobby.Leave(entityID)

	b.resolver[entityID] = id
	defer delete(b.resolver, entityID)

	defer b.TellEveryoneBut(id, fmt.Sprintf("%s %s has disconnected", "Server", id))
	b.TellEveryoneBut(id, fmt.Sprintf("%s %s has connected", "Server", id))

	for {
		select {
		case cmd := <-c.R:
			log.Println("Server got from conn", cmd)
			//b.TellEveryoneBut(id, fmt.Sprintf("%s %s", id, cmd))
			b.lobby.Handle(entityID, cmd)
		case err := <-c.E:
			log.Println(err)
			return
		}
	}
}
