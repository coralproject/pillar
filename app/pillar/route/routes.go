package route

import (
	"github.com/coralproject/pillar/app/pillar/handler"
	"github.com/coralproject/pillar/pkg/web"
)

//Route defines mappings of end-points to handler methods
type Route struct {
	Method      string
	Pattern     string
	HandlerFunc web.HandlerFunc
}

var routes = []Route{
	//Generic or Common ones
	{"GET", "/about", handler.AboutThisApp},

	//Import Handlers
	{"POST", "/api/import/asset", handler.ImportAsset},
	{"POST", "/api/import/user", handler.ImportUser},
	{"POST", "/api/import/comment", handler.ImportComment},
	{"POST", "/api/import/action", handler.ImportAction},
	{"POST", "/api/import/note", handler.ImportNote},

	//Tag Handlers
	{"GET", "/api/tags", handler.GetTags},
	{"POST", "/api/tag", handler.CreateUpdateTag},
	{"DELETE", "/api/tag", handler.DeleteTag},

	//Search Handlers
	{"GET", "/api/searches", handler.GetSearches},
	{"GET", "/api/search/{id}", handler.GetSearch},
	{"PUT", "/api/search", handler.CreateUpdateSearch},
	{"POST", "/api/search", handler.CreateUpdateSearch},
	{"DELETE", "/api/search/{id}", handler.DeleteSearch},

	//Manage User Activities
	{"POST", "/api/cay/useraction", handler.CreateUserAction},

	//Create/Update Handlers
	{"POST", "/api/author", handler.CreateUpdateAuthor},
	{"POST", "/api/section", handler.CreateUpdateSection},
	{"POST", "/api/asset", handler.CreateUpdateAsset},
	{"POST", "/api/user", handler.CreateUpdateUser},
	{"POST", "/api/comment", handler.CreateUpdateComment},
	{"POST", "/api/index", handler.CreateIndex},
	{"POST", "/api/metadata", handler.UpdateMetadata},

	// Forms
	{"POST", "/api/form", handler.CreateUpdateForm},
	{"PUT", "/api/form", handler.CreateUpdateForm},
	{"PUT", "/api/form/{id}/status/{status}", handler.UpdateFormStatus},
	{"GET", "/api/forms", handler.GetForms},
	{"GET", "/api/form/{id}", handler.GetForm},
	{"DELETE", "/api/form/{id}", handler.DeleteForm},

	// Form Submissions
	{"POST", "/api/form_submission/{form_id}", handler.CreateFormSubmission},
	{"PUT", "/api/form_submission/{id}/status/{status}", handler.UpdateFormSubmissionStatus},
	{"GET", "/api/form_submissions/{form_id}", handler.GetFormSubmissionsByForm},
	{"GET", "/api/form_submission/{id}", handler.GetFormSubmission},
	{"GET", "/api/form_submissions/search/{search}", handler.SearchFormSubmissions},
	{"PUT", "/api/form_submission/{id}/{answer_id}", handler.EditFormSubmissionAnswer},
	{"PUT", "/api/form_submission/{id}/flag/{flag}", handler.AddFlagToFormSubmission},
	{"DELETE", "/api/form_submission/{id}/flag/{flag}", handler.RemoveFlagFromFormSubmission},
	{"DELETE", "/api/form_submission/{id}", handler.DeleteFormSubmission},

	// Form Galleries
	{"GET", "/api/form_gallery/{id}", handler.GetFormGallery},
	{"GET", "/api/form_galleries/{form_id}", handler.GetFormGalleriesByForm},
	{"GET", "/api/form_galleries/form/{form_id}", handler.GetFormGalleriesByForm}, // a more explicit version of the above for clarity
	{"PUT", "/api/form_gallery/{id}/add/{submission_id}/{answer_id}", handler.AddAnswerToFormGallery},
	{"DELETE", "/api/form_gallery/{id}/remove/{submission_id}/{answer_id}", handler.RemoveAnswerFromFormGallery},
}
