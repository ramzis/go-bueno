package room

import "github.com/ramzis/bueno/internal/pkg/lobby/room/game/bueno"

func (r *room) CreateGame() {
	r.game = bueno.New(r)
}
