package main

import (
	"github.com/ramzis/bueno/internal/app/server"
	"github.com/ramzis/bueno/internal/pkg/connection"
	"github.com/ramzis/bueno/internal/pkg/lobby"
	"net/url"
)

func main() {
	lobby := lobby.New()
	server := server.New(lobby)
	connection.Listen(url.URL{Host: ":8080"}, server.HandleConnection)
}
