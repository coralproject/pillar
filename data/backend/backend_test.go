package backend

import (
	"testing"

	"github.com/coralproject/pillar/data"
	"github.com/coralproject/pillar/data/backend/mongodb"
)

// testBackend is a method for testing Backend interfaces.
func testBackend(b Backend, t *testing.T) {

	testComment := &data.Comment{
		ID: "1234",
	}

	if err := b.SetComment(testComment); err != nil {
		t.Error(err)
	}

	retrievedComment, err := b.Comment(testComment.ID)
	if err != nil {
		t.Error(err)
	}

	if retrievedComment.ID != testComment.ID {
		t.Errorf("Retrieved comment ID, %s, did not match the test comment ID, %s", retrievedComment.ID, testComment.ID)
	}

	if err := b.DeleteComment(testComment.ID); err != nil {
		t.Error(err)
	}

	if err := b.Close(); err != nil {
		t.Error(err)
	}

}

func TestMongoDBBackend(t *testing.T) {
	b, err := mongodb.NewMongoDBBackend("mongodb://localhost:27017", "pillar")
	if err != nil {
		t.Fatal(err)
	}

	// Run the generic backend test.
	testBackend(b, t)
}
