package db

import (
	"gopkg.in/mgo.v2"

	"github.com/coralproject/pillar/pkg/backend/iterator"
	"github.com/coralproject/pillar/pkg/model"
	"log"
	"os"
)

const (
	defaultMongoURL string = "mongodb://localhost:27017/coral"
)

var (
	mgoSession *mgo.Session
)

// MongoDB encapsulates a mongo session with all relevant collections
type MongoDB struct {
	Session        *mgo.Session
	Assets         *mgo.Collection
	Users          *mgo.Collection
	Actions        *mgo.Collection
	Comments       *mgo.Collection
	Tags           *mgo.Collection
	Authors        *mgo.Collection
	Sections       *mgo.Collection
	TagTargets     *mgo.Collection
	CayUserActions *mgo.Collection
}

//Close closes the mongodb session; must be called, else the session remains open
func (m *MongoDB) Close() {
	m.Session.Close()
}

//Upsert upserts a specific entity into the given collection
func (m *MongoDB) Upsert(objectType string, id, object interface{}) error {
	session := m.Session.Clone()
	defer session.Close()

	_, err := m.Session.DB("").C(objectType).UpsertId(id, object)
	return err
}

//Find finds from a collection using the query
func (m *MongoDB) Find(objectType string, query map[string]interface{}) (iterator.Iterator, error) {
	if err := model.ValidateObjectType(objectType); err != nil {
		return nil, err
	}

	session := m.Session.Clone()
	return &iter{
		session: session,
		iter:    session.DB("").C(objectType).Find(query).Iter(),
		result: func() interface{} {
			return model.ObjectTypeInstance(objectType)
		},
	}, nil
}

func init() {
	url := os.Getenv("MONGODB_URL")
	if url == "" {
		log.Printf("$MONGODB_URL not found, trying to connect locally [%s]", defaultMongoURL)
		url = defaultMongoURL
	}

	session, err := mgo.Dial(url)
	if err != nil {
		log.Fatalf("Error connecting to mongo database: %s", err)
	}

	// Ensure indicies are built.
	for _, one := range model.Indicies {
		c := session.DB("").C(one.Target)
		if err := c.EnsureIndex(one.Index); err != nil {
			log.Fatalf("Error building indicies: %s", err)
		}
	}

	//keep it for reuse
	mgoSession = session
}

//NewMongoDB returns a cloned MongoDB Session
func NewMongoDB() *MongoDB {
	//	if mgoSession == nil {
	//		initDB()
	//	}
	//
	db := MongoDB{}
	db.Session = mgoSession.Clone() //must clone
	db.Users = db.Session.DB("").C(model.Users)
	db.Assets = db.Session.DB("").C(model.Assets)
	db.Actions = db.Session.DB("").C(model.Actions)
	db.Comments = db.Session.DB("").C(model.Comments)
	db.Authors = db.Session.DB("").C(model.Authors)
	db.Sections = db.Session.DB("").C(model.Sections)
	db.Tags = db.Session.DB("").C(model.Tags)
	db.TagTargets = db.Session.DB("").C(model.TagTargets)
	db.CayUserActions = db.Session.DB("").C(model.CayUserActions)

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
