package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
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
func CreateUpdateTag(context *web.AppContext) (*model.Tag, *web.AppError) {
	var input model.Tag
	context.Unmarshall(&input)

	//old-name is empty, upsert one
	if input.OldName == "" {
		return upsertTag(context.DB, &input)
	}

	//since old-name is passed, this implies a rename
	return renameTag(context.DB, &input)
}

// creates a new Tag
func upsertTag(db *db.MongoDB, object *model.Tag) (*model.Tag, *web.AppError) {
	if object.Name == "" {
		message := fmt.Sprintf("Invalid Tag Name [%s]", object.Name)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	var dbEntity model.Tag
	if db.Tags.FindId(object.Name).One(&dbEntity); dbEntity.Name == "" {
		object.DateCreated = time.Now()
	}

	object.DateUpdated = time.Now()
	if _, err := db.Tags.UpsertId(object.Name, object); err != nil {
		message := fmt.Sprintf("Error creating tag [%+v]", object)
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return object, nil
}

// updates an existing Tag
func renameTag(db *db.MongoDB, object *model.Tag) (*model.Tag, *web.AppError) {

	var dbEntity model.Tag
	if db.Tags.FindId(object.OldName).One(&dbEntity); dbEntity.Name == "" {
		message := fmt.Sprintf("Cannot update, tag not found: [%s]", object.OldName)
		return nil, &web.AppError{nil, message, http.StatusInternalServerError}
	}

	var newTag model.Tag
	newTag.Name = object.Name
	newTag.OldName = object.OldName
	newTag.Description = dbEntity.Description
	if object.Description != "" {
		newTag.Description = object.Description
	}
	newTag.DateCreated = dbEntity.DateCreated
	newTag.DateUpdated = time.Now()

	//remove the old one
	if err := db.Tags.RemoveId(object.OldName); err != nil {
		message := fmt.Sprintf("Error removing old tag [%s]", object.OldName)
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	//insert new tag
	if err := db.Tags.Insert(newTag); err != nil {
		message := fmt.Sprintf("Error creating tag [%+v]", newTag)
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	var user model.User
	iter := db.Users.Find(bson.M{"tags": object.OldName}).Iter()
	for iter.Next(&user) {
		tags := make([]string, len(user.Tags))
		for _, one := range user.Tags {
			if one != object.OldName {
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
func GetTags(context *web.AppContext) ([]model.Tag, *web.AppError) {

	//set created-date for the new ones
	all := make([]model.Tag, 0)
	if err := context.DB.Tags.Find(nil).All(&all); err != nil {
		message := fmt.Sprintf("Error fetching tags")
		return nil, &web.AppError{err, message, http.StatusInternalServerError}
	}

	return all, nil
}

// DeleteTag deletes a tag
func DeleteTag(context *web.AppContext) *web.AppError {
	var input model.Tag
	context.Unmarshall(&input)

	//we must have the tag name for deletion
	if input.Name == "" {
		message := fmt.Sprintf("Cannot delete an invalid tag [%v]", input)
		return &web.AppError{nil, message, http.StatusInternalServerError}
	}

	//delete
	if err := context.DB.Tags.RemoveId(input.Name); err != nil {
		message := fmt.Sprintf("Error deleting tag [%v]", input)
		return &web.AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
