package room

import (
	"github.com/ramzis/bueno/internal/pkg/lobby/entity"
)

type Room interface {
	// Adds an entity to a room
	Join(entity entity.Entity)
	// Removes an entity from a room
	Leave(ID entity.ID) bool
	GetID() ID
	Write(string)
	GetEntities() []entity.ID
}

type Room2Lobby interface {
	SendMessageToRoomAll(from entity.ID, room ID, msg string)
	SendMessageToRoomOne(from entity.ID, to entity.ID, room ID, msg string)
}
