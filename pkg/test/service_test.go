package service

import (
	"encoding/json"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/service"
	"github.com/coralproject/pillar/pkg/web"
	"log"
	"os"
	"testing"
)

const (
	dataUsers    = "fixtures/import/users.json"
	dataAssets   = "fixtures/import/assets.json"
	dataComments = "fixtures/import/comments.json"
	dataActions  = "fixtures/import/actions.json"

	dataTags        = "fixtures/crud/tags.json"
	dataSections    = "fixtures/crud/sections.json"
	dataIndexes     = "fixtures/crud/indexes.json"
	dataMetadata    = "fixtures/crud/metadata.json"
	dataNewTags     = "fixtures/crud/tags_rename.json"
	dataUserActions = "fixtures/crud/user-actions.json"
	dataSearches    = "fixtures/crud/searches.json"
)

func init() {
	db := db.NewMongoDB()
	defer db.Close()

	//Empty all test data
	db.Tags.RemoveAll(nil)
	db.Sections.RemoveAll(nil)
	db.Authors.RemoveAll(nil)
	db.Actions.RemoveAll(nil)
	db.Comments.RemoveAll(nil)
	db.Users.RemoveAll(nil)
	db.Assets.RemoveAll(nil)
	db.CayUserActions.RemoveAll(nil)
	db.TagTargets.RemoveAll(nil)
	db.Searches.RemoveAll(nil)
}

func TestCreateSections(t *testing.T) {
	file, err := os.Open(dataSections)
	if err != nil {
		log.Fatalf("opening config file", err.Error())
	}

	objects := []model.Section{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading tags ", err.Error())
	}

	c := web.NewContext(nil, nil)
	defer c.Close()

	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.CreateUpdateSection(c); err != nil {
			t.Fail()
		}
	}
}

func TestCreateTags(t *testing.T) {
	file, err := os.Open(dataTags)
	if err != nil {
		log.Fatalf("opening config file", err.Error())
	}

	objects := []model.Tag{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading tags ", err.Error())
	}

	c := web.NewContext(nil, nil)
	defer c.Close()

	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.CreateUpdateTag(c); err != nil {
			t.Fail()
		}
	}
}

func TestCreateSearches(t *testing.T) {
	file, err := os.Open(dataSearches)
	if err != nil {
		log.Fatalf("opening config file", err.Error())
	}

	objects := []model.Search{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading userGroups ", err.Error())
	}

	c := web.NewContext(nil, nil)
	defer c.Close()

	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.CreateUpdateSearch(c); err != nil {
			t.Fail()
		}
	}

	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.CreateUpdateSearch(c); err != nil {
			t.Fail()
		}
	}
}

func TestImportAssets(t *testing.T) {
	file, err := os.Open(dataAssets)
	if err != nil {
		log.Fatalf("opening config file", err.Error())
	}

	objects := []model.Asset{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading assets data", err.Error())
	}

	c := web.NewContext(nil, nil)
	defer c.Close()

	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.ImportAsset(c); err != nil {
			t.Fail()
		}
	}

	//try the same set again and it shouldn't fail
	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.ImportAsset(c); err != nil {
			t.Fail()
		}
	}
}

func TestImportUsers(t *testing.T) {
	file, err := os.Open(dataUsers)
	if err != nil {
		log.Fatalf("opening config file", err.Error())
	}

	objects := []model.User{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading users data", err.Error())
	}

	c := web.NewContext(nil, nil)
	defer c.Close()

	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.ImportUser(c); err != nil {
			t.Fail()
		}
	}
	//try the same set again and it shouldn't fail
	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.ImportUser(c); err != nil {
			t.Fail()
		}
	}
}

func TestImportComments(t *testing.T) {
	file, err := os.Open(dataComments)
	if err != nil {
		log.Fatalf("opening config file", err.Error())
	}

	objects := []model.Comment{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading comments data", err.Error())
	}

	c := web.NewContext(nil, nil)
	defer c.Close()

	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.ImportComment(c); err != nil {
			log.Fatalf("Error: %s\n", err)
			t.Fail()
		}
	}
	//try the same set again and it shouldn't fail
	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.ImportComment(c); err != nil {
			log.Fatalf("Error: %s\n", err)
			t.Fail()
		}
	}
}

func TestImportActions(t *testing.T) {
	file, err := os.Open(dataActions)
	if err != nil {
		log.Fatalf("opening config file", err.Error())
	}

	objects := []model.Action{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading user data", err.Error())
	}

	c := web.NewContext(nil, nil)
	defer c.Close()

	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.ImportAction(c); err != nil {
			log.Fatalf("Error: %s\n", err)
			t.Fail()
		}
	}

	//Try again with the same data and it should all fail
	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.ImportAction(c); err == nil {
			log.Fatalf("Error: %s\n", err)
			t.Fail()
		}
	}
}

func TestCreateIndexes(t *testing.T) {
	file, err := os.Open(dataIndexes)
	if err != nil {
		log.Fatalf("opening config file", err.Error())
	}

	objects := []model.Index{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading index data", err.Error())
	}

	c := web.NewContext(nil, nil)
	defer c.Close()

	for _, one := range objects {
		c.Marshall(one)
		if err := service.CreateIndex(c); err != nil {
			log.Fatalf("Error: %s\n", err)
			t.Fail()
		}
	}
}

func TestUpdateMetadata(t *testing.T) {
	file, err := os.Open(dataMetadata)
	if err != nil {
		log.Fatalf("opening config file", err.Error())
	}

	objects := []model.Metadata{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading metadata ", err.Error())
	}

	c := web.NewContext(nil, nil)
	defer c.Close()

	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.UpdateMetadata(c); err != nil {
			log.Fatalf("Error: %s\n", err)
			t.Fail()
		}
	}
}

func TestUserActions(t *testing.T) {
	file, err := os.Open(dataUserActions)
	if err != nil {
		log.Fatalf("opening config file", err.Error())
	}

	objects := []model.CayUserAction{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading user-actions ", err.Error())
	}

	c := web.NewContext(nil, nil)
	defer c.Close()

	for _, one := range objects {
		c.Marshall(one)
		if err := service.CreateUserAction(c); err != nil {
			log.Fatalf("Error: %s\n", err)
			t.Fail()
		}
	}
}

func TestRenameTags(t *testing.T) {
	file, err := os.Open(dataNewTags)
	if err != nil {
		log.Fatalf("opening config file", err.Error())
	}

	objects := []model.Tag{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading tags ", err.Error())
	}

	c := web.NewContext(nil, nil)
	defer c.Close()

	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.CreateUpdateTag(c); err != nil {
			log.Fatalf("Error: %s\n", err)
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
