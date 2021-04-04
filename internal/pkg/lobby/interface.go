package lobby

type Lobby interface {
	PlayerJoin()
	PlayerLeave()
	GetMessageChan() chan string
}
