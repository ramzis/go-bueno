package server

import (
	"github.com/ramzis/bueno/internal/pkg/connection"
	"github.com/ramzis/bueno/internal/pkg/lobby"
	"github.com/ramzis/bueno/internal/pkg/lobby/entity"
	"net"
)

type server struct {
	conns     map[string]*connection.Connection
	listeners []func()
	lobby     lobby.Lobby
	resolver  map[entity.ID]string
}

type Server interface {
	HandleConnection(conn net.Conn)
	RegisterConn(id string, conn *connection.Connection)
}

func New(lobby lobby.Lobby) Server {
	s := &server{
		conns:     make(map[string]*connection.Connection, 0),
		listeners: make([]func(), 0),
		lobby:     lobby,
		resolver:  make(map[entity.ID]string),
	}
	s.HandleLobbyMessages()
	return s
}

func (s *server) RegisterConn(id string, conn *connection.Connection) {
	s.conns[id] = conn
	s.OnUpdateConns()
}

func (s *server) RemoveConn(id string) {
	if _, ok := s.conns[id]; ok {
		delete(s.conns, id)
	}
	s.OnUpdateConns()
}

func (s *server) OnUpdateConns() {
	for _, notify := range s.listeners {
		notify()
	}
}

func (s *server) RegisterConnListener(notify func()) {
	s.listeners = append(s.listeners, notify)
}

func (s *server) TellEveryone(msg string) {
	for _, conn := range s.conns {
		conn.W <- msg
	}
}

func (s *server) TellEveryoneBut(id string, msg string) {
	for connId, conn := range s.conns {
		if connId == id {
			continue
		}
		conn.W <- msg
	}
}

func (s *server) TellOne(id string, msg string) {
	for connId, conn := range s.conns {
		if connId == id {
			conn.W <- msg
			break
		}
	}
}
