package handy

// Card is a card
type Card struct {
	ID uint `json:"id"`
	DeckID uint `json:"deckId"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
	IPA         string `json:"ipa"`
}

var SampleCards = []Card{}

func init() {
	SampleCards = append(SampleCards, Card{1, 1, "handy", "удобный", "ˈhændɪ"})
	SampleCards = append(SampleCards, Card{1, 2, "through", "через", "θruː"})
	SampleCards = append(SampleCards, Card{1, 2, "tough", "жесткий", "tʌf"})
	SampleCards = append(SampleCards, Card{1, 2, "thorough", "полный", "ˈθʌrə"})
	SampleCards = append(SampleCards, Card{1, 2, "thought", "мысль", "θɔːt"})
	SampleCards = append(SampleCards, Card{1, 2, "though", "хотя", "ðəʊ"})
}