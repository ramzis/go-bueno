package bueno

type Connection struct {
	R, W, E chan string
}

type Bueno struct {
	conns     map[string]Connection
	listeners []func()
}

func NewBueno() *Bueno {
	return &Bueno{conns: make(map[string]Connection, 0)}
}

func (b *Bueno) RegisterConn(id string, conn Connection) {
	b.conns[id] = conn
	b.OnUpdateConns()
}

func (b *Bueno) RemoveConn(id string) {
	if _, ok := b.conns[id]; ok {
		delete(b.conns, id)
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
	for _, conn := range b.conns {
		conn.W <- s
	}
}

func (b *Bueno) TellEveryoneBut(id string, s string) {
	for connId, conn := range b.conns {
		if connId == id {
			continue
		}
		conn.W <- s
	}
}
