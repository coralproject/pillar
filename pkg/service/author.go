package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
	"net/http"
)

// CreateUpdateAuthor creates/updates an author resource
func CreateUpdateAuthor(context *web.AppContext) (*model.Author, *web.AppError) {
	var input model.Author
	if err := UnmarshallAndValidate(context, &input); err != nil {
		return nil, err
	}

	if input.ID == "" {
		message := fmt.Sprintf("Invalid Author Id [%s]", input.ID)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	db := context.DB
	if _, err := db.Authors.UpsertId(input.ID, input); err != nil {
		message := fmt.Sprintf("Error creating/updating Author")
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return &input, nil
}
