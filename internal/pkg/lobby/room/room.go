package room

import (
	"github.com/google/uuid"
	"github.com/ramzis/bueno/internal/pkg/lobby/entity"
	"github.com/ramzis/bueno/internal/pkg/lobby/room/game"
	"log"
	"strings"
)

type room struct {
	ID       ID
	lobby    Room2Lobby
	game     game.Game
	entities map[entity.ID]entity.Entity
}

func (r *room) Write(msg string) {
	r.lobby.Write(msg)
}

func (r *room) Join(entity entity.Entity) {
	r.entities[entity.GetID()] = entity
	log.Printf("%s has joined the room %s", entity.GetID(), r.ID)
}

func (r *room) Leave(ID entity.ID) bool {
	// TODO: remove from games
	left := false
	if _, found := r.entities[ID]; found {
		delete(r.entities, ID)
		left = found
	}
	return left
}

func (r *room) GetID() ID {
	return r.ID
}

func (r *room) GetEntities() []entity.ID {
	entities := make([]entity.ID, 0, len(r.entities))
	for _, e := range r.entities {
		entities = append(entities, e.GetID())
	}
	return entities
}

func New(lobby Room2Lobby) Room {
	return &room{
		ID:       ID(strings.Split(uuid.NewString(), "-")[0]),
		lobby:    lobby,
		entities: make(map[entity.ID]entity.Entity),
	}
}
