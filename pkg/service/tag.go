package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
	"github.com/coralproject/pillar/pkg/db"
)

// CreateTagTargets creates TagTarget entries for various tags on an entity
func CreateTagTargets(db *db.MongoDB, tags []string, tt *model.TagTarget) error {

	for _, name := range tags {

		tt.ID = bson.NewObjectId()
		tt.Name = name
		tt.DateCreated = time.Now()

		//skip the same entry, if exists
		dbEntity := model.TagTarget{}
		db.TagTargets.Find(bson.M{"target_id": tt.TargetID, "name": name, "target": tt.Target}).One(&dbEntity)
		if dbEntity.ID != "" {
			continue
		}

		if err := db.TagTargets.Insert(tt); err != nil {
			return err
		}
	}

	return nil
}

// CreateUpdateTag adds or updates a tag
func CreateUpdateTag(context *AppContext) (*model.Tag, *AppError) {
	db := context.DB
	input := context.Input.(model.Tag)

	//old-name is empty, upsert one
	if input.Old_Name == "" {
		return upsertTag(db, &input)
	}

	//since old-name is passed, this implies a rename
	return renameTag(db, &input)
}

// creates a new Tag
func upsertTag(db *db.MongoDB, object *model.Tag) (*model.Tag, *AppError) {
	if object.Name == "" {
		message := fmt.Sprintf("Invalid Tag Name [%s]", object.Name)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	var dbEntity model.Tag
	if db.Tags.FindId(object.Name).One(&dbEntity); dbEntity.Name == "" {
		object.DateCreated = time.Now()
	}

	object.DateUpdated = time.Now()
	if _, err := db.Tags.UpsertId(object.Name, object); err != nil {
		message := fmt.Sprintf("Error creating tag [%+v]", object)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	return object, nil
}

// updates an existing Tag
func renameTag(db *db.MongoDB, object *model.Tag) (*model.Tag, *AppError) {

	var dbEntity model.Tag
	if db.Tags.FindId(object.Old_Name).One(&dbEntity); dbEntity.Name == "" {
		message := fmt.Sprintf("Cannot update, tag not found: [%s]", object.Old_Name)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	var newTag model.Tag
	newTag.Name = object.Name
	newTag.Old_Name = object.Old_Name
	newTag.Description = dbEntity.Description
	if object.Description != "" {
		newTag.Description = object.Description
	}
	newTag.DateCreated = dbEntity.DateCreated
	newTag.DateUpdated = time.Now()

	//remove the old one
	if err := db.Tags.RemoveId(object.Old_Name); err != nil {
		message := fmt.Sprintf("Error removing old tag [%s]", object.Old_Name)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	//insert new tag
	if err := db.Tags.Insert(newTag); err != nil {
		message := fmt.Sprintf("Error creating tag [%+v]", newTag)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	var user model.User
	iter := db.Users.Find(bson.M{"tags": object.Old_Name}).Iter()
	for iter.Next(&user) {
		tags := make([]string, len(user.Tags))
		for _, one := range user.Tags {
			if one != object.Old_Name {
				tags = append(tags, one)
			}
		}
		tags = append(tags, object.Name)
		db.Users.Update(
			bson.M{"_id": user.ID},
			bson.M{"$set": bson.M{"tags": tags}},
		)
	}

	return &newTag, nil
}

// GetTags returns an array of tags
func GetTags(context *AppContext) ([]model.Tag, *AppError) {

	//set created-date for the new ones
	all := make([]model.Tag, 0)
	if err := context.DB.Tags.Find(nil).All(&all); err != nil {
		message := fmt.Sprintf("Error fetching tags")
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	return all, nil
}

// DeleteTag deletes a tag
func DeleteTag(context *AppContext) *AppError {
	db := context.DB
	object := context.Input.(model.Tag)

	//we must have the tag name for deletion
	if object.Name == "" {
		message := fmt.Sprintf("Cannot delete an invalid tag [%v]", object)
		return &AppError{nil, message, http.StatusInternalServerError}
	}

	//delete
	if err := db.Tags.RemoveId(object.Name); err != nil {
		message := fmt.Sprintf("Error deleting tag [%v]", object)
		return &AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
