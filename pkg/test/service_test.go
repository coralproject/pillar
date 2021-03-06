package test

import (
	"encoding/json"
	"fmt"
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
		objectS []model.FormSubmission
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
		c.Marshall(objectS[0])
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
				expectedDescription := "of the rest of your life"
				header := result.Header.(map[string]interface{})
				Expect(header["description"]).Should(Equal(expectedDescription))
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

var _ = Describe("Save", func() {

	var (
		object model.Form
		result *model.Form
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
		result, _ = service.CreateUpdateForm(c)
	})

	AfterEach(func() {
		// empty database
		emptyDB()

		// recover MONGODB_URL
		recoverEnvVariables()
	})

	Describe("a form gallery", func() {

		Context("with appropiate context", func() {
			var (
				gal1, gal2 *model.FormGallery
				errG       *web.AppError
			)

			JustBeforeEach(func() {

				// get the fixtures from appropiate json file for data form submission
				file, e := ioutil.ReadFile(dataFormGalleriesIds)
				if e != nil {
					log.Fatalf("opening config file %v", e.Error())
				}

				var gal model.FormGallery
				e = json.Unmarshal(file, &gal)
				if e != nil {
					log.Fatalf("Error reading forms galleries. %v", e.Error())
				}

				// new context
				c := web.NewContext(nil, nil)
				defer c.Close()

				// get the gallery to create
				c.SetValue("form_id", result.ID.Hex())
				c.Marshall(gal)

				// create gallery
				gal1, errG = service.CreateFormGallery(c)
				Expect(errG).Should(BeNil())

				// update fields
				gal1.Description = "New Description"

				// new context
				c2 := web.NewContext(nil, nil)
				defer c2.Close()

				// get the gallery to create
				c2.SetValue("gallery_id", gal1.ID.Hex())
				c2.Marshall(gal1)

				// update it
				gal2, errG = service.UpdateFormGallery(c2)
			})

			It("should return save the form gallery without errors", func() {
				Expect(errG).Should(BeNil())
			})
			It("should not create a new one", func() {
				Expect(gal2.Id()).Should(Equal(gal1.Id()))
			})
			It("should update the new fields", func() {
				Expect(gal2.Description).Should(Equal("New Description"))
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

			var result map[string]interface{}

			JustBeforeEach(func() {
				// create the context for this form
				c := web.NewContext(nil, nil)
				defer c.Close()
				c.SetValue("form_id", formid)
				result, err = service.GetFormSubmissionsByForm(c)

			})
			It("should return at least a submission to a form and no error", func() {
				Expect(err).Should(BeNil(), "is giving an error")
				Expect(len(result["submissions"].([]model.FormSubmission))).ShouldNot(Equal(0), "it should return at least one submission")
			})
		})
	})

	Describe("submissions to a form", func() {
		Context("when using limit", func() {

			var result map[string]interface{}

			JustBeforeEach(func() {
				// create the context for this form
				c := web.NewContext(nil, nil)
				defer c.Close()
				c.SetValue("form_id", formid)
				c.SetValue("limit", "3")
				result, err = service.GetFormSubmissionsByForm(c)
			})
			It("should not give back an error and return right amounts", func() {
				countsbyflag := result["counts"].(map[string]interface{})["search_by_flag"].(map[string]int)
				Expect(err).Should(BeNil())
				Expect(len(result["submissions"].([]model.FormSubmission))).Should(Equal(3), "on limit on submissions")
				Expect(countsbyflag["flagged"]).Should(Equal(5), "on flagged")
				Expect(countsbyflag["test_another_flag"]).Should(Equal(1), "on test_another_flag")
				Expect(countsbyflag["gophers"]).Should(Equal(3), "on gophers")
				Expect(countsbyflag["pythoners"]).Should(Equal(1), "on pythoners")
			})
		})
	})

	Describe("submissions to a form", func() {
		Context("order by date", func() {

			var result map[string]interface{}

			JustBeforeEach(func() {
				// create the context for this form
				c := web.NewContext(nil, nil)
				defer c.Close()
				c.SetValue("form_id", formid)
				c.SetValue("orderby", "asc")
				result, err = service.GetFormSubmissionsByForm(c)
			})
			It("should order by date", func() {
				fss := result["submissions"].([]model.FormSubmission)
				Expect(err).Should(BeNil())
				Expect(len(fss)).ShouldNot(Equal(0))

				// order by date in an ascending mode
				d := fss[1].DateCreated.Sub(fss[0].DateCreated).Seconds() >= 0
				Expect(d).Should(BeTrue(), fmt.Sprintf("%v is newer than %v", fss[0].DateCreated, fss[1].DateCreated))
			})
		})
	})

	Describe("submissions to a form", func() {
		Context("with the right context", func() {

			var result map[string]interface{}

			JustBeforeEach(func() {
				// create the context for this form
				c := web.NewContext(nil, nil)
				defer c.Close()
				c.SetValue("form_id", formid)
				c.SetValue("filterby", "flagged")
				result, err = service.GetFormSubmissionsByForm(c)
			})
			It("should bring back total of submissions to that form", func() {
				Expect(err).Should(BeNil())

				counts := result["counts"].(map[string]interface{})
				expectedCounts := 9 // total of all submissions for that form
				Expect(counts["total_submissions"]).Should(Equal(expectedCounts))
			})
		})
	})

	Describe("submissions to a form", func() {
		Context("filtering by not having a flag", func() {

			var result map[string]interface{}

			JustBeforeEach(func() {
				// create the context for this form
				c := web.NewContext(nil, nil)
				defer c.Close()
				c.SetValue("form_id", formid)
				c.SetValue("filterby", "-flagged")
				result, err = service.GetFormSubmissionsByForm(c)
			})
			It("should bring back total of submissions to that form", func() {
				Expect(err).Should(BeNil())

				counts := result["counts"].(map[string]interface{})
				expectedCounts := 9 // total of all submissions for that form
				Expect(counts["total_submissions"]).Should(Equal(expectedCounts))

				for _, s := range result["submissions"].([]model.FormSubmission) {
					Expect(s.Flags).ShouldNot(ContainElement("flagged"))
				}

				// total of submissions that do not have flagged as a flag
				expectedLen := 5 // there are 5 submissions with flag flagged
				Expect(len(result["submissions"].([]model.FormSubmission))).Should(Equal(expectedLen))

			})
		})
	})

	Describe("submissions to a form", func() {
		Context("filtering and with search parameters", func() {

			var result map[string]interface{}

			JustBeforeEach(func() {
				// create the context for this form
				c := web.NewContext(nil, nil)
				defer c.Close()
				c.SetValue("form_id", formid)
				c.SetValue("search", "Gophers")   //2 - 577c197810780b3401e7a1cf & 577c197810780b3401e7a3af
				c.SetValue("filterby", "flagged") //1 - 577c197810780b3401e7a3af
				result, err = service.GetFormSubmissionsByForm(c)
			})
			It("should bring no errors and all the right numbers", func() {
				Expect(err).Should(BeNil())

				counts := result["counts"].(map[string]interface{})

				expectedTotalSubmissions := 9 // total of all submissions for that form
				Expect(counts["total_submissions"]).Should(Equal(expectedTotalSubmissions), "total submissions for this form")

				expectedTotalSearch := 2
				// total of all submissions with the search applied, without filtering
				Expect(counts["total_search"]).Should(Equal(expectedTotalSearch), "total submissions for this specific search and filterby")

				expectedSearchByFlag := map[string]int{"flagged": 2, "something_else": 1}
				// total of all submissions with the search applied, without filtering, by flag
				Expect(counts["search_by_flag"].(map[string]int)["flagged"]).Should(Equal(expectedSearchByFlag["flagged"]), "total submissions by flag flagged")
				Expect(counts["search_by_flag"].(map[string]int)["something_else"]).Should(Equal(expectedSearchByFlag["something_else"]), "total submissions by flag something_else")
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
		objectS []model.FormSubmission
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
		c.Marshall(objectS[0])
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
				context.SetValue("flag", "flagged")
				sub, err = service.AddFlagToFormSubmission(context)

				Expect(err).Should(BeNil())
				Expect(len(sub.Flags)).Should(Equal(fCount + 1))
			})

			It("should not be able to add a flag twice", func() {
				context.SetValue("flag", "flagged")
				service.AddFlagToFormSubmission(context)
				sub, err = service.AddFlagToFormSubmission(context)
				Expect(err).Should(BeNil())
				Expect(len(sub.Flags)).Should(Equal(fCount + 1))
			})

			It("should be able to add a second flag", func() {
				context.SetValue("flag", "flagged")
				service.AddFlagToFormSubmission(context)
				context.SetValue("flag", "test_another__flag")
				_, err = service.AddFlagToFormSubmission(context)
				Expect(err).Should(BeNil())
			})

			It("should be able to remove a flag to a gallery after removing it", func() {
				context.SetValue("flag", "flagged")
				service.AddFlagToFormSubmission(context)
				_, err = service.RemoveFlagFromFormSubmission(context)
				Expect(err).Should(BeNil())
			})

			It("should get the right count after adding and removing", func() {
				context.SetValue("flag", "flagged")
				service.AddFlagToFormSubmission(context)
				sub, _ = service.RemoveFlagFromFormSubmission(context)
				Expect(len(sub.Flags)).ShouldNot(Equal(fCount + 1))
			})

			It("should not be able to remove a flag twice", func() {
				context.SetValue("flag", "flagged")
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
		objectS []model.FormSubmission
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
		c.Marshall(objectS[0])
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
		erra    *web.AppError
		objects []model.User
		assets  []model.Asset
		c       *web.AppContext
	)

	BeforeEach(func() {
		setTestDatabase()

		// get the fixtures from appropiate json file for data assets
		assets = getDataAssets(dataAssets)

		ca := web.NewContext(nil, nil)
		defer ca.Close()
		ca.Marshall(assets[0])

		_, erra = service.ImportAsset(ca)

		// get the fixtures from appropiate json file for data users
		objects = getDataUsers(dataUsers)

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

	Describe("an asset", func() {
		Context("with appropiate context", func() {
			It("should not return an Error", func() {
				Expect(erra).Should(BeNil())
			})
		})
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
