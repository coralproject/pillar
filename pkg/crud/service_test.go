package crud

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

const (
	dataUsers    = "fixtures/users.json"
	dataAssets   = "fixtures/assets.json"
	dataComments = "fixtures/comments.json"
	dataActions  = "fixtures/actions.json"
	dataIndexes  = "fixtures/indexes.json"
	dataMetadata = "fixtures/metadata.json"
	dataTags     = "fixtures/tags.json"
)

func init() {
	mm := GetMongoManager()
	mm.TagTargets.RemoveAll(nil)
	mm.Tags.RemoveAll(nil)
	mm.Actions.RemoveAll(nil)
	mm.Comments.RemoveAll(nil)
	mm.Users.RemoveAll(nil)
	mm.Assets.RemoveAll(nil)
}

func TestCreateAsset(t *testing.T) {
	file, err := os.Open(dataAssets)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []Asset{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading asset data", err.Error())
	}

	for _, one := range objects {
		_, err := CreateAsset(&one)
		if err != nil {
			t.Fail()
		}
	}
}

func TestCreateUser(t *testing.T) {
	file, err := os.Open(dataUsers)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []User{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading user data", err.Error())
	}

	for _, one := range objects {
		_, err := CreateUser(&one)
		if err != nil {
			t.Fail()
		}
	}
}

func TestCreateComments(t *testing.T) {
	file, err := os.Open(dataComments)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []Comment{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading user data", err.Error())
	}

	for _, one := range objects {
		_, err := CreateComment(&one)
		if err != nil {
			t.Fail()
		}
	}
}

func TestCreateActions(t *testing.T) {
	file, err := os.Open(dataActions)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []Action{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading user data", err.Error())
	}

	for _, one := range objects {
		_, err := CreateAction(&one)
		if err != nil {
			t.Fail()
		}
	}
}

func TestCreateIndexes(t *testing.T) {
	file, err := os.Open(dataIndexes)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []Index{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading index data", err.Error())
	}

	for _, one := range objects {
		err := CreateIndex(&one)
		if err != nil {
			t.Fail()
		}
	}
}

func TestUpdateMetadata(t *testing.T) {
	file, err := os.Open(dataMetadata)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []Metadata{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading metadata ", err.Error())
	}

	for _, one := range objects {
		_, err := UpdateMetadata(&one)
		if err != nil {
			t.Fail()
		}
	}
}

func TestUpsertTag(t *testing.T) {
	file, err := os.Open(dataTags)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []Tag{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading tags ", err.Error())
	}

	for _, one := range objects {
		t.Logf("Tag: %+v\n\n", one)
		_, err := UpsertTag(&one)
		if err != nil {
			t.Fail()
		}
	}
}
