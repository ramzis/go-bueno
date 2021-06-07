package deck

import (
	"github.com/pkg/errors"
	"math/rand"
	"time"
)

var NoCardsErr = errors.New("no cards")
var NotEnoughCardsErr = errors.New("not enough cards")

type deck struct {
	draw      []string
	discard   []string
	initCards func() []string
}

func (d *deck) Draw(count int) ([]string, error) {
	// If there are not enough cards to draw.
	if len(d.draw) < count {
		// If it is not possible to merge discard to draw.
		if (len(d.discard) - 1 + len(d.draw)) < count {
			return []string{}, NotEnoughCardsErr
		}

		// Merge all but last from discard into draw and shuffle.
		d.draw = append(d.draw, d.discard[:len(d.discard)-1]...)
		d.discard = d.discard[len(d.discard)-1:]
		d.shuffleDraw()
	}

	drawnCards := d.draw[:count]
	d.draw = d.draw[count:]

	return drawnCards, nil
}

func (d *deck) initDraw() {
	d.draw = d.initCards()
}

func (d *deck) initDiscard() {
	d.discard = make([]string, 0)
}

func (d *deck) shuffleDraw() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.draw), func(i, j int) { d.draw[i], d.draw[j] = d.draw[j], d.draw[i] })
}

func (d *deck) Discard(cards []string) {
	d.discard = append(d.discard, cards...)
}

func (d *deck) Reset() {
	d.initDraw()
	d.initDiscard()
	d.shuffleDraw()
}

func (d *deck) GetCurrentCard() (string, error) {
	if cardsLeft := len(d.discard); cardsLeft < 1 {
		return "", NoCardsErr
	} else {
		return d.discard[cardsLeft-1], nil
	}
}

func New(initCards func() []string) Deck {
	deck := &deck{
		initCards: initCards,
	}
	deck.Reset()
	return deck
}
