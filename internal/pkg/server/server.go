package server

import (
	"log"
	"net"
	"net/url"
)

func Listen(URL url.URL, connectionHandler func(net.Conn)) {
	ln, err := net.Listen("tcp", URL.Host)
	if err != nil {
		panic(err)
	}
	log.Print("Server listening on ", ln.Addr())
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Failed to handle connection", err)
			continue
		}

		go func() {
			defer conn.Close()
			defer log.Println("Connection to", conn.RemoteAddr(), "ended")
			log.Println(conn.RemoteAddr(), "connected")
			connectionHandler(conn)
		}()
	}
}
