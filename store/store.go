package store

type Card struct {
	UUID     UUID `json:"uuid"`
	DeckUUID UUID          `json:"deckUuid"`
	Type     string        `json:"type"`
	Viewed   uint64        `json:"viewed"`
	Opened   uint64        `json:"opened"`
}

type CardStore interface {
	Get(UUID) (Card, error)
	Save(Card) error
}

type Deck struct {
	UUID        UUID   `json:"uuid"`
	UserID      UUID   `json:"userId"`
	Name        string `json:"name"`
	TypeOfCards string `json:"typeOfCards"`
}

type DeckStore interface {
	Get(UUID) (Deck, error)
	Save(Deck) error
	Delete(Deck) error
}

type Translation struct {
	UUID        UUID   `json:"uuid"`
	From        string `json:"from"`
	To          string `json:"to"`
	Word        string `json:"word"`
	Translation string `json:"translation"`
	IPA         string `json:"ipa"`
}

type TranslationStore interface {
	Get(UUID) (Translation, error)
	Save(Translation) error
	Delete(UUID) error
}

