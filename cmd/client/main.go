package main

import (
	"github.com/ramzis/bueno/internal/app/client"
	"github.com/ramzis/bueno/internal/pkg/dialer"
	"net/url"
)

func main() {
	dialer.DialTcp(url.URL{Host: "172.28.231.89:8080"}, client.HandleConnection)
}
