package bueno

import (
	"github.com/ramzis/bueno/internal/pkg/lobby/entity"
	"github.com/ramzis/bueno/internal/pkg/lobby/room/game"
	"github.com/ramzis/bueno/internal/pkg/lobby/room/game/deck"
	"log"
)

type bueno struct {
	game.Game
	state                State
	deck                 deck.Deck
	pendingCardsToDraw   int
	isNextPlayerSkipped  bool
	didPlayerRequestDraw bool
}

func (b *bueno) Start() {
	if b.state != nil {
		return
	}
	b.state = NewState()
	b.stateMachine()
}

func (b *bueno) Stop() {
	if b.state == nil {
		return
	}
	b.state.Set(STATE_Stop)
	b.state = nil
}

func (b *bueno) stateMachine() {
	go func() {
		stateChange := b.state.OnChangeChan()
		for {
			state, ok := <-stateChange
			if !ok {
				return
			}
			switch state {
			case STATE_NewGame:
				b.state.Set(STATE_InitVars)
			case STATE_InitVars:
				b.pendingCardsToDraw = 0
				b.isNextPlayerSkipped = false
				b.didPlayerRequestDraw = false
				b.state.Set(STATE_InitDeck)
			case STATE_InitDeck:
				b.deck.Reset()
				b.state.Set(STATE_DealStartingHands)
			case STATE_DealStartingHands:
				b.dealStartingHands()
			case STATE_DrawFirstCard:
				b.drawFirstCard()
				b.state.Set(STATE_StartGame)
			case STATE_StartGame:
				b.gameLoop()
				b.state.Set(STATE_GameRunning)
			case STATE_GameRunning:
			case STATE_GameRunning_WaitingForMove:
			case STATE_GameRunning_MakingMove:
			case STATE_GameRunning_ResolvingMoveEffects:
			case STATE_Null:
			default:
			}
		}
	}()
}

func (b *bueno) dealStartingHands() {
	panic("Not implemented")
}

func (b *bueno) drawFirstCard() {
	panic("Not implemented")
}

func (b *bueno) gameLoop() {
	panic("Not implemented")
}

func (b *bueno) Handle(ID entity.ID, msg string) {
	log.Printf("[BUENO] Handled message for %s: %s", ID, msg)
}

func New(room game.Game2Room) Bueno {
	return &bueno{
		Game:                 game.New(room, GAME_TYPE_BUENO),
		state:                nil,
		deck:                 deck.New(initDeck),
		pendingCardsToDraw:   0,
		isNextPlayerSkipped:  false,
		didPlayerRequestDraw: false,
	}
}
