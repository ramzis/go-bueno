package main

import (
	"net/url"

	"github.com/ramzis/bueno/internal/app/client/connection"
	"github.com/ramzis/bueno/internal/pkg/dialer"
)

func main() {
	dialer.DialTcp(url.URL{Host: "172.28.231.89:8080"}, connection.HandleConnection)
}
