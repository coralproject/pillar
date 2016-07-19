package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
)

const (
	dataUsers    = "fixtures/import/users.json"
	dataAssets   = "fixtures/import/assets.json"
	dataComments = "fixtures/import/comments.json"
	dataActions  = "fixtures/import/actions.json"

	dataTags        = "fixtures/crud/tags.json"
	dataSections    = "fixtures/crud/sections.json"
	dataIndexes     = "fixtures/crud/indexes.json"
	dataFormIndexes = "fixtures/crud/form_indexes.json"
	dataMetadata    = "fixtures/crud/metadata.json"
	dataNewTags     = "fixtures/crud/tags_rename.json"
	dataUserActions = "fixtures/crud/user-actions.json"
	dataSearches    = "fixtures/crud/searches.json"

	dataForms           = "fixtures/crud/forms.json"
	dataFormSubmissions = "fixtures/crud/form_submissions.json"

	dataFormsIds           = "fixtures/crud/forms_with_id.json"
	dataFormSubmissionsIds = "fixtures/crud/form_submissions_with_id.json"

	dataFormGalleriesIds = "fixtures/crud/form_galleries_with_id.json"
	dataFormGalleries    = "fixtures/crud/form_galleries.json"
)

var (
	envmongodburl string
	testDBname    string
)

// remove the test database
func emptyDB() {
	deb := db.NewMongoDB(testDBname)
	defer deb.Close()

	if e := deb.DB.DropDatabase(); e != nil {
		log.Fatalf("Fail removing DropDatabase %s. Error: %v", deb.DB.Name, e)
	}
}

func recoverEnvVariables() {
	e := os.Setenv("MONGODB_URL", envmongodburl)
	if e != nil {
		log.Fatal("Error when setting back the MONGODB_URL environment variable.", e)
	}
}

// load forms fixtures
func loadformfixtures() {

	// connect to test mongo database
	deb := db.NewMongoDB(os.Getenv("MONGODB_URL"))
	defer deb.Close()

	// get the fixtures from appropiate json file for data form submission
	file, e := ioutil.ReadFile(dataFormsIds)
	if e != nil {
		log.Fatalf("opening config file %v", e.Error())
	}

	var objects model.Form
	e = json.Unmarshal(file, &objects)
	if e != nil {
		log.Fatalf("Error reading forms. %v", e.Error())
	}

	// insert in bulk into mongo database
	b := deb.DB.C(model.Forms).Bulk()
	b.Unordered()
	b.Insert(objects)

	if _, e = b.Run(); e != nil {
		log.Fatalf("Error when loading fixtures for %s. Error: %v", model.Forms, e)
	}

	// get the fixtures from appropiate json file for data form submission
	file, e = ioutil.ReadFile(dataFormSubmissionsIds)
	if e != nil {
		log.Fatalf("opening config file %v", e.Error())
	}

	var objectsS []model.FormSubmission
	e = json.Unmarshal(file, &objectsS)
	if e != nil {
		log.Fatalf("Error reading forms submisions. %v", e.Error())
	}

	// insert in bulk into mongo database
	b = deb.DB.C(model.FormSubmissions).Bulk()
	b.Unordered()
	for _, o := range objectsS {
		b.Insert(o)
	}

	if _, e = b.Run(); e != nil {
		log.Fatalf("Error when loading fixtures for %s. Error: %v", model.FormSubmissions, e)
	}
}

// load form galleries fixtures
func loadformgalleriesfixtures() {

	// connect to test mongo database
	deb := db.NewMongoDB(os.Getenv("MONGODB_URL"))
	defer deb.Close()

	// get the fixtures from appropiate json file for data form submission
	file, e := ioutil.ReadFile(dataFormGalleriesIds)
	if e != nil {
		log.Fatalf("opening config file %v", e.Error())
	}

	var objects model.FormGallery
	e = json.Unmarshal(file, &objects)
	if e != nil {
		log.Fatalf("Error reading forms galleries. %v", e.Error())
	}

	// insert in bulk into mongo database
	b := deb.DB.C(model.FormGalleries).Bulk()
	b.Unordered()
	b.Insert(objects)

	if _, e = b.Run(); e != nil {
		log.Fatalf("Error when loading fixtures for %s. Error: %v", model.FormGalleries, e)
	}
}

// search for a string in a struct
func find(s string, r []model.FormSubmission) bool {
	found := false

	for _, i := range r {
		for _, f := range i.Answers {
			switch k := f.Answer.(type) {
			case string:
				found = strings.Contains(k, s)
				if found {
					return found
				}
			}
		}

		switch k := i.Footer.(type) {
		case string:
			found = strings.Contains(k, s)
			if found {
				return found
			}
		}

		fmt.Println(i.Header, reflect.TypeOf(i.Header))

		switch k := i.Header.(type) {
		case bson.M:
			fmt.Println("string ", k)
			for _, p := range k {
				found = strings.Contains(p.(string), s)
				if found {
					return found
				}
			}
		}

		for _, k := range i.Flags {
			found = strings.Contains(k, s)
			if found {
				return found
			}
		}
		switch k := i.Header.(type) {
		case string:
			found = strings.Contains(k, s)
			if found {
				return found
			}
		}
		found = strings.Contains(i.Status, s)
		if found {
			return found
		}
	}

	return found
}

// set the environment variables for a test database
// in the tests we are creating a unique database for testing
// and removing it after running all the tests
func setTestDatabase() {

	// save back the MONGODB_URL env variable that it may be used for production
	// this is sketchy and should be consider in the refactoring. We have the db package and web package getting MONGODB_URL everywhere
	envmongodburl = os.Getenv("MONGODB_URL")
	if envmongodburl == "" {
		log.Fatal("Setup environmental variable MONGODB_URL with connection string.")
	}
	// the name of the database the instance use plus test plus a timestamp
	testDBname = envmongodburl + strconv.FormatInt(time.Now().UTC().Unix(), 10)
	e := os.Setenv("MONGODB_URL", testDBname)
	if e != nil {
		log.Fatal("Error when setting environment test ", e)
	}
}

func getDataFormSubmissions(fileName string) []model.FormSubmission {

	var objects []model.FormSubmission

	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Fatalf("opening config file %v", e.Error())
	}

	e = json.Unmarshal(file, &objects)
	if e != nil {
		log.Fatalf("Error reading forms submissions. %v", e.Error())
	}

	return objects
}

func getDataTags(fileName string) []model.Tag {
	var objects []model.Tag

	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Fatalf("opening config file %v", e.Error())
	}

	e = json.Unmarshal(file, &objects)
	if e != nil {
		log.Fatalf("Error reading tags. %v", e.Error())
	}

	return objects
}

func getMetadata(fileName string) []model.Metadata {
	var objects []model.Metadata
	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Fatalf("opening config file %v", e.Error())
	}

	e = json.Unmarshal(file, &objects)
	if e != nil {
		log.Fatalf("Error reading metadata. %v", e.Error())
	}
	return objects
}

func getDataAssets(fileName string) []model.Asset {
	var assets []model.Asset

	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Fatalf("When opening fixtures on assets: %v", e.Error())
	}

	e = json.Unmarshal(file, &assets)
	if e != nil {
		log.Fatalf("Error reading assets: %v", e.Error())
	}

	return assets
}

func getDataUsers(fileName string) []model.User {
	var users []model.User

	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		log.Fatalf("opening config file %v", e.Error())
	}

	e = json.Unmarshal(file, &users)
	if e != nil {
		log.Fatalf("Error reading assets. %v", e.Error())
	}

	return users
}
