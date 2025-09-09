package storage

import (
	"ccfactory/server/factory"
)

func Into[T any](r factory.RawMessage) (T, error) {
	return factory.Into[T](r)
}
