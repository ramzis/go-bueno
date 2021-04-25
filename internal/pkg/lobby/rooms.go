package lobby

import "github.com/ramzis/bueno/internal/pkg/lobby/room"

func (l *lobby) GetDefaultRoom() room.Room {
	return l.defaultRoom
}
