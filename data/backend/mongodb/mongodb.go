package mongodb

import (
	"github.com/coralproject/pillar/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	commentsCollection string = "comments"
)

var (
	indicies = []mgo.Index{
		mgo.Index{
			Key:      []string{"id"},
			Unique:   true,
			DropDups: true,
		},
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
func NewMongoDBBackend(url, database string) (*MongoDBBackend, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}

	m := &MongoDBBackend{
		database: database,
		session:  session,
	}

	// Ensure indicies are built.
	for _, index := range indicies {
		if err := m.comments().EnsureIndex(index); err != nil {
			return nil, err
		}
	}

	return m, nil
}

// comments is shorthand for the comments collection object.
func (m *MongoDBBackend) comments() *mgo.Collection {
	return m.session.DB(m.database).C(commentsCollection)
}

// Comment returns the comment with the specified ID if it exists.
func (m *MongoDBBackend) Comment(id string) (*data.Comment, error) {
	comment := &data.Comment{}
	if err := m.comments().Find(bson.M{"id": id}).One(comment); err != nil {
		return nil, err
	}
	return comment, nil
}

// SetComment creates or updates a comment with the ID of the specified
// comment object.
func (m *MongoDBBackend) SetComment(comment *data.Comment) error {
	if _, err := m.comments().Upsert(bson.M{"id": comment.ID}, comment); err != nil {
		return err
	}
	return nil
}

// DeleteComment removes a comment with the specified ID.
func (m *MongoDBBackend) DeleteComment(id string) error {
	if err := m.comments().Remove(bson.M{"id": id}); err != nil {
		return err
	}
	return nil
}

// Close closes the underlying session.
func (m *MongoDBBackend) Close() error {
	m.session.Close()
	return nil
}
