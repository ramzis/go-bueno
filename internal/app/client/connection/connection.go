package connection

import (
	handler "github.com/ramzis/bueno/internal/pkg/connection"
	"log"
	"net"
)

func HandleConnection(conn net.Conn) {
	r, w, e := handler.HandleConnection(conn, false)

	for {
		select {
		case <-r:
			w <- "Test"
			//if len(cmd) > 2 {
			//	from := cmd[1]
			//	msg := strings.Join(cmd[2:], " ")
			//	log.Printf("[%s]: %s", from, msg)
			//}
		case err := <-e:
			log.Println(err)
			return
		}
	}
}
