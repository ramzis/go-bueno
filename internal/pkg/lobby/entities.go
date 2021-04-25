package lobby

import (
	"github.com/ramzis/bueno/internal/pkg/lobby/entity"
)

func (l *lobby) Join() entity.ID {
	e := entity.New()
	ID := e.GetID()
	l.defaultRoom.Join(e)
	return ID
}

func (l *lobby) Leave(ID entity.ID) {
	for _, room := range l.rooms {
		if left := room.Leave(ID); left {
			break
		}
	}
}
