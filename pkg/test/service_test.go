package service

import (
	"encoding/json"
	"fmt"
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

	dataForms           = "fixtures/crud/forms.json"
	dataFormSubmissions = "fixtures/crud/form_submissions.json"
)

func init() {
	db := db.NewMongoDB(os.Getenv("MONGODB_URL"))
	defer db.Close()

	//Empty all test data
	db.DB.C(model.Forms).RemoveAll(nil)
	db.DB.C(model.FormSubmissions).RemoveAll(nil)
	db.DB.C(model.FormGalleries).RemoveAll(nil)
	db.DB.C(model.Tags).RemoveAll(nil)
	db.DB.C(model.Sections).RemoveAll(nil)
	db.DB.C(model.Authors).RemoveAll(nil)
	db.DB.C(model.Actions).RemoveAll(nil)
	db.DB.C(model.Comments).RemoveAll(nil)
	db.DB.C(model.Users).RemoveAll(nil)
	db.DB.C(model.Assets).RemoveAll(nil)
	db.DB.C(model.CayUserActions).RemoveAll(nil)
	db.DB.C(model.TagTargets).RemoveAll(nil)
	db.DB.C(model.Searches).RemoveAll(nil)
}

func TestCreateForms(t *testing.T) {
	file, err := os.Open(dataForms)
	if err != nil {
		log.Fatalf("opening config file", err.Error())
	}

	objects := []model.Form{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading forms ", err.Error())
	}

	c := web.NewContext(nil, nil)
	defer c.Close()

	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.CreateUpdateForm(c); err != nil {
			t.Fail()
		}
	}
}

func getAForm(t *testing.T) *model.Form {

	c := web.NewContext(nil, nil)
	defer c.Close()

	fs := []model.Form{}

	// let's see if we have forms to reply to
	fs, err := service.GetForms(c)
	if err != nil {
		log.Fatalf("Could not load forms for the test", err)
		t.Fail()
	}

	// Get the first form in the collection to reply to
	return &fs[0]

}

func getASubmissionToAForm(f *model.Form, t *testing.T) *model.FormSubmission {

	// create the context for this form
	c := web.NewContext(nil, nil)
	defer c.Close()
	c.SetValue("form_id", f.ID.Hex())

	s, err := service.GetFormSubmissionsByForm(c)
	if err != nil {
		log.Fatalf("Could not load forms submissions for the test", err)
		t.Fail()
	}

	return &s[0]
}

func getAGalleryFormAForm(f *model.Form, t *testing.T) *model.FormGallery {

	// create the context for this form
	c := web.NewContext(nil, nil)
	defer c.Close()
	c.SetValue("form_id", f.ID.Hex())

	g, err := service.GetFormGalleriesByForm(c)
	if err != nil {
		log.Fatalf("Could not load forms galleries for the test", err)
		t.Fail()
	}

	return &g[0]
}

func TestFormSubmissionsSchenanigans(t *testing.T) {

	file, err := os.Open(dataFormSubmissions)
	if err != nil {
		log.Fatalf("opening config file", err.Error())
	}

	objects := []model.FormSubmission{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		log.Fatalf("Error reading forms submissions ", err.Error())
	}

	f := getAForm(t)
	fId := f.ID.Hex()

	c := web.NewContext(nil, nil)
	defer c.Close()

	c.SetValue("form_id", fId)

	for _, one := range objects {
		c.Marshall(one)

		if _, err := service.CreateFormSubmission(c); err != nil {
			fmt.Println(err)
			t.Fail()
		}

		// let's create another just so we can delete it
		c.Marshall(one)
		s, err := service.CreateFormSubmission(c)
		if err != nil {

			fmt.Println(err)
			t.Fail()
		}

		c := web.NewContext(nil, nil)
		defer c.Close()
		c.SetValue("id", s.ID.Hex())

		err = service.DeleteFormSubmission(c)
		if err != nil {

			fmt.Println(err)
			t.Fail()
		}

	}

}

// let's test adding and removing answers to a gallery
//  Galleries are made up of
//  Answers in submissions to forms
func TestAddingAndRemovingAnswersToGallery(t *testing.T) {

	// so let's get a form
	f := getAForm(t)

	// one of its submissions
	s := getASubmissionToAForm(f, t)

	// and a gallery
	g := getAGalleryFormAForm(f, t)

	// set the values into a context
	c := web.NewContext(nil, nil)
	defer c.Close()

	c.SetValue("id", g.ID.Hex())
	c.SetValue("submission_id", s.ID.Hex())

	// and for each submission answer
	for _, i := range s.Answers {
		c.SetValue("answer_id", i.WidgetId)

		// do a complex dance to test each permutation
		_, err := service.RemoveAnswerFromFormGallery(c)
		if err == nil {
			log.Fatalln("We shouldn't be able to remove an answer that isn't in the gallery")
			t.Fail()
		}

		_, err = service.AddAnswerToFormGallery(c)
		if err != nil {
			log.Fatalln("We should be able to add an answer to an empty gallery.", err)
			t.Fail()
		}

		_, err = service.AddAnswerToFormGallery(c)
		if err == nil {
			log.Fatalln("We shouldn't be able to add an answer twice to a gallery.")
			t.Fail()
		}

		_, err = service.RemoveAnswerFromFormGallery(c)
		if err != nil {
			log.Fatalln("We should be able to remove an answer that's already in a gallery")
			t.Fail()
		}

		_, err = service.AddAnswerToFormGallery(c)
		if err != nil {
			log.Fatalln("Should be able to add an answer to a gallery after removing it")
			t.Fail()
		}
	}

}

func TestFlagFormSubmissions(t *testing.T) {

	// so let's get a form n' submission
	f := getAForm(t)
	s := getASubmissionToAForm(f, t)

	// new context for the form submission
	c := web.NewContext(nil, nil)
	defer c.Close()
	c.SetValue("id", s.ID.Hex())

	fCount := len(s.Flags)

	c.SetValue("flag", "test_the_flag")
	s, err := service.AddFlagToFormSubmission(c)
	if err != nil {
		log.Fatalln("Should be able to add a flag to a gallery after removing it")
		t.Fail()
	}

	if len(s.Flags) != (fCount + 1) {
		log.Fatalln("Flag count should increment after add")
		t.Fail()
	}

	s, err = service.AddFlagToFormSubmission(c)
	if len(s.Flags) != (fCount + 1) {
		log.Fatalln("Should not be able to add a flag twice")
		t.Fail()
	}

	c.SetValue("flag", "test_another__flag")
	s, err = service.AddFlagToFormSubmission(c)
	if err != nil {
		log.Fatalln("Should be able to add a flag to a gallery after removing it")
		t.Fail()
	}

	s, err = service.RemoveFlagFromFormSubmission(c)
	if err != nil {
		log.Fatalln("Should be able to remove a flag to a gallery after removing it")
		t.Fail()
	}

	if len(s.Flags) != (fCount + 1) {
		log.Fatalln("Should get the right count after adding and removing")
		t.Fail()
	}

	s, err = service.RemoveFlagFromFormSubmission(c)
	if err != nil {
		log.Fatalln("Should not be able to remove a flag twice")
		t.Fail()
	}

	if len(s.Flags) != (fCount + 1) {
		log.Fatalln("Should get the right count after adding and removing")
		t.Fail()
	}

}

func TestEditFormSubmissionAnswer(t *testing.T) {

	// so let's get a form
	f := getAForm(t)

	// one of its submissions
	s := getASubmissionToAForm(f, t)

	c := web.NewContext(nil, nil)
	c.SetValue("id", s.ID.Hex())

	for _, a := range s.Answers {
		// set the answer id into the context
		c.SetValue("answer_id", a.WidgetId)

		// set the context body to something ridiculous
		body := model.FormSubmissionEditInput{"This is an edit! Purple Monkey Dishwasher."}
		c.Marshall(body)

		// and edit the answer
		_, err := service.EditFormSubmissionAnswer(c)
		if err != nil {
			log.Fatalf("Could not edit a form submission answer: %+v", err)
		}

	}

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
		search, err := service.CreateUpdateSearch(c)
		if err != nil {
			t.Fail()
		}
		//make sure upsert on the same ID works
		c.Marshall(search)
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

	//Try again with the same data and it shouldn't fail
	for _, one := range objects {
		c.Marshall(one)
		if _, err := service.ImportAction(c); err != nil {
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
