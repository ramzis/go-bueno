package server

import (
	"github.com/ramzis/bueno/internal/pkg/connection"
	"github.com/ramzis/bueno/internal/pkg/lobby"
	"github.com/ramzis/bueno/internal/pkg/lobby/entity"
	"net"
)

type server struct {
	conns     map[string]*handler.Connection
	listeners []func()
	lobby     lobby.Lobby
	resolver  map[entity.ID]string
}

type Server interface {
	HandleConnection(conn net.Conn)
	RegisterConn(id string, conn *handler.Connection)
}

func New(lobby lobby.Lobby) Server {
	s := &server{
		conns:     make(map[string]*handler.Connection, 0),
		listeners: make([]func(), 0),
		lobby:     lobby,
		resolver:  make(map[entity.ID]string),
	}
	s.HandleLobbyMessages()
	return s
}

func (b *server) RegisterConn(id string, conn *handler.Connection) {
	b.conns[id] = conn
	b.OnUpdateConns()
}

func (b *server) RemoveConn(id string) {
	if _, ok := b.conns[id]; ok {
		delete(b.conns, id)
	}
	b.OnUpdateConns()
}

func (b *server) OnUpdateConns() {
	for _, notify := range b.listeners {
		notify()
	}
}

func (b *server) RegisterConnListener(notify func()) {
	b.listeners = append(b.listeners, notify)
}

func (b *server) TellEveryone(s string) {
	for _, conn := range b.conns {
		conn.W <- s
	}
}

func (b *server) TellEveryoneBut(id string, s string) {
	for connId, conn := range b.conns {
		if connId == id {
			continue
		}
		conn.W <- s
	}
}

func (b *server) TellOne(id string, s string) {
	for connId, conn := range b.conns {
		if connId == id {
			conn.W <- s
			break
		}
	}
}
