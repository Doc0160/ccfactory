package factory

import (
	"encoding/json"
)

type RawMessage = json.RawMessage

func Into[T any](r RawMessage) (T, error) {
	var v T
	err := json.Unmarshal(r, &v)
	return v, err
}
