package lobby

import (
	"fmt"
	"github.com/ramzis/bueno/internal/pkg/lobby/entity"
	"github.com/ramzis/bueno/internal/pkg/lobby/room"
	"log"
	"strings"
)

type lobby struct {
	c           chan string
	rooms       map[room.ID]room.Room
	defaultRoom room.Room
}

func (l *lobby) GetMessageChan() chan string {
	return l.c
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
		l.Write(fmt.Sprintf("E %s R %s MSG %s %s", from, room, to, msg))
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
			l.Write(fmt.Sprintf("E %s R %s MSG %s %s", from, room, to, msg))
			break
		}
	}
}

func (l *lobby) Handle(from entity.ID, msg string) {
	words := strings.Split(msg, " ")
	if len(words) < 7 {
		log.Println("Invalid format length")
		return
	}
	if words[0] != "E" {
		log.Println("Invalid format E")
		return
	}
	e := words[1]
	if entity.ID(e) != from {
		log.Println("Invalid sender")
		return
	}
	if words[2] != "R" {
		log.Println("Invalid format R")
		return
	}
	r := words[3]
	if words[4] != "MSG" {
		log.Printf("Unhandled command type %s", words[4])
		return
	}
	switch words[5] {
	case "ALL":
		l.SendMessageToRoomAll(entity.ID(e), room.ID(r), strings.Join(words[6:], " "))
	default:
		l.SendMessageToRoomOne(entity.ID(e), entity.ID(words[5]), room.ID(r), strings.Join(words[6:], " "))
	}
}

func (l *lobby) Write(msg string) {
	l.c <- msg
}

func New() Lobby {
	l := &lobby{
		c:     make(chan string, 0),
		rooms: make(map[room.ID]room.Room, 0),
	}
	l.defaultRoom = room.New(l)
	l.rooms[l.defaultRoom.GetID()] = l.defaultRoom
	return l
}
