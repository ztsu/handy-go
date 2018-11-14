package handy

type Card struct {
	UUID     UUID   `json:"uuid"`
	DeckUUID UUID   `json:"deckUuid"`
	Type     string `json:"type"`
	Viewed   uint64 `json:"viewed"`
	Opened   uint64 `json:"opened"`
}

type CardStore interface {
	Get(UUID) (Card, error)
	Save(Card) error
}