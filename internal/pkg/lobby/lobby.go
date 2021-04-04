package lobby

import "log"

type lobby struct {
	c chan string
}

func (l *lobby) GetMessageChan() chan string {
	return l.c
}

func (l *lobby) PlayerJoin() {
}

func (l *lobby) PlayerLeave() {
}

func (l *lobby) Read() {
	go func() {
		for {
			select {
			case msg := <-l.c:
				log.Printf("LOBBY got: %s", msg)
			}
		}
	}()
}

func (l *lobby) Write(msg string) {
	l.c <- msg
}

func New() Lobby {
	l := &lobby{
		c: make(chan string, 0),
	}
	l.Read()
	return l
}
