package handy

type UserDecks struct {
	UserID UUID   `json:"userId"`
	Decks  []UUID `json:"decks"`
}

type UserDecksStore interface {
	Get(userID UUID)
	Save(UserDecks)
	Delete(userID UUID)
}