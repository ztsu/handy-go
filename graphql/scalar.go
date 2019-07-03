package graphql

import (
	"io"
	"strconv"
)

type Test int

func (t Test) MarshalGQL(w io.Writer) {
	w.Write([]byte(strconv.Itoa(int(t))))
}

func (t *Test) UnmarshalGQL(v interface{}) error {
	*t = 123
	return nil
}