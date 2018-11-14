package handy

type Deck struct {
	UUID        UUID   `json:"uuid"`
	UserID      UUID   `json:"userId"`
	Name        string `json:"name"`
	TypeOfCards string `json:"typeOfCards"`
}

type DeckStore interface {
	Get(UUID) (Deck, error)
	Save(Deck) error
	Delete(UUID) error
}
