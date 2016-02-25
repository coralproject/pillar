package mongodb

import (
	"crypto/tls"
	"net"

	"gopkg.in/mgo.v2"

	"github.com/coralproject/pillar/pkg/backend/iterator"
	"github.com/coralproject/pillar/pkg/model"
)

const (
// commentsCollection string = "comments"
)

var (
	indexMap = map[string][]mgo.Index{
	// mgo.Index{
	//   Key:      []string{"id"},
	//   Unique:   true,
	//   DropDups: true,
	// },
	}
)

// MongoDBBackend represents a MongoDB backend.
type MongoDBBackend struct {
	database string
	session  *mgo.Session
}

// NewMongoDBBackend creates a new MongoDBBackendBackend object using a
// MongoDB conection string. Any database defined in the connection string
// will be overridden by the argument-specified database.
func NewMongoDBBackend(addrs []string, database, username, password string, ssl bool) (*MongoDBBackend, error) {

	// Build a DialInfo object using the provided arguments.
	dialInfo := &mgo.DialInfo{
		Addrs:    addrs,
		Database: database,
		Username: username,
		Password: password,
	}

	// Determine whether or not to use TLS.
	if ssl {
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		}
	}

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil, err
	}

	m := &MongoDBBackend{
		database: database,
		session:  session,
	}

	// Ensure indicies are built.
	for collection, indicies := range indexMap {
		c := m.db().C(collection)
		for _, index := range indicies {
			if err := c.EnsureIndex(index); err != nil {
				return nil, err
			}
		}
	}

	return m, nil
}

func (m *MongoDBBackend) db() *mgo.Database {
	return m.session.DB(m.database)
}

func (m *MongoDBBackend) c(name string) *mgo.Collection {
	return m.db().C(name)
}

type iter struct {
	done, closed bool
	iter         *mgo.Iter
	result       func() interface{}
}

func (i *iter) Next() (interface{}, bool, error) {

	if !i.done {

		// Get a new instance.
		r := i.result()
		i.done = !(i.iter.Next(r))
		return r, i.done, nil
	}

	if !i.closed {
		i.closed = true
		if err := i.iter.Close(); err != nil {
			return nil, true, err
		}
	}

	return nil, true, nil
}

func (m *MongoDBBackend) Find(objectType string, query map[string]interface{}) (iterator.Iterator, error) {
	if err := model.ValidateObjectType(objectType); err != nil {
		return nil, err
	}

	return &iter{
		iter: m.c(objectType).Find(query).Iter(),
		result: func() interface{} {
			return model.ObjectTypeInstance(objectType)
		},
	}, nil
}

// Close closes the underlying session.
func (m *MongoDBBackend) Close() error {
	m.session.Close()
	return nil
}
