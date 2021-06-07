package deck

import (
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func initCards() []string {
	return []string{
		"{'color': 'yellow', 'type': '1'}",
		"{'color': 'green', 'type': '2'}",
		"{'color': 'red', 'type': '3'}",
		"{'color': 'blue', 'type': '4'}",
	}
}

func TestNew(t *testing.T) {
	_ = New(initCards)
}

func TestDeck_GetCurrentCard(t *testing.T) {
	testCases := []struct {
		name     string
		action   func(deck Deck)
		expected string
		error    error
	}{
		{
			name:   "returns nil for no card",
			action: func(deck Deck) {},
			error:  NoCardsErr,
		},
		{
			name: "returns card when it exists",
			action: func(deck Deck) {
				deck.Discard([]string{
					initCards()[0],
				})
			},
			expected: initCards()[0],
		},
	}

	assert := testify.New(t)
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			d := New(func() []string { return []string{} })
			tc.action(d)
			actual, err := d.GetCurrentCard()
			assert.Equal(tc.expected, actual)
			assert.Equal(tc.error, err)
		})
	}
}

func TestDeck_Reset(t *testing.T) {
	assert := testify.New(t)

	d := New(initCards)

	var err error
	var drawn []string

	testCases := []struct {
		name   string
		action func()
	}{
		{
			"exhaust deck",
			func() {
				drawn, err = d.Draw(4)
				assert.NoError(err)
				assert.Len(drawn, 4)
			},
		},
		{
			"discard should be empty",
			func() {
				_, err = d.GetCurrentCard()
				assert.EqualError(err, NoCardsErr.Error())
			},
		},
		{
			"discard some",
			func() {
				d.Discard(drawn)
				_, err = d.GetCurrentCard()
				assert.NoError(err)
			},
		},
		{
			"draw is empty, combined not enough to auto merge",
			func() {
				_, err = d.Draw(4)
				assert.EqualError(err, NotEnoughCardsErr.Error())
			},
		},
		{
			"manual reset",
			d.Reset,
		},
		{
			"discard should be empty again",
			func() {
				_, err = d.GetCurrentCard()
				assert.Error(NoCardsErr, err)
			},
		},
		{
			"draw should have cards again",
			func() {
				drawn, err = d.Draw(4)
				assert.NoError(err)
				assert.Len(drawn, 4)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.action()
		})
	}
}

func TestDeck_Discard(t *testing.T) {
	assert := testify.New(t)

	d := New(initCards)

	var err error
	var drawn []string

	testCases := []struct {
		name   string
		action func()
	}{
		{
			"discard should be empty",
			func() {
				_, err := d.GetCurrentCard()
				assert.EqualError(err, NoCardsErr.Error())
			},
		},
		{
			"discard one",
			func() {
				drawn, err = d.Draw(1)
				assert.NoError(err)
				d.Discard(drawn)
			},
		},
		{
			"discard should have last discarded card",
			func() {
				actual, err := d.GetCurrentCard()
				assert.NoError(err)
				assert.Equal(drawn[0], actual)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.action()
		})
	}
}

func TestDeck_Draw(t *testing.T) {
	assert := testify.New(t)

	d := New(initCards)

	var err error
	var drawn []string

	testCases := []struct {
		name   string
		action func()
	}{
		{
			"should not be able to draw more than available",
			func() {
				drawn, err = d.Draw(len(initCards()) + 1)
				assert.EqualError(err, NotEnoughCardsErr.Error())
				assert.Len(drawn, 0)
			},
		},
		{
			"draw enough when available",
			func() {
				drawn, err = d.Draw(4)
				assert.NoError(err)
				assert.Len(drawn, 4)
			},
		},
		{
			"deck should be exhausted, no auto merge with empty discard",
			func() {
				_drawn, err := d.Draw(1)
				assert.EqualError(err, NotEnoughCardsErr.Error())
				assert.Len(_drawn, 0)
			},
		},
		{
			"discard drawn",
			func() {
				d.Discard(drawn)
			},
		},
		{
			"should not auto merge too small discard",
			func() {
				drawn, err = d.Draw(4)
				assert.EqualError(err, NotEnoughCardsErr.Error())
				assert.Len(drawn, 0)
			},
		},
		{
			"should auto merge discard and draw",
			func() {
				drawn, err = d.Draw(3)
				assert.NoError(err)
				assert.Len(drawn, 3)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.action()
		})
	}
}
