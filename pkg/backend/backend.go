package backend

import (
	"errors"

	"github.com/coralproject/pillar/pkg/backend/iterator"
)

type Query map[interface{}]interface{}

var (
	BackendNotInitializedError = errors.New("backend not initialized")
	BackendTypeError           = errors.New("object type is incorrect")
)

type Backend interface {
	Find(objectType string, query map[string]interface{}) (iterator.Iterator, error)
	FindID(objectType string, id interface{}) (interface{}, error)
	Upsert(objectType string, query map[string]interface{}, object interface{}) error
	UpsertID(objectType string, id, object interface{}) error
	Close() error
}
