package backend

import (
	"crypto/tls"
	"net"

	"gopkg.in/mgo.v2"

	"github.com/coralproject/pillar/pkg/model"
)

const (
// commentsCollection string = "comments"
)

var (
	indexMap = map[string][]mgo.Index{
		"comments": []mgo.Index{
			mgo.Index{
				Key: []string{"user_id"},
			},
		},
		"actions": []mgo.Index{
			mgo.Index{
				Key: []string{"user_id"},
			},
			mgo.Index{
				Key: []string{"target", "target_id"},
			},
		},
		"dimensions": []mgo.Index{
			mgo.Index{
				Key:      []string{"name"},
				Unique:   true,
				DropDups: true,
			},
		},
	}
)

// MongoDBBackend represents a MongoDB backend.
type MongoDBBackend struct {
	database string
	session  *mgo.Session
	sem      chan struct{}
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
		c := session.DB(m.database).C(collection)
		for _, index := range indicies {
			if err := c.EnsureIndex(index); err != nil {
				return nil, err
			}
		}
	}

	return m, nil
}

// type session struct {
// 	session *mgo.Session
// 	sem
// }

func (m *MongoDBBackend) newSession() *mgo.Session {
	return m.session.Copy()
}

func (m *MongoDBBackend) Upsert(objectType string, query map[string]interface{}, object interface{}) error {
	session := m.newSession()
	defer session.Close()

	_, err := session.DB(m.database).C(objectType).Upsert(query, object)
	return err
}

func (m *MongoDBBackend) UpsertID(objectType string, id, object interface{}) error {
	session := m.newSession()
	defer session.Close()

	_, err := session.DB(m.database).C(objectType).UpsertId(id, object)
	return err
}

type iter struct {
	done    bool
	iter    *mgo.Iter
	result  func() interface{}
	session *mgo.Session
}

func (i *iter) Next() (interface{}, bool, error) {

	if !i.done {

		// Get a new instance.
		r := i.result()
		i.done = !(i.iter.Next(r))

		// Check for a timeout. If one occured, re-run Next (to theoretically
		// reconnect).
		if i.done && i.iter.Timeout() {
			i.done = !(i.iter.Next(r))
		}

		// If ther iterator is done, we can close it and the underlying session.
		var err error
		if i.done {
			err = i.iter.Close()
			i.session.Close()
		}

		return r, i.done, err
	}

	return nil, true, nil
}

func (m *MongoDBBackend) Find(objectType string, query map[string]interface{}) (Iterator, error) {
	if err := model.ValidateObjectType(objectType); err != nil {
		return nil, err
	}

	session := m.newSession()
	return &iter{
		session: session,
		iter:    session.DB(m.database).C(objectType).Find(query).Iter(),
		result: func() interface{} {
			return model.ObjectTypeInstance(objectType)
		},
	}, nil
}

func (m *MongoDBBackend) FindID(objectType string, id interface{}) (interface{}, error) {
	if err := model.ValidateObjectType(objectType); err != nil {
		return nil, err
	}

	session := m.newSession()
	defer session.Close()
	result := model.ObjectTypeInstance(objectType)
	err := session.DB(m.database).C(objectType).FindId(id).One(result)
	return result, err
}

// Close closes the underlying session.
func (m *MongoDBBackend) Close() error {
	m.session.Close()
	return nil
}
