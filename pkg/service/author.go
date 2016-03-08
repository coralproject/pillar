package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"net/http"
)

// CreateUpdateAuthor creates/updates an author resource
func CreateUpdateAuthor(context *AppContext) (*model.Author, *AppError) {
	var input model.Author
	context.Unmarshall(&input)

	if input.ID == "" {
		message := fmt.Sprintf("Invalid Author Id [%s]", input.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	db := context.DB
	if _, err := db.Authors.UpsertId(input.ID, input); err != nil {
		message := fmt.Sprintf("Error creating/updating Author")
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	return &input, nil
}
