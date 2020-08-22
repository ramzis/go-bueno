package bueno

import (
	"bufio"
	"net"
)

type Bueno struct {
	conns     []net.Conn
	listeners []func()
}

func NewBueno() *Bueno {
	return &Bueno{conns: make([]net.Conn, 0)}
}

func (b *Bueno) RegisterConn(conn net.Conn) {
	b.conns = append(b.conns, conn)
	b.OnUpdateConns()
}

func (b *Bueno) RemoveConn(conn net.Conn) {
	for i := len(b.conns) - 1; i >= 0; i-- {
		if b.conns[i] == conn {
			b.conns = append(b.conns[:i], b.conns[i+1:]...)
			break
		}
	}
	b.OnUpdateConns()
}

func (b *Bueno) OnUpdateConns() {
	for _, notify := range b.listeners {
		notify()
	}
}

func (b *Bueno) RegisterConnListener(notify func()) {
	b.listeners = append(b.listeners, notify)
}

func (b *Bueno) TellEveryone(s string) {
	for _, c := range b.conns {
		w := bufio.NewWriter(c)
		w.WriteString(s)
		w.WriteByte(0x0)
		w.Flush()
	}
}
