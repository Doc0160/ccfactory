package storage

import (
	"ccfactory/server/factory"
)

type Factory = factory.Factory
type BusAccess = factory.BusAccess
type Storage = factory.Storage

func Into[T any](r factory.RawMessage) (T, error) {
	return factory.Into[T](r)
}
