package server

import (
	"fmt"
	"github.com/ramzis/bueno/internal/pkg/lobby/entity"
	"log"
	"strings"
)

func (b *server) HandleLobbyMessages() {
	go func() {
		lobbyMsgChan := b.lobby.GetMessageChan()

		for {
			select {
			case lobbyMsg, ok := <-lobbyMsgChan:
				if !ok {
					log.Println("lobbyMsgChan reading closed by server")
					return
				}
				log.Println("Server got from lobby", lobbyMsg)

				split := strings.Split(lobbyMsg, " ")
				if len(split) < 7 {
					log.Println("Invalid length in lobby msg chan handler")
					continue
				}
				to := b.resolver[entity.ID(split[5])]
				b.TellOne(to, fmt.Sprintf("%s@%s %s", split[1], split[3], strings.Join(split[6:], " ")))
			}
		}
	}()
}
