package lobby

type lobby struct {
}

func (l *lobby) PlayerJoin() {
}

func (l *lobby) PlayerLeave() {
}

func New() Lobby {
	return &lobby{}
}
