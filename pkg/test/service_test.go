package test

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/service"
	"github.com/coralproject/pillar/pkg/web"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Create", func() {

	var (
		err     *web.AppError
		errS    *web.AppError
		object  model.Form
		result  *model.Form
		objectS model.FormSubmission
		resultS *model.FormSubmission
	)

	BeforeEach(func() {
		setTestDatabase()

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

		// get the fixtures from appropiate json file for data form submission
		objectS = getDataFormSubmissions(dataFormSubmissionsIds)

		c = web.NewContext(nil, nil)
		defer c.Close()

		c.SetValue("form_id", result.ID.Hex())
		c.Marshall(objectS)
		resultS, errS = service.CreateFormSubmission(c)
	})

	AfterEach(func() {
		// empty database
		emptyDB()

		// recover MONGODB_URL
		recoverEnvVariables()
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

			It("should not give an error", func() {
				Expect(errS).Should(BeNil())
			})
		})
	})

	Describe("a form and then delete", func() {
		Context("with appropiate context", func() {

			JustBeforeEach(func() {
				co := web.NewContext(nil, nil)
				defer co.Close()

				co.SetValue("id", resultS.ID.Hex())

				err = service.DeleteFormSubmission(co)
			})

			It("should not give an error", func() {
				Expect(err).Should(BeNil())
			})
		})
	})

	Describe("an answer to a gallery", func() {
		Context("with appropiate context", func() {

			var (
				gallery model.FormGallery
				context *web.AppContext
				err     *web.AppError
			)

			JustBeforeEach(func() {

				loadformgalleriesfixtures()

				// create a gallery
				c := web.NewContext(nil, nil)
				defer c.Close()
				c.SetValue("form_id", result.ID.Hex())

				var g []model.FormGallery
				g, err = service.GetFormGalleriesByForm(c)

				gallery = g[0]

				// set the values into a context
				context = web.NewContext(nil, nil)
				defer context.Close()

				context.SetValue("id", gallery.ID.Hex())
				context.SetValue("submission_id", resultS.ID.Hex())
				c.SetValue("answer_id", resultS.Answers[0].WidgetId)

			})

			It("should return a gallery", func() {
				Expect(err).Should(BeNil(), "Could not load forms galleries")
				Expect(gallery).ShouldNot(BeNil())
			})

			It("should not be able to have an answer that is not in the gallery removed", func() {
				_, e := service.RemoveAnswerFromFormGallery(context)
				Expect(e).ShouldNot(BeNil())
			})

			It("should be able to add an answer to an empty gallery", func() {
				_, e := service.AddAnswerToFormGallery(context)
				Expect(e).Should(BeNil())
			})

			It("shouldn't be able to add an answer twice to a gallery.", func() {
				service.AddAnswerToFormGallery(context)
				_, e := service.AddAnswerToFormGallery(context)
				Expect(e).ShouldNot(BeNil())
			})

			It("should be able to remove an answer that's already in a gallery", func() {
				service.AddAnswerToFormGallery(context)
				_, e := service.RemoveAnswerFromFormGallery(context)
				Expect(e).Should(BeNil())
			})

			It("should be able to add an answer to a gallery after removing it", func() {
				service.AddAnswerToFormGallery(context)
				service.RemoveAnswerFromFormGallery(context)
				_, e := service.AddAnswerToFormGallery(context)
				Expect(e).Should(BeNil())
			})
		})
	})

	Describe("a section", func() {
		Context("with appropiate context", func() {
			var (
				objects []model.Section
				err     *web.AppError
			)

			JustBeforeEach(func() {
				// get the fixtures from appropiate json file for data sections
				file, e := ioutil.ReadFile(dataSections)
				if e != nil {
					log.Fatalf("opening config file %v", e.Error())
				}

				e = json.Unmarshal(file, &objects)
				if e != nil {
					log.Fatalf("Error reading sections. %v", e.Error())
				}

				c := web.NewContext(nil, nil)
				defer c.Close()

				c.Marshall(objects[0])

				_, err = service.CreateUpdateSection(c)
			})

			It("should return no error", func() {
				Expect(err).Should(BeNil(), "Could not load sections")
			})
		})
	})

	Describe("a tag", func() {
		Context("with appropiate context", func() {
			var (
				objects []model.Tag
				err     *web.AppError
			)

			JustBeforeEach(func() {
				// get the fixtures from appropiate json file for data tags
				objects = getDataTags(dataTags)

				c := web.NewContext(nil, nil)
				defer c.Close()

				for _, o := range objects {
					c.Marshall(o)
					_, err = service.CreateUpdateTag(c)
					Expect(err).Should(BeNil(), "Could not load tags")
				}
			})

			It("should rename without error", func() {
				// get the fixtures from appropiate json file for data tags
				objects = getDataTags(dataNewTags)

				c := web.NewContext(nil, nil)
				defer c.Close()

				for _, one := range objects {
					c.Marshall(one)
					_, err := service.CreateUpdateTag(c)
					Expect(err).Should(BeNil())
				}
			})

			It("should delete tags without error", func() {

				c := web.NewContext(nil, nil)
				defer c.Close()

				for _, o := range objects {
					c.Marshall(o)
					err := service.DeleteTag(c)
					Expect(err).Should(BeNil())
				}
				//
				// //	objects, err := GetTags()
				// //	if err != nil || len(objects) != 0 {
				// //		t.Fail()
				// //	}
				// //}
			})
		})
	})

	Describe("a search", func() {
		Context("with appropiate context", func() {
			var (
				objects []model.Search
				err     *web.AppError
				errS    *web.AppError
				search  *model.Search
			)

			JustBeforeEach(func() {
				// get the fixtures from appropiate json file for data searchs
				file, e := ioutil.ReadFile(dataSearches)
				if e != nil {
					log.Fatalf("opening config file %v", e.Error())
				}

				e = json.Unmarshal(file, &objects)
				if e != nil {
					log.Fatalf("Error reading searches. %v", e.Error())
				}

				c := web.NewContext(nil, nil)
				defer c.Close()

				c.Marshall(objects[0])

				search, err = service.CreateUpdateSearch(c)

				c.Marshall(search)
				_, errS = service.CreateUpdateSearch(c)
			})

			It("should return no error", func() {
				Expect(err).Should(BeNil(), "Could not load searchs")
			})

			It("should make sure upsert on the same ID works", func() {
				Expect(errS).Should(BeNil())
			})
		})
	})

	Describe("an index", func() {
		Context("with appropiate context", func() {
			var (
				objects []model.Index
				err     *web.AppError
			)

			JustBeforeEach(func() {
				// get the fixtures from appropiate json file for data indexes
				file, e := ioutil.ReadFile(dataIndexes)
				if e != nil {
					log.Fatalf("opening config file %v", e.Error())
				}

				e = json.Unmarshal(file, &objects)
				if e != nil {
					log.Fatalf("Error reading indexes. %v", e.Error())
				}

				c := web.NewContext(nil, nil)
				defer c.Close()

				c.Marshall(objects[0])

				err = service.CreateIndex(c)
			})

			It("should return no error", func() {
				Expect(err).Should(BeNil(), "Could not load searchs")
			})
		})
	})

	Describe("user actions", func() {
		Context("with appropiate context", func() {
			var (
				usera []model.CayUserAction
				err   *web.AppError
			)

			JustBeforeEach(func() {

				file, e := ioutil.ReadFile(dataUserActions)
				if e != nil {
					log.Fatalf("opening config file %v", e.Error())
				}

				e = json.Unmarshal(file, &usera)
				if e != nil {
					log.Fatalf("Error reading user actions. %v", e.Error())
				}

				c := web.NewContext(nil, nil)
				defer c.Close()

				for _, one := range usera {
					c.Marshall(one)
					err = service.CreateUserAction(c)
				}
			})

			It("should return no error", func() {
				Expect(err).Should(BeNil(), "Could not load searchs")
			})
		})
	})
})

var _ = Describe("Get", func() {

	var (
		err    *web.AppError
		formid string
	)

	BeforeEach(func() {
		setTestDatabase()

		// get the fixtures for forms and forms submissions
		loadformfixtures()
		formid = "577c18f4a969c805f7f8c889"

		loadformgalleriesfixtures()
	})

	AfterEach(func() {
		// empty database
		emptyDB()

		// recover MONGODB_URL
		recoverEnvVariables()
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
			It("should return at least a submission to a form and no error", func() {
				Expect(err).Should(BeNil())
				Expect(len(fss)).ShouldNot(Equal(0))
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

var _ = Describe("Search", func() {

	const (
		SEARCH = "Gophers"
	)

	var (
		result []model.FormSubmission
		err    *web.AppError
	)

	BeforeEach(func() {
		setTestDatabase()

		// add submissions from fixtures
		loadformfixtures()

		// prep a context for the search
		c := web.NewContext(nil, nil)
		defer c.Close()

		c.SetValue("search", SEARCH)

		result, err = service.SearchFormSubmissions(c)
	})

	AfterEach(func() {
		// empty database
		emptyDB()

		// recover MONGODB_URL
		recoverEnvVariables()
	})

	Describe("the form submissions ", func() {
		Context("with an existing string", func() {
			It("should give no error and give back the result we are looking for", func() {
				Expect(err).Should(BeNil())
				Expect(find(SEARCH, result)).Should(BeTrue())
			})
		})
	})

})

var _ = Describe("Flag", func() {

	var (
		err     *web.AppError
		object  model.Form
		result  *model.Form
		objectS model.FormSubmission
		resultS *model.FormSubmission
	)

	BeforeEach(func() {
		setTestDatabase()

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

		objectS = getDataFormSubmissions(dataFormSubmissionsIds)

		c = web.NewContext(nil, nil)
		defer c.Close()

		c.SetValue("form_id", result.ID.Hex())
		c.Marshall(objectS)
		resultS, err = service.CreateFormSubmission(c)
	})

	AfterEach(func() {
		// empty database
		emptyDB()

		// recover MONGODB_URL
		recoverEnvVariables()
	})

	Describe("a form submission", func() {
		Context("with appropiate context", func() {

			var (
				sub     *model.FormSubmission
				fCount  int
				context *web.AppContext
			)
			JustBeforeEach(func() {
				context = web.NewContext(nil, nil)
				defer context.Close()

				context.SetValue("id", resultS.ID.Hex())

				fCount = len(resultS.Flags)
			})

			It("should not give an error and should increment flag count after add", func() {
				context.SetValue("flag", "test_the_flag")
				sub, err = service.AddFlagToFormSubmission(context)

				Expect(err).Should(BeNil())
				Expect(len(sub.Flags)).Should(Equal(fCount + 1))
			})

			It("should not be able to add a flag twice", func() {
				context.SetValue("flag", "test_the_flag")
				service.AddFlagToFormSubmission(context)
				sub, err = service.AddFlagToFormSubmission(context)
				Expect(err).Should(BeNil())
				Expect(len(sub.Flags)).Should(Equal(fCount + 1))
			})

			It("should be able to add a second flag", func() {
				context.SetValue("flag", "test_the_flag")
				service.AddFlagToFormSubmission(context)
				context.SetValue("flag", "test_another__flag")
				_, err = service.AddFlagToFormSubmission(context)
				Expect(err).Should(BeNil())
			})

			It("should be able to remove a flag to a gallery after removing it", func() {
				context.SetValue("flag", "test_the_flag")
				service.AddFlagToFormSubmission(context)
				_, err = service.RemoveFlagFromFormSubmission(context)
				Expect(err).Should(BeNil())
			})

			It("should get the right count after adding and removing", func() {
				context.SetValue("flag", "test_the_flag")
				service.AddFlagToFormSubmission(context)
				sub, _ = service.RemoveFlagFromFormSubmission(context)
				Expect(len(sub.Flags)).ShouldNot(Equal(fCount + 1))
			})

			It("should not be able to remove a flag twice", func() {
				context.SetValue("flag", "test_the_flag")
				service.AddFlagToFormSubmission(context)
				service.RemoveFlagFromFormSubmission(context)
				sub, err = service.RemoveFlagFromFormSubmission(context)
				Expect(err).Should(BeNil())
				Expect(len(sub.Flags)).Should(Equal(fCount))
			})

		})
	})
})

var _ = Describe("Edit", func() {

	var (
		err     *web.AppError
		object  model.Form
		result  *model.Form
		objectS model.FormSubmission
		resultS *model.FormSubmission
	)

	BeforeEach(func() {
		setTestDatabase()

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

		// get the fixtures from appropiate json file for data form submission
		file, e = ioutil.ReadFile(dataFormSubmissionsIds)
		if e != nil {
			log.Fatalf("opening config file %v", e.Error())
		}

		e = json.Unmarshal(file, &objectS)
		if e != nil {
			log.Fatalf("Error reading forms. %v", e.Error())
		}

		c = web.NewContext(nil, nil)
		defer c.Close()

		c.SetValue("form_id", result.ID.Hex())
		c.Marshall(objectS)
		resultS, err = service.CreateFormSubmission(c)

	})

	AfterEach(func() {
		// empty database
		emptyDB()

		// recover MONGODB_URL
		recoverEnvVariables()
	})

	Describe("a form submission answer", func() {
		Context("with appropiate context", func() {

			var (
				context *web.AppContext
			)
			JustBeforeEach(func() {
				context = web.NewContext(nil, nil)
				defer context.Close()

				context.SetValue("id", resultS.ID.Hex())

			})

			It("should not edit a form submission answer", func() {
				// set the answer id into the context

				context.SetValue("answer_id", resultS.Answers[0].WidgetId)

				// set the context body to something ridiculous
				body := model.FormSubmissionEditInput{EditedAnswer: "This is an edit! Purple Monkey Dishwasher."}
				context.Marshall(body)

				// and edit the answer
				_, err = service.EditFormSubmissionAnswer(context)

				Expect(err).Should(BeNil())
			})

		})
	})
})

var _ = Describe("Import", func() {

	var (
		err     *web.AppError
		objects []model.Asset
	)

	BeforeEach(func() {
		setTestDatabase()

		// get the fixtures from appropiate json file for data assets
		file, e := ioutil.ReadFile(dataAssets)
		if e != nil {
			log.Fatalf("opening config file %v", e.Error())
		}

		e = json.Unmarshal(file, &objects)
		if e != nil {
			log.Fatalf("Error reading assets. %v", e.Error())
		}

		c := web.NewContext(nil, nil)
		defer c.Close()

		c.Marshall(objects[0])

		_, err = service.ImportAsset(c)
	})

	AfterEach(func() {
		// empty database
		emptyDB()

		// recover MONGODB_URL
		recoverEnvVariables()
	})

	Describe("an asset", func() {
		Context("with appropiate context", func() {
			It("should not return an Error", func() {
				Expect(err).Should(BeNil())
			})
		})
	})
})

var _ = Describe("Import", func() {

	var (
		err     *web.AppError
		objects []model.User
		c       *web.AppContext
	)

	BeforeEach(func() {
		setTestDatabase()

		// get the fixtures from appropiate json file for data users
		file, e := ioutil.ReadFile(dataUsers)
		if e != nil {
			log.Fatalf("opening config file %v", e.Error())
		}

		e = json.Unmarshal(file, &objects)
		if e != nil {
			log.Fatalf("Error reading assets. %v", e.Error())
		}

		c = web.NewContext(nil, nil)
		defer c.Close()

		for _, i := range objects {
			c.Marshall(i)

			_, err = service.ImportUser(c)
		}
	})

	AfterEach(func() {
		// empty database
		emptyDB()

		// recover MONGODB_URL
		recoverEnvVariables()
	})

	Describe("a user", func() {
		Context("with appropiate context", func() {
			It("should not return an Error", func() {
				Expect(err).Should(BeNil())
			})
		})
	})

	Describe("a user again", func() {
		Context("with appropiate context", func() {
			It("should not return an Error", func() {
				c.Marshall(objects[0])
				_, err = service.ImportUser(c)
				Expect(err).Should(BeNil())
			})
		})
	})

	Describe("a comment", func() {
		Context("with appropiate context", func() {
			var (
				comments []model.Comment
				assets   []model.Asset
			)
			JustBeforeEach(func() {

				// get the fixtures for the assets

				// get the fixtures from appropiate json file for data assets
				file, e := ioutil.ReadFile(dataAssets)
				if e != nil {
					log.Fatalf("opening config file %v", e.Error())
				}

				e = json.Unmarshal(file, &assets)
				if e != nil {
					log.Fatalf("Error reading assets. %v", e.Error())
				}

				c1 := web.NewContext(nil, nil)
				defer c1.Close()

				c1.Marshall(assets[0])

				service.ImportAsset(c1)

				// get the fixtures from appropiate json file for data comment
				file, e = ioutil.ReadFile(dataComments)
				if e != nil {
					log.Fatalf("opening config file %v", e.Error())
				}

				e = json.Unmarshal(file, &comments)
				if e != nil {
					log.Fatalf("Error reading comments. %v", e.Error())
				}

				c = web.NewContext(nil, nil)
				defer c.Close()

				c.Marshall(comments[0])

				_, err = service.ImportComment(c)
			})

			It("should not return an error", func() {
				Expect(err).Should(BeNil())
			})

			It("should not return an error to import the comment again", func() {
				c.Marshall(comments[0])
				_, err = service.ImportComment(c)
				Expect(err).Should(BeNil())
			})
		})
	})

	Describe("an action", func() {
		Context("with appropiate context", func() {

			var (
				actions  []model.Action
				assets   []model.Asset
				comments []model.Comment
			)
			JustBeforeEach(func() {

				// get the fixtures from appropiate json file for data assets
				file, e := ioutil.ReadFile(dataAssets)
				if e != nil {
					log.Fatalf("opening config file %v", e.Error())
				}

				e = json.Unmarshal(file, &assets)
				if e != nil {
					log.Fatalf("Error reading assets. %v", e.Error())
				}

				c1 := web.NewContext(nil, nil)
				defer c1.Close()

				for _, i := range assets {
					c1.Marshall(i)
					_, e := service.ImportAsset(c1)
					Expect(e).Should(BeNil())
				}

				// get the fixtures from appropiate json file for data comment
				file, e = ioutil.ReadFile(dataComments)
				if e != nil {
					log.Fatalf("opening config file %v", e.Error())
				}

				e = json.Unmarshal(file, &comments)
				if e != nil {
					log.Fatalf("Error reading comments. %v", e.Error())
				}

				c2 := web.NewContext(nil, nil)
				defer c2.Close()

				for _, i := range comments {
					c2.Marshall(i)
					_, e := service.ImportComment(c2)
					Expect(e).Should(BeNil())
				}

				// get the fixtures from appropiate json file for data actions
				file, e = ioutil.ReadFile(dataActions)
				if e != nil {
					log.Fatalf("opening config file %v", e.Error())
				}

				e = json.Unmarshal(file, &actions)
				if e != nil {
					log.Fatalf("Error reading actions. %v", e.Error())
				}

				c = web.NewContext(nil, nil)
				defer c.Close()

				c.Marshall(actions[0])

				_, err = service.ImportAction(c)
			})

			It("should not return an error", func() {
				Expect(err).Should(BeNil())
			})

			It("and an action again should not return an error", func() {
				c.Marshall(actions[0])
				_, err = service.ImportAction(c)
				Expect(err).Should(BeNil())
			})
		})
	})

})

var _ = Describe("Update", func() {
	var (
		objects []model.Metadata
	)
	BeforeEach(func() {
		setTestDatabase()

		// get the fixtures from appropiate json file for data assets
		file, e := ioutil.ReadFile(dataAssets)
		if e != nil {
			log.Fatalf("opening config file %v", e.Error())
		}

		assets := []model.Asset{}
		e = json.Unmarshal(file, &assets)
		if e != nil {
			log.Fatalf("Error reading assets. %v", e.Error())
		}

		c1 := web.NewContext(nil, nil)
		defer c1.Close()

		for _, i := range assets {
			c1.Marshall(i)
			service.ImportAsset(c1)
		}

		objects = getMetadata(dataMetadata)

	})
	AfterEach(func() {
		// empty database
		emptyDB()

		// recover MONGODB_URL
		recoverEnvVariables()
	})

	Describe("metadata", func() {
		Context("with appropiate context", func() {
			It("should not return an Error", func() {
				c := web.NewContext(nil, nil)
				defer c.Close()

				for _, one := range objects {
					c.Marshall(one)
					_, err := service.UpdateMetadata(c)
					Expect(err).Should(BeNil())
				}
			})
		})
	})
})
