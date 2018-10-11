package handy

// Deck is a bunch of cards
type Deck struct {
	ID          uint	`json:"id"`
	Name        string  `json:"name"`
	TypeOfCards string  `json:"typeOfCards"`
}

var SampleDecks = []Deck{}

func init() {
	SampleDecks = append(SampleDecks, Deck{
		ID:          1,
		Name:        "Test",
		TypeOfCards: "words",
	})

	SampleDecks = append(SampleDecks, Deck{
		ID:          2,
		Name:        "Yes, English can be weird. It can be understood",
		TypeOfCards: "words",
	})
}