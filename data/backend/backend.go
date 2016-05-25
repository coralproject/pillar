package backend

import (
	"github.com/coralproject/pillar/data"
)

// Backend defines methods necessary for data storage and retrieval.
//
// Notably, methods that perform function such as index creation have been
// ommitted; the thinking being that index creation should be internal to the
// Backend implementation (if it's necessary at all).
type Backend interface {
	Comment(string) (*data.Comment, error)
	SetComment(*data.Comment) error
	DeleteComment(string) error

	// Close should terminate any open connections and perform any necessary
	// cleanup.
	Close() error
}
