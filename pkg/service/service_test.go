package service

import (
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"os"
	"testing"
)

const (
	dataUsers       = "fixtures/users.json"
	dataAssets      = "fixtures/assets.json"
	dataComments    = "fixtures/comments.json"
	dataActions     = "fixtures/actions.json"
	dataIndexes     = "fixtures/indexes.json"
	dataMetadata    = "fixtures/metadata.json"
	dataTags        = "fixtures/tags.json"
	dataNewTags     = "fixtures/tags_rename.json"
	dataUserActions = "fixtures/user-actions.json"
)

func init() {
	db := db.NewMongoDB()
	defer db.Close()

	//Empty all test data
	db.TagTargets.RemoveAll(nil)
	db.Tags.RemoveAll(nil)
	db.Actions.RemoveAll(nil)
	db.Comments.RemoveAll(nil)
	db.Users.RemoveAll(nil)
	db.Assets.RemoveAll(nil)
	db.CayUserActions.RemoveAll(nil)
}

func TestCreateTag(t *testing.T) {
	file, err := os.Open(dataTags)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []model.Tag{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading tags ", err.Error())
	}

	c := NewContext()
	defer c.Close()
	for _, one := range objects {
		c.Input = one
		_, err := CreateUpdateTag(c)
		if err != nil {
			t.Fail()
		}
	}
}

func TestCreateAsset(t *testing.T) {
	file, err := os.Open(dataAssets)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []model.Asset{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading asset data", err.Error())
	}

	c := NewContext()
	defer c.Close()
	for _, one := range objects {
		c.Input = one
		_, err := CreateAsset(c)
		if err != nil {
			t.Fail()
		}
	}

	//try the same set again and it shouldn't fail
	for _, one := range objects {
		c.Input = one
		_, err := CreateAsset(c)
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

	objects := []model.User{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading user data", err.Error())
	}

	c := NewContext()
	defer c.Close()
	for _, one := range objects {
		c.Input = one
		_, err := CreateUser(c)
		if err != nil {
			t.Fail()
		}
	}
	//try the same set again and it shouldn't fail
	for _, one := range objects {
		c.Input = one
		_, err := CreateUser(c)
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

	objects := []model.Comment{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading user data", err.Error())
	}

	c := NewContext()
	defer c.Close()
	for _, one := range objects {
		c.Input = one
		_, err := CreateComment(c)
		if err != nil {
			t.Fail()
		}
	}
	//try the same set again and it shouldn't fail
	for _, one := range objects {
		c.Input = one
		_, err := CreateComment(c)
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

	objects := []model.Action{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading user data", err.Error())
	}

	c := NewContext()
	defer c.Close()
	for _, one := range objects {
		c.Input = one
		_, err := CreateAction(c)
		if err != nil {
			t.Fail()
		}
	}

	//Try again with the same data and it should all fail
	for _, one := range objects {
		c.Input = one
		_, err := CreateAction(c)
		if err == nil {
			t.Fail()
		}
	}
}

func TestCreateIndexes(t *testing.T) {
	file, err := os.Open(dataIndexes)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []model.Index{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading index data", err.Error())
	}

	c := NewContext()
	defer c.Close()
	for _, one := range objects {
		c.Input = one
		err := CreateIndex(c)
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

	objects := []model.Metadata{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading metadata ", err.Error())
	}

	c := NewContext()
	defer c.Close()
	for _, one := range objects {
		c.Input = one
		_, err := UpdateMetadata(c)
		if err != nil {
			t.Fail()
		}
	}
}

func TestUserActions(t *testing.T) {
	file, err := os.Open(dataUserActions)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []model.CayUserAction{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading user-actions ", err.Error())
	}

	c := NewContext()
	defer c.Close()
	for _, one := range objects {
		c.Input = one
		err := CreateUserAction(c)
		if err != nil {
			fmt.Printf("%+v", err)
			t.Fail()
		}
	}
}

func TestRenameTags(t *testing.T) {
	file, err := os.Open(dataNewTags)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []model.Tag{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading tags ", err.Error())
	}

	c := NewContext()
	defer c.Close()
	for _, one := range objects {
		c.Input = one
		_, err := CreateUpdateTag(c)
		if err != nil {
			t.Fail()
		}
	}
}

//func TestDeleteAllTag(t *testing.T) {
//	tags, err := GetTags()
//	if err != nil || len(tags) == 0 {
//		t.Fail()
//	}
//
//	for _, one := range tags {
//		err := DeleteTag(&one)
//		if err != nil {
//			t.Fail()
//		}
//	}
//
//	objects, err := GetTags()
//	if err != nil || len(objects) != 0 {
//		t.Fail()
//	}
//}
