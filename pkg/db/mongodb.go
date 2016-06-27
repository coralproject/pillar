package db

import (
	"log"
	"sync"

	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2"
)

var (
	mgoSession *mgo.Session
	mu         sync.Mutex
)

// MongoDB encapsulates a mongo database and session
type MongoDB struct {
	Session *mgo.Session
	DB      *mgo.Database
}

//Close closes the mongodb session; must be called, else the session remains open
func (m *MongoDB) Close() {
	if m.Session == nil {
		return
	}
	m.Session.Close()
}

//IsValid returns true for a valid session, false otherwise
func (m *MongoDB) IsValid() bool {
	return m.Session != nil
}

//Upsert upserts a specific entity into the given collection
func (m *MongoDB) Upsert(objectType string, id, object interface{}) error {
	session := m.Session.Copy()
	defer session.Close()

	_, err := m.Session.DB("").C(objectType).UpsertId(id, object)
	return err
}

//Find finds from a collection using the query
func (m *MongoDB) Find(objectType string, query map[string]interface{}) (Iterator, error) {
	if err := model.ValidateObjectType(objectType); err != nil {
		return nil, err
	}

	session := m.Session.Copy()
	return &iter{
		session: session,
		iter:    session.DB("").C(objectType).Find(query).Iter(),
		result: func() interface{} {
			return model.ObjectTypeInstance(objectType)
		},
	}, nil
}

func connect(url string) *mgo.Session {
	mu.Lock()
	defer mu.Unlock()

	if mgoSession != nil {
		return mgoSession
	}

	session, err := mgo.Dial(url)
	if err != nil {
		log.Printf("Error connecting to Mongo Database [%v]", err)
		return nil
	}

	// Ensure indicies are built.
	for _, one := range model.Indicies {
		c := session.DB("").C(one.Target)
		if err := c.EnsureIndex(one.Index); err != nil {
			log.Fatalf("Error building indicies: %s", err)
		}
	}

	//save the main session for reuse
	mgoSession = session
	return mgoSession
}

func NewMongoDB(url string) *MongoDB {
	db := MongoDB{}

	s := connect(url)
	if s != nil {
		db.Session = s.Copy()
		db.DB = s.DB("")
	}

	return &db
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
