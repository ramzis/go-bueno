package main

import (
	"net/url"

	"github.com/ramzis/bueno/internal/app/server/bueno"
	"github.com/ramzis/bueno/internal/pkg/server"
)

func main() {
	bueno := bueno.NewBueno()
	server.Listen(url.URL{Host: ":8080"}, bueno.HandleConnection)
}
