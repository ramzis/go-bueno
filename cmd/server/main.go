package main

import (
	"github.com/ramzis/bueno/internal/app/server"
	"github.com/ramzis/bueno/internal/pkg/lobby"
	"net/url"

	listener "github.com/ramzis/bueno/internal/pkg/server"
)

func main() {
	lobby := lobby.New()
	server := server.New(lobby)
	listener.Listen(url.URL{Host: ":8080"}, server.HandleConnection)
}
