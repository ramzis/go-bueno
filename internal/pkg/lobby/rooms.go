package lobby

import (
	"fmt"
	"github.com/ramzis/bueno/internal/pkg/lobby/entity"
	"github.com/ramzis/bueno/internal/pkg/lobby/room"
	"log"
)

func (l *lobby) GetDefaultRoom() room.Room {
	return l.defaultRoom
}

// TODO: make into a chan
func (l *lobby) SendMessageToRoomAll(from entity.ID, room room.ID, msg string) {
	r, ok := l.rooms[room]
	if !ok {
		log.Println("Invalid room")
		return
	}
	for _, to := range r.GetEntities() {
		if to == from {
			continue
		}
		l.write(fmt.Sprintf("E %s R %s MSG %s %s", from, room, to, msg))
	}
}

func (l *lobby) SendMessageToRoomOne(from entity.ID, to entity.ID, room room.ID, msg string) {
	r, ok := l.rooms[room]
	if !ok {
		log.Println("Invalid room")
		return
	}
	for _, _to := range r.GetEntities() {
		if _to == to {
			l.write(fmt.Sprintf("E %s R %s MSG %s %s", from, room, to, msg))
			break
		}
	}
}
