package jsonutil

import (
	"encoding/json"
)

// Marshal returns the JSON encoding of v.
func Marshal(v interface{}) []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	return bytes
}
