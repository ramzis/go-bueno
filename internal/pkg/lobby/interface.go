package lobby

import "github.com/ramzis/bueno/internal/pkg/lobby/entity"

// The Lobby is the message broker for all the room.Room's
// By default, there exists the 'lobby' room.Room which can be used to join
// other ones made for games.
type Lobby interface {
	Join() entity.ID
	Leave(ID entity.ID)
	// Used to notify the lobby and read its messages
	GetMessageChan() chan string
	Handle(ID entity.ID, msg string)
}
