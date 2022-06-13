package main

import (
	"github.com/ramzis/bueno/internal/app/client"
	"github.com/ramzis/bueno/internal/pkg/connection"
	"net/url"
)

func main() {
	connection.DialTcp(url.URL{Host: ":8080"}, client.HandleConnection)
}
