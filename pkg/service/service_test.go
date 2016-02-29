package service

import (
	"encoding/json"
	"fmt"
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
	mm := GetMongoManager()
	mm.TagTargets.RemoveAll(nil)
	mm.Tags.RemoveAll(nil)
	mm.Actions.RemoveAll(nil)
	mm.Comments.RemoveAll(nil)
	mm.Users.RemoveAll(nil)
	mm.Assets.RemoveAll(nil)
	mm.CayUserActions.RemoveAll(nil)
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

	for _, one := range objects {
		_, err := CreateUpdateTag(&one)
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

	for _, one := range objects {
		_, err := CreateAsset(&one)
		if err != nil {
			t.Fail()
		}
	}

	//try the same set again and it shouldn't fail
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

	objects := []model.User{}
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
	//try the same set again and it shouldn't fail
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

	objects := []model.Comment{}
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
	//try the same set again and it shouldn't fail
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

	objects := []model.Action{}
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

	//Try again with the same data and it should all fail
	for _, one := range objects {
		_, err := CreateAction(&one)
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

	objects := []model.Metadata{}
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

	for _, one := range objects {
		err := CreateUserAction(&one)
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

	for _, one := range objects {
		_, err := CreateUpdateTag(&one)
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
