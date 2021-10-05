package dialer

import (
	"log"
	"net"
	"net/url"
)

func DialTcp(URL url.URL, connectionHandler func(net.Conn)) {
	log.Println("Connecting to server...")
	conn, err := net.Dial("tcp", URL.Host)
	if err != nil {
		log.Println("Failed to establish connection to", URL.Host, err)
		return
	}

	defer log.Println("Connection to", conn.RemoteAddr(), "ended")
	log.Println("Connected to", conn.RemoteAddr())
	connectionHandler(conn)
}
