package game

import "github.com/ramzis/bueno/internal/pkg/lobby/entity"

type Game interface {
	AddPlayer(entity entity.Entity)
	MessageHandler
	Game2Room
}

type Game2Room interface {
	SendMessageToRoomOne(to entity.ID, msg string)
	SendMessageToRoomAll(msg string)
}

type MessageHandler interface {
	Handle(ID entity.ID, msg string)
}
