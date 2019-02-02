package http

import (
	"encoding/json"
	"fmt"
	"github.com/ztsu/handy-go"
	"net/http"
)

func NewCreateTranslationHandler(repository handy.TranslationStore) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var tr handy.Translation

		err := json.NewDecoder(r.Body).Decode(&tr)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		err = repository.Save(tr)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("Error: %s", err)))
			return
		}

		json.NewEncoder(w).Encode(tr)
	}
}