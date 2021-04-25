package bueno

import (
	"github.com/ramzis/bueno/internal/pkg/lobby/entity"
	"github.com/ramzis/bueno/internal/pkg/lobby/room/game"
	"log"
)

type bueno struct {
	game.Game
}

func (b *bueno) HandleMessage(ID entity.ID, msg string) {
	log.Printf("[BUENO] Handled message for %s: %s", ID, msg)
}

func New(room game.Game2Room) Bueno {
	return &bueno{
		Game: game.New(room, GAME_TYPE_BUENO),
	}
}
