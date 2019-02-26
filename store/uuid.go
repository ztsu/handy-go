package store

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
)

type UUID uuid.UUID

func (id *UUID) UnmarshalJSON(b []byte) error {
	parsed, err := uuid.Parse(string(b))

	*id = UUID(parsed)
	if err != nil {
		return errors.New("could not parse UUID")
	}

	return nil
}

func (id UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(uuid.UUID(id).String())
}

func (id UUID) MarshalBinary() []byte {
	return id[:]
}

func (id UUID) String() string {
	return uuid.UUID(id).String()
}
