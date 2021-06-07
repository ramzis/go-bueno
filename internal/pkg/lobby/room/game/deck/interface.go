package deck

type Deck interface {
	Draw(count int) ([]string, error)
	Discard(cards []string)
	Reset()
	GetCurrentCard() (string, error)
}
