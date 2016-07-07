package service_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/service"
	"github.com/coralproject/pillar/pkg/web"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Searching", func() {

	var (
		mongodb_url string
		result      []model.FormSubmission
		err         *web.AppError
	)

	BeforeEach(func() {

		// set test database
		mongodb_url = os.Getenv("MONGODB_URL")
		if mongodb_url == "" {
			log.Fatal("MONGODB_URL needs to be setup.")
		}
		e := os.Setenv("MONGODB_URL", mongodb_url+"_test") //os.Setenv("MONGODB_URL", mongodb_url+"_"+strconv.Itoa(rand.Intn(10))+"_test")
		if e != nil {
			fmt.Println("Error when setting environment test ", e)
		}

		// add submissions from fixtures
		loadformfixtures()
	})

	AfterEach(func() {
		// empty database
		emptydb()

		// restore initial database
		e := os.Setenv("MONGODB_URL", mongodb_url)
		if e != nil {
			fmt.Println("Error when setting environment back ", e)
		}
	})

	Describe("the form submissions ", func() {
		Context("with an existing string", func() {

			JustBeforeEach(func() {
				// prep a context for the search
				c := web.NewContext(nil, nil)
				defer c.Close()

				c.SetValue("search", "Gophers")

				result, err = service.SearchFormSubmissions(c)
			})

			It("should give back a result and no error", func() {
				Expect(len(result)).Should(Equal(1))
				Expect(err).Should(BeNil())
				Expect(find("Gophers", result)).Should(BeTrue())
			})
		})
	})

})

var _ = Describe("Get", func() {

	var (
		err         *web.AppError
		formid      string
		mongodb_url string
	)

	BeforeEach(func() {

		// set test database
		mongodb_url = os.Getenv("MONGODB_URL")
		if mongodb_url == "" {
			log.Fatal("MONGODB_URL needs to be setup.")
		}
		e := os.Setenv("MONGODB_URL", mongodb_url+"_test") //os.Setenv("MONGODB_URL", mongodb_url+"_"+strconv.Itoa(rand.Intn(10))+"_test")
		if e != nil {
			fmt.Println("Error when setting environment test ", e)
		}

		// get the fixtures for forms and forms submissions
		loadformfixtures()
		formid = "577c1c95a969c805f7f8c88a"

		loadformgalleriesfixtures()
	})

	AfterEach(func() {
		// empty database
		emptydb()

		// restore initial database
		e := os.Setenv("MONGODB_URL", mongodb_url)
		if e != nil {
			fmt.Println("Error when setting environment back ", e)
		}
	})

	Describe("forms", func() {
		Context("with appropiate context", func() {

			var fs []model.Form

			JustBeforeEach(func() {
				c := web.NewContext(nil, nil)
				defer c.Close()

				fs = []model.Form{}

				// let's see if we have forms to reply to
				fs, err = service.GetForms(c)
			})
			It("should return at least a form and no error", func() {
				Expect(len(fs)).ShouldNot(Equal(0))
				Expect(err).Should(BeNil())
			})
		})
	})

	Describe("submissions to a form", func() {
		Context("with appropiate context", func() {

			var fss []model.FormSubmission

			JustBeforeEach(func() {
				// create the context for this form
				c := web.NewContext(nil, nil)
				defer c.Close()
				c.SetValue("form_id", formid)

				fss, err = service.GetFormSubmissionsByForm(c)
			})
			It("should return at least a form and no error", func() {
				Expect(len(fss)).ShouldNot(Equal(0))
				Expect(err).Should(BeNil())
			})
		})
	})

	Describe("galleries to a form", func() {
		Context("with appropiate context", func() {

			var g []model.FormGallery

			JustBeforeEach(func() {
				// create the context for this form
				c := web.NewContext(nil, nil)
				defer c.Close()
				c.SetValue("form_id", formid)

				g, err = service.GetFormGalleriesByForm(c)

			})

			It("should return at least a gallery and no error", func() {
				Expect(len(g)).ShouldNot(Equal(0))
				Expect(err).Should(BeNil())
			})
		})
	})
})

var _ = Describe("Create", func() {

	var (
		mongodb_url string
		err         *web.AppError
		object      model.Form
		result      *model.Form
	)

	BeforeEach(func() {
		// set test database
		mongodb_url = os.Getenv("MONGODB_URL")
		if mongodb_url == "" {
			log.Fatal("MONGODB_URL needs to be setup.")
		}
		e := os.Setenv("MONGODB_URL", mongodb_url+"_test") //os.Setenv("MONGODB_URL", mongodb_url+"_"+strconv.Itoa(rand.Intn(10))+"_test")
		if e != nil {
			fmt.Println("Error when setting environment test ", e)
		}

		// get the fixtures from appropiate json file for data form submission
		file, e := ioutil.ReadFile(dataFormsIds)
		if e != nil {
			log.Fatalf("opening config file %v", e.Error())
		}

		e = json.Unmarshal(file, &object)
		if e != nil {
			log.Fatalf("Error reading forms. %v", e.Error())
		}

		c := web.NewContext(nil, nil)
		defer c.Close()

		c.Marshall(object)
		result, err = service.CreateUpdateForm(c)
	})

	AfterEach(func() {
		// empty database
		emptydb()

		// restore initial database
		e := os.Setenv("MONGODB_URL", mongodb_url)
		if e != nil {
			fmt.Println("Error when setting environment back ", e)
		}
	})

	Describe("a form", func() {
		Context("with appropiate context", func() {

			It("should not give an error", func() {
				Expect(err).Should(BeNil())
				Expect(result).ShouldNot(BeNil())
			})
		})
	})

	Describe("a form submission", func() {
		Context("with appropiate context", func() {

			var (
				objectS model.FormSubmission
			)

			JustBeforeEach(func() {

				// get the fixtures from appropiate json file for data form submission
				file, e := ioutil.ReadFile(dataFormSubmissionsIds)
				if e != nil {
					log.Fatalf("opening config file %v", e.Error())
				}

				e = json.Unmarshal(file, &objectS)
				if e != nil {
					log.Fatalf("Error reading forms. %v", e.Error())
				}

				c := web.NewContext(nil, nil)
				defer c.Close()

				c.SetValue("form_id", result.ID.Hex())
				c.Marshall(objectS)
				_, err = service.CreateFormSubmission(c)
			})

			It("should not give an error", func() {
				Expect(err).Should(BeNil())
			})
		})
	})
})

// var _ = Describe("Delete", func() {
//
// 	var (
// 		mongodb_url string
// 		err         *web.AppError
// 	)
//
// 	BeforeEach(func() {
// 		// set test database
// 		mongodb_url = os.Getenv("MONGODB_URL")
// 		if mongodb_url == "" {
// 			log.Fatal("MONGODB_URL needs to be setup.")
// 		}
// 		e := os.Setenv("MONGODB_URL", mongodb_url+"_test") //os.Setenv("MONGODB_URL", mongodb_url+"_"+strconv.Itoa(rand.Intn(10))+"_test")
// 		if e != nil {
// 			fmt.Println("Error when setting environment test ", e)
// 		}
//
// 		loadformfixtures()
// 	})
//
// 	AfterEach(func() {
// 		// empty database
// 		emptydb()
//
// 		// restore initial database
// 		e := os.Setenv("MONGODB_URL", mongodb_url)
// 		if e != nil {
// 			fmt.Println("Error when setting environment back ", e)
// 		}
// 	})
//
// 	Describe("a form", func() {
// 		Context("with appropiate context", func() {
//
// 			var (
// 				object model.Form
// 				formid string
// 			)
//
// 			JustBeforeEach(func() {
// 				// get the fixtures from appropiate json file for data form submission
// 				file, e := ioutil.ReadFile(dataFormsIds)
// 				if e != nil {
// 					log.Fatalf("opening config file %v", e.Error())
// 				}
//
// 				e = json.Unmarshal(file, &object)
// 				if e != nil {
// 					log.Fatalf("Error reading forms. %v", e.Error())
// 				}
//
// 				formid = "577c1c95a969c805f7f8c88a"
//
// 				c := web.NewContext(nil, nil)
// 				defer c.Close()
//
// 				c.SetValue("formid", formid)
//
// 				c.Marshall(object)
// 				r, _ := service.CreateUpdateForm(c)
//
// 				co := web.NewContext(nil, nil)
// 				defer co.Close()
//
// 				co.SetValue("id", r.ID.Hex())
//
// 				err = service.DeleteFormSubmission(co)
// 			})
//
// 			It("should not give an error", func() {
// 				Expect(err).Should(BeNil())
// 			})
// 		})
// 	})
// })

// func TestFormSubmissionsSchenanigans(t *testing.T) {
//
// 	file, err := os.Open(dataFormSubmissions)
// 	if err != nil {
// 		log.Fatalf("opening config file %v", err.Error())
// 	}
//
// 	objects := []model.FormSubmission{}
// 	jsonParser := json.NewDecoder(file)
// 	if err = jsonParser.Decode(&objects); err != nil {
// 		log.Fatalf("Error reading forms submissions %v", err.Error())
// 	}
//
// 	f := getAForm(t)
// 	fID := f.ID.Hex()
//
// 	c := web.NewContext(nil, nil)
// 	defer c.Close()
//
// 	c.SetValue("form_id", fID)
//
// 	for _, one := range objects {
// 		c.Marshall(one)
//
// 		if _, err := service.CreateFormSubmission(c); err != nil {
// 			fmt.Println(err)
// 			t.Fail()
// 		}
//
// 		// let's create another just so we can delete it
// 		c.Marshall(one)
// 		s, err := service.CreateFormSubmission(c)
// 		if err != nil {
//
// 			fmt.Println(err)
// 			t.Fail()
// 		}
//
// 		c := web.NewContext(nil, nil)
// 		defer c.Close()
// 		c.SetValue("id", s.ID.Hex())
//
// 		err = service.DeleteFormSubmission(c)
// 		if err != nil {
//
// 			fmt.Println(err)
// 			t.Fail()
// 		}
//
// 	}
//
// }
//
// //
// // // let's test adding and removing answers to a gallery
// // //  Galleries are made up of
// // //  Answers in submissions to forms
// // func TestAddingAndRemovingAnswersToGallery(t *testing.T) {
// //
// // 	// so let's get a form
// // 	f := getAForm(t)
// //
// // 	// one of its submissions
// // 	s := getASubmissionToAForm(f, t)
// //
// // 	// and a gallery
// // 	g := getAGalleryFormAForm(f, t)
// //
// // 	// set the values into a context
// // 	c := web.NewContext(nil, nil)
// // 	defer c.Close()
// //
// // 	c.SetValue("id", g.ID.Hex())
// // 	c.SetValue("submission_id", s.ID.Hex())
// //
// // 	// and for each submission answer
// // 	for _, i := range s.Answers {
// // 		c.SetValue("answer_id", i.WidgetId)
// //
// // 		// do a complex dance to test each permutation
// // 		_, err := service.RemoveAnswerFromFormGallery(c)
// // 		if err == nil {
// // 			log.Fatalln("We shouldn't be able to remove an answer that isn't in the gallery")
// // 			t.Fail()
// // 		}
// //
// // 		_, err = service.AddAnswerToFormGallery(c)
// // 		if err != nil {
// // 			log.Fatalln("We should be able to add an answer to an empty gallery.", err)
// // 			t.Fail()
// // 		}
// //
// // 		_, err = service.AddAnswerToFormGallery(c)
// // 		if err == nil {
// // 			log.Fatalln("We shouldn't be able to add an answer twice to a gallery.")
// // 			t.Fail()
// // 		}
// //
// // 		_, err = service.RemoveAnswerFromFormGallery(c)
// // 		if err != nil {
// // 			log.Fatalln("We should be able to remove an answer that's already in a gallery")
// // 			t.Fail()
// // 		}
// //
// // 		_, err = service.AddAnswerToFormGallery(c)
// // 		if err != nil {
// // 			log.Fatalln("Should be able to add an answer to a gallery after removing it")
// // 			t.Fail()
// // 		}
// // 	}
// //
// // }
//
// func TestFlagFormSubmissions(t *testing.T) {
//
// 	// so let's get a form n' submission
// 	f := getAForm(t)
// 	s := getASubmissionToAForm(f, t)
//
// 	// new context for the form submission
// 	c := web.NewContext(nil, nil)
// 	defer c.Close()
// 	c.SetValue("id", s.ID.Hex())
//
// 	fCount := len(s.Flags)
//
// 	c.SetValue("flag", "test_the_flag")
// 	s, err := service.AddFlagToFormSubmission(c)
// 	if err != nil {
// 		log.Fatalln("Should be able to add a flag to a gallery after removing it")
// 		t.Fail()
// 	}
//
// 	if len(s.Flags) != (fCount + 1) {
// 		log.Fatalln("Flag count should increment after add")
// 		t.Fail()
// 	}
//
// 	s, err = service.AddFlagToFormSubmission(c)
// 	if len(s.Flags) != (fCount + 1) {
// 		log.Fatalln("Should not be able to add a flag twice")
// 		t.Fail()
// 	}
//
// 	c.SetValue("flag", "test_another__flag")
// 	s, err = service.AddFlagToFormSubmission(c)
// 	if err != nil {
// 		log.Fatalln("Should be able to add a flag to a gallery after removing it")
// 		t.Fail()
// 	}
//
// 	s, err = service.RemoveFlagFromFormSubmission(c)
// 	if err != nil {
// 		log.Fatalln("Should be able to remove a flag to a gallery after removing it")
// 		t.Fail()
// 	}
//
// 	if len(s.Flags) != (fCount + 1) {
// 		log.Fatalln("Should get the right count after adding and removing")
// 		t.Fail()
// 	}
//
// 	s, err = service.RemoveFlagFromFormSubmission(c)
// 	if err != nil {
// 		log.Fatalln("Should not be able to remove a flag twice")
// 		t.Fail()
// 	}
//
// 	if len(s.Flags) != (fCount + 1) {
// 		log.Fatalln("Should get the right count after adding and removing")
// 		t.Fail()
// 	}
//
// }
//
// func TestEditFormSubmissionAnswer(t *testing.T) {
//
// 	// so let's get a form
// 	f := getAForm(t)
//
// 	// one of its submissions
// 	s := getASubmissionToAForm(f, t)
//
// 	c := web.NewContext(nil, nil)
// 	c.SetValue("id", s.ID.Hex())
//
// 	for _, a := range s.Answers {
// 		// set the answer id into the context
// 		c.SetValue("answer_id", a.WidgetId)
//
// 		// set the context body to something ridiculous
// 		body := model.FormSubmissionEditInput{"This is an edit! Purple Monkey Dishwasher."}
// 		c.Marshall(body)
//
// 		// and edit the answer
// 		_, err := service.EditFormSubmissionAnswer(c)
// 		if err != nil {
// 			log.Fatalf("Could not edit a form submission answer: %+v", err)
// 		}
//
// 	}
//
// }
//
// func TestCreateSections(t *testing.T) {
// 	file, err := os.Open(dataSections)
// 	if err != nil {
// 		log.Fatalf("opening config file %v", err.Error())
// 	}
//
// 	objects := []model.Section{}
// 	jsonParser := json.NewDecoder(file)
// 	if err = jsonParser.Decode(&objects); err != nil {
// 		log.Fatalf("Error reading tags %v", err.Error())
// 	}
//
// 	c := web.NewContext(nil, nil)
// 	defer c.Close()
//
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if _, err := service.CreateUpdateSection(c); err != nil {
// 			t.Fail()
// 		}
// 	}
// }
//
// func TestCreateTags(t *testing.T) {
// 	file, err := os.Open(dataTags)
// 	if err != nil {
// 		log.Fatalf("opening config file %v", err.Error())
// 	}
//
// 	objects := []model.Tag{}
// 	jsonParser := json.NewDecoder(file)
// 	if err = jsonParser.Decode(&objects); err != nil {
// 		log.Fatalf("Error reading tags %v", err.Error())
// 	}
//
// 	c := web.NewContext(nil, nil)
// 	defer c.Close()
//
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if _, err := service.CreateUpdateTag(c); err != nil {
// 			t.Fail()
// 		}
// 	}
// }
//
// func TestCreateSearches(t *testing.T) {
// 	file, err := os.Open(dataSearches)
// 	if err != nil {
// 		log.Fatalf("opening config file %v", err.Error())
// 	}
//
// 	objects := []model.Search{}
// 	jsonParser := json.NewDecoder(file)
// 	if err = jsonParser.Decode(&objects); err != nil {
// 		log.Fatalf("Error reading userGroups %v", err.Error())
// 	}
//
// 	c := web.NewContext(nil, nil)
// 	defer c.Close()
//
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		search, err := service.CreateUpdateSearch(c)
// 		if err != nil {
// 			t.Fail()
// 		}
// 		//make sure upsert on the same ID works
// 		c.Marshall(search)
// 		if _, err := service.CreateUpdateSearch(c); err != nil {
// 			t.Fail()
// 		}
// 	}
// }
//
// func TestImportAssets(t *testing.T) {
// 	file, err := os.Open(dataAssets)
// 	if err != nil {
// 		log.Fatalf("opening config file %v", err.Error())
// 	}
//
// 	objects := []model.Asset{}
// 	jsonParser := json.NewDecoder(file)
// 	if err = jsonParser.Decode(&objects); err != nil {
// 		log.Fatalf("Error reading assets data %v", err.Error())
// 	}
//
// 	c := web.NewContext(nil, nil)
// 	defer c.Close()
//
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if _, err := service.ImportAsset(c); err != nil {
// 			t.Fail()
// 		}
// 	}
//
// 	//try the same set again and it shouldn't fail
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if _, err := service.ImportAsset(c); err != nil {
// 			t.Fail()
// 		}
// 	}
// }
//
// func TestImportUsers(t *testing.T) {
// 	file, err := os.Open(dataUsers)
// 	if err != nil {
// 		log.Fatalf("opening config file %v", err.Error())
// 	}
//
// 	objects := []model.User{}
// 	jsonParser := json.NewDecoder(file)
// 	if err = jsonParser.Decode(&objects); err != nil {
// 		log.Fatalf("Error reading users data %v", err.Error())
// 	}
//
// 	c := web.NewContext(nil, nil)
// 	defer c.Close()
//
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if _, err := service.ImportUser(c); err != nil {
// 			t.Fail()
// 		}
// 	}
// 	//try the same set again and it shouldn't fail
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if _, err := service.ImportUser(c); err != nil {
// 			t.Fail()
// 		}
// 	}
// }
//
// func TestImportComments(t *testing.T) {
// 	file, err := os.Open(dataComments)
// 	if err != nil {
// 		log.Fatalf("opening config file %v", err.Error())
// 	}
//
// 	objects := []model.Comment{}
// 	jsonParser := json.NewDecoder(file)
// 	if err = jsonParser.Decode(&objects); err != nil {
// 		log.Fatalf("Error reading comments data %v", err.Error())
// 	}
//
// 	c := web.NewContext(nil, nil)
// 	defer c.Close()
//
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if _, err := service.ImportComment(c); err != nil {
// 			log.Fatalf("Error: %v\n", err)
// 			t.Fail()
// 		}
// 	}
// 	//try the same set again and it shouldn't fail
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if _, err := service.ImportComment(c); err != nil {
// 			log.Fatalf("Error: %v\n", err)
// 			t.Fail()
// 		}
// 	}
// }
//
// func TestImportActions(t *testing.T) {
// 	file, err := os.Open(dataActions)
// 	if err != nil {
// 		log.Fatalf("opening config file %v", err.Error())
// 	}
//
// 	objects := []model.Action{}
// 	jsonParser := json.NewDecoder(file)
// 	if err = jsonParser.Decode(&objects); err != nil {
// 		log.Fatalf("Error reading user data %v", err.Error())
// 	}
//
// 	c := web.NewContext(nil, nil)
// 	defer c.Close()
//
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if _, err := service.ImportAction(c); err != nil {
// 			log.Fatalf("Error: %v\n", err)
// 			t.Fail()
// 		}
// 	}
//
// 	//Try again with the same data and it shouldn't fail
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if _, err := service.ImportAction(c); err != nil {
// 			log.Fatalf("Error: %v\n", err)
// 			t.Fail()
// 		}
// 	}
// }
//
// func TestCreateIndexes(t *testing.T) {
// 	file, err := os.Open(dataIndexes)
// 	if err != nil {
// 		log.Fatalf("opening config file %v", err.Error())
// 	}
//
// 	objects := []model.Index{}
// 	jsonParser := json.NewDecoder(file)
// 	if err = jsonParser.Decode(&objects); err != nil {
// 		log.Fatalf("Error reading index data %v", err.Error())
// 	}
//
// 	c := web.NewContext(nil, nil)
// 	defer c.Close()
//
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if err := service.CreateIndex(c); err != nil {
// 			log.Fatalf("Error: %v\n", err)
// 			t.Fail()
// 		}
// 	}
// }
//
// func TestUpdateMetadata(t *testing.T) {
// 	file, err := os.Open(dataMetadata)
// 	if err != nil {
// 		log.Fatalf("opening config file %v", err.Error())
// 	}
//
// 	objects := []model.Metadata{}
// 	jsonParser := json.NewDecoder(file)
// 	if err = jsonParser.Decode(&objects); err != nil {
// 		log.Fatalf("Error reading metadata %v", err.Error())
// 	}
//
// 	c := web.NewContext(nil, nil)
// 	defer c.Close()
//
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if _, err := service.UpdateMetadata(c); err != nil {
// 			log.Fatalf("Error: %v\n", err)
// 			t.Fail()
// 		}
// 	}
// }
//
// func TestUserActions(t *testing.T) {
// 	file, err := os.Open(dataUserActions)
// 	if err != nil {
// 		log.Fatalf("opening config file %v", err.Error())
// 	}
//
// 	objects := []model.CayUserAction{}
// 	jsonParser := json.NewDecoder(file)
// 	if err = jsonParser.Decode(&objects); err != nil {
// 		log.Fatalf("Error reading user-actions %v", err.Error())
// 	}
//
// 	c := web.NewContext(nil, nil)
// 	defer c.Close()
//
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if err := service.CreateUserAction(c); err != nil {
// 			log.Fatalf("Error: %v\n", err)
// 			t.Fail()
// 		}
// 	}
// }
//
// func TestRenameTags(t *testing.T) {
// 	file, err := os.Open(dataNewTags)
// 	if err != nil {
// 		log.Fatalf("opening config file %v", err.Error())
// 	}
//
// 	objects := []model.Tag{}
// 	jsonParser := json.NewDecoder(file)
// 	if err = jsonParser.Decode(&objects); err != nil {
// 		log.Fatalf("Error reading tags %v", err.Error())
// 	}
//
// 	c := web.NewContext(nil, nil)
// 	defer c.Close()
//
// 	for _, one := range objects {
// 		c.Marshall(one)
// 		if _, err := service.CreateUpdateTag(c); err != nil {
// 			log.Fatalf("Error: %v\n", err)
// 			t.Fail()
// 		}
// 	}
// }
//
// //func TestDeleteAllTag(t *testing.T) {
// //	tags, err := GetTags()
// //	if err != nil || len(tags) == 0 {
// //		t.Fail()
// //	}
// //
// //	for _, one := range tags {
// //		err := DeleteTag(&one)
// //		if err != nil {
// //			t.Fail()
// //		}
// //	}
// //
// //	objects, err := GetTags()
// //	if err != nil || len(objects) != 0 {
// //		t.Fail()
// //	}
// //}
