package game

import "github.com/ramzis/bueno/internal/pkg/lobby/entity"

type Game interface {
	AddPlayer(entity entity.Entity)
	Send(ID entity.ID, msg string)
	Receive(ID entity.ID, msg string)
}

type Game2Room interface {
	Write(msg string)
}

type MessageHandler interface {
	Handle(ID entity.ID, msg string)
}
