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
	Route{"GET", "/about", handler.AboutThisApp},

	//Import Handlers
	Route{"POST", "/api/import/asset", handler.ImportAsset},
	Route{"POST", "/api/import/user", handler.ImportUser},
	Route{"POST", "/api/import/comment", handler.ImportComment},
	Route{"POST", "/api/import/action", handler.ImportAction},
	Route{"POST", "/api/import/note", handler.ImportNote},

	//Tag Handlers
	Route{"GET", "/api/tags", handler.GetTags},
	Route{"POST", "/api/tag", handler.CreateUpdateTag},
	Route{"DELETE", "/api/tag", handler.DeleteTag},

	//Search Handlers
	Route{"GET", "/api/searches", handler.GetSearches},
	Route{"GET", "/api/search/{id}", handler.GetSearch},
	Route{"PUT", "/api/search", handler.CreateUpdateSearch},
	Route{"POST", "/api/search", handler.CreateUpdateSearch},
	Route{"DELETE", "/api/search/{id}", handler.DeleteSearch},

	//Manage User Activities
	Route{"POST", "/api/cay/useraction", handler.CreateUserAction},

	//Create/Update Handlers
	Route{"POST", "/api/author", handler.CreateUpdateAuthor},
	Route{"POST", "/api/section", handler.CreateUpdateSection},
	Route{"POST", "/api/asset", handler.CreateUpdateAsset},
	Route{"POST", "/api/user", handler.CreateUpdateUser},
	Route{"POST", "/api/comment", handler.CreateUpdateComment},
	Route{"POST", "/api/index", handler.CreateIndex},
	Route{"POST", "/api/metadata", handler.UpdateMetadata},

	// Forms
	Route{"POST", "/api/form", handler.CreateUpdateForm},
	Route{"PUT", "/api/form", handler.CreateUpdateForm},
	Route{"PUT", "/api/form/{id}/status/{status}", handler.UpdateFormStatus},
	Route{"GET", "/api/forms", handler.GetForms},
	Route{"GET", "/api/form/{id}", handler.GetForm},
	Route{"DELETE", "/api/form/{id}", handler.DeleteForm},

	// Form Submissions
	Route{"POST", "/api/form_submission/{form_id}", handler.CreateFormSubmission},
	Route{"PUT", "/api/form_submission/{id}/status/{status}", handler.UpdateFormSubmissionStatus},
	Route{"GET", "/api/form_submissions/{form_id}", handler.GetFormSubmissionsByForm},
	Route{"GET", "/api/form_submission/{id}", handler.GetFormSubmission},
	Route{"PUT", "/api/form_submission/{id}/{answer_id}", handler.EditFormSubmissionAnswer},
	Route{"PUT", "/api/form_submission/{id}/flag/{flag}", handler.AddFlagToFormSubmission},
	Route{"DELETE", "/api/form_submission/{id}/flag/{flag}", handler.RemoveFlagFromFormSubmission},
	Route{"DELETE", "/api/form_submission/{id}", handler.DeleteFormSubmission},

	// Form Galleries
	Route{"PUT", "/api/form_gallery/{id}/add/{submission_id}/{answer_id}", handler.AddAnswerToFormGallery},
	Route{"DELETE", "/api/form_gallery/{id}/remove/{submission_id}/{answer_id}", handler.RemoveAnswerFromFormGallery},
}
