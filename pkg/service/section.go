package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"net/http"
	"time"
)

// CreateSection creates Section
func CreateSection(context *AppContext) (*model.Section, *AppError) {
	var input model.Section
	context.Unmarshall(&input)

	if input.Name == "" {
		message := fmt.Sprintf("Invalid Section Name [%s]", input.Name)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	var dbEntity model.Section
	if context.DB.Sections.FindId(input.Name).One(&dbEntity); dbEntity.Name == "" {
		input.DateCreated = time.Now()
	}

	input.DateUpdated = time.Now()
	if _, err := context.DB.Sections.UpsertId(input.Name, input); err != nil {
		message := fmt.Sprintf("Error creating/updating Section")
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	return &input, nil
}

// GetSections returns an array of Sections
func GetSections(context *AppContext) ([]model.Section, *AppError) {

	all := make([]model.Section, 0)
	if err := context.DB.Sections.Find(nil).All(&all); err != nil {
		message := fmt.Sprintf("Error fetching Sections")
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	return all, nil
}


