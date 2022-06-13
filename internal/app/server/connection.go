package server

import (
	"fmt"
	"github.com/ramzis/bueno/internal/pkg/connection"
	"log"
	"net"
)

func (s *server) HandleConnection(conn net.Conn) {
	c := connection.HandleConnection(conn, true)

	id := conn.RemoteAddr().String()

	defer s.RemoveConn(id)
	s.RegisterConn(id, c)

	entityID := s.lobby.Join()
	defer s.lobby.Leave(entityID)

	s.resolver[entityID] = id
	defer delete(s.resolver, entityID)

	defer s.TellEveryoneBut(id, fmt.Sprintf("%s %s has disconnected", "Server", id))
	s.TellEveryoneBut(id, fmt.Sprintf("%s %s has connected", "Server", id))

	for {
		select {
		case cmd := <-c.R:
			log.Println("Server got from conn", cmd)
			//s.TellEveryoneBut(id, fmt.Sprintf("%s %s", id, cmd))
			s.lobby.Handle(entityID, cmd)
		case err := <-c.E:
			log.Println(err)
			return
		}
	}
}
