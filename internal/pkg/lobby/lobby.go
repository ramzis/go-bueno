package lobby

import (
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

// Handle reads an incoming message from the connection handler
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

// write sends an outgoing message to the connection handler
func (l *lobby) write(msg string) {
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
