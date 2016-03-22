package backend

import (
	"errors"
	"sync"

	"golang.org/x/net/context"

	"github.com/coralproject/pillar/pkg/backend/iterator"
)

const (
	backendContextKey = "backend"
)

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

func NewBackendContext(ctx context.Context, b Backend) context.Context {
	return context.WithValue(ctx, backendContextKey, b)
}

func BackendFromContext(ctx context.Context) (Backend, error) {
	b, ok := ctx.Value(backendContextKey).(Backend)
	if !ok {
		return nil, BackendNotInitializedError
	}
	return b, nil
}

type IdentityMap struct {
	Backend
	objects map[string]map[interface{}]interface{}
	mu      sync.Mutex
}

func NewIdentityMap(b Backend) *IdentityMap {
	return &IdentityMap{
		Backend: b,
		objects: make(map[string]map[interface{}]interface{}),
	}
}

func (m *IdentityMap) FindID(objectType string, id interface{}) (interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.objects[objectType]; !ok {
		m.objects[objectType] = make(map[interface{}]interface{})
	}

	if value, ok := m.objects[objectType][id]; ok {
		return value, nil
	}

	value, err := m.Backend.FindID(objectType, id)
	if err != nil {
		return nil, err
	}

	m.objects[objectType][id] = value
	return value, nil
}
