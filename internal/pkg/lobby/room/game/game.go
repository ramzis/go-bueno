package game

import (
	"github.com/ramzis/bueno/internal/pkg/lobby/entity"
	"log"
)

type game struct {
	Type    Type
	room    Game2Room
	players map[entity.ID]entity.Entity
}

func (g *game) StartGame() {
}

func (g *game) PlayMove(entity entity.Entity) {
}

func (g *game) AddPlayer(entity entity.Entity) {
	g.players[entity.GetID()] = entity
}

func (g *game) SendMessageToRoomOne(to entity.ID, msg string) {
	g.room.SendMessageToRoomOne(to, msg)
}

func (g *game) SendMessageToRoomAll(msg string) {
	g.room.SendMessageToRoomAll(msg)
}

func (g *game) Receive(ID entity.ID, msg string) {
	g.HandleMessage(ID, msg)
}

func (g *game) HandleMessage(ID entity.ID, msg string) {
	log.Println("Unhandled message in game")
}

func New(room Game2Room, t Type) Game {
	return &game{
		Type:    t,
		room:    room,
		players: make(map[entity.ID]entity.Entity, 0),
	}
}
