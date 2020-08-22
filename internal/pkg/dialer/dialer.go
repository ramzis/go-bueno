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
		log.Println("Failed to establish connection to", URL.Host)
	} else {
		log.Println("Connected to", conn.RemoteAddr().String())
		finished := make(chan bool)
		go func() {
			defer log.Println("Connection to", conn.RemoteAddr().String(), "ended")
			connectionHandler(conn)
			finished <- true
		}()
		<-finished
	}
}
