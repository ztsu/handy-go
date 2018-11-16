package store

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