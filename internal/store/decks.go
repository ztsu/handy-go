package store

type Deck struct {
	UUID        UUID   `json:"uuid"`
	UserID      UUID   `json:"userId"`
	Name        string `json:"name"`
	TypeOfCards string `json:"typeOfCards"`
}

type userDecks struct {
	UserID UUID   `json:"userId"`
	Decks  []UUID `json:"decks"`
}

type DeckStore interface {
	Get(UUID) (Deck, error)
	Save(Deck) error
	Delete(Deck) error
}
