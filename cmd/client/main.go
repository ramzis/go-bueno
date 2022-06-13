package main

import (
	"github.com/ramzis/bueno/internal/app/client"
	"github.com/ramzis/bueno/internal/pkg/connection"
	"net/url"
)

func main() {
	connection.DialTcp(url.URL{Host: "172.28.231.89:8080"}, client.HandleConnection)
}
