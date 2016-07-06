package service_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

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
)

// check if we are working with the test db
func checkIsTestDB() bool {
	s := strings.Split(os.Getenv("MONGODB_URL"), "_")
	return (s[len(s)-1] == "test")
}

// remove the test database
func emptydb() {
	if test := checkIsTestDB(); !test {
		log.Fatalf("Fail to setup test database with %s", os.Getenv("MONGODB_URL"))
	}
	deb := db.NewMongoDB(os.Getenv("MONGODB_URL"))
	defer deb.Close()

	if e := deb.DB.DropDatabase(); e != nil {
		log.Fatalf("Fail removing DropDatabase %s. Error: %v", deb.DB.Name, e)
	}
}

// load forms fixtures
func loadformfixtures() {
	if test := checkIsTestDB(); !test {
		log.Fatalf("Fail to setup test database with %s", os.Getenv("MONGODB_URL"))
	}

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

	_, e = b.Run()
	if e != nil {
		log.Fatalf("Error when loading fixtures for %s and %s. Error: %v", model.Forms, model.FormSubmissions, e)
	}

	// get the fixtures from appropiate json file for data form submission
	file, e = ioutil.ReadFile(dataFormSubmissionsIds)
	if e != nil {
		log.Fatalf("opening config file %v", e.Error())
	}

	var objectsS model.FormSubmission
	e = json.Unmarshal(file, &objectsS)
	if e != nil {
		log.Fatalf("Error reading forms submisions. %v", e.Error())
	}

	// insert in bulk into mongo database
	b = deb.DB.C(model.FormSubmissions).Bulk()
	b.Insert(objectsS)

	_, e = b.Run()
	if e != nil {
		log.Fatalf("Error when loading fixtures for %s and %s. Error: %v", model.Forms, model.FormSubmissions, e)
	}

	// create text indexes
	file, e = ioutil.ReadFile(dataFormIndexes) // os.Open(dataIndexes)
	if e != nil {
		log.Fatalf("opening config file %v", e.Error())
	}

	objectsI := []model.Index{}
	e = json.Unmarshal(file, &objectsI)
	if e != nil {
		log.Fatalf("Error reading index data %v", e.Error())
	}

	for _, i := range objectsI {
		err := deb.DB.C(i.Target).EnsureIndex(i.Index)
		if err != nil {
			log.Fatalf("Error %v creating index [%+v]", err, i)
		}
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
