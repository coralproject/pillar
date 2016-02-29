package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

// CreateTagTargets creates TagTarget entries for various tags on an entity
func CreateTagTargets(manager *MongoManager, tags []string, tt *model.TagTarget) error {

	for _, name := range tags {

		tt.ID = bson.NewObjectId()
		tt.Name = name
		tt.DateCreated = time.Now()

		//skip the same entry, if exists
		dbEntity := model.TagTarget{}
		manager.TagTargets.Find(bson.M{"target_id": tt.TargetID, "name": name, "target": tt.Target}).One(&dbEntity)
		if dbEntity.ID != "" {
			continue
		}

		if err := manager.TagTargets.Insert(tt); err != nil {
			return err
		}
	}

	return nil
}

// CreateUpdateTag adds or updates a tag
func CreateUpdateTag(object *model.Tag) (*model.Tag, *AppError) {
	manager := GetMongoManager()
	defer manager.Close()

	//Create a new one, set created-date for the new ones
	if object.Old_Name == "" {
		return createTag(object, manager)
	}

	return updateTag(object, manager)
}

// creates a new Tag
func createTag(object *model.Tag, manager *MongoManager) (*model.Tag, *AppError) {
	var dbEntity model.Tag

	//return, if exists
	if manager.Tags.FindId(object.Name).One(&dbEntity); dbEntity.Name != "" {
		message := fmt.Sprintf("Tag exists [%+v]\n", object)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	object.DateCreated = time.Now()
	if err := manager.Tags.Insert(object); err != nil {
		message := fmt.Sprintf("Error creating tag [%+v]", object)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	return object, nil
}

// updates an existing Tag
func updateTag(object *model.Tag, manager *MongoManager) (*model.Tag, *AppError) {

	var dbEntity model.Tag
	if manager.Tags.FindId(object.Old_Name).One(&dbEntity); dbEntity.Name == "" {
		message := fmt.Sprintf("Cannot update, tag not found: [%+v]", object)
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
	if err := manager.Tags.RemoveId(object.Old_Name); err != nil {
		message := fmt.Sprintf("Error removing old tag [%+v]", object)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	//insert new tag
	if err := manager.Tags.Insert(newTag); err != nil {
		message := fmt.Sprintf("Error creating tag [%+v]", newTag)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	var user model.User
	iter := manager.Users.Find(bson.M{"tags": object.Old_Name}).Iter()
	for iter.Next(&user) {
		tags := make([]string, len(user.Tags))
		for _, one := range user.Tags {
			if one != object.Old_Name {
				tags = append(tags, one)
			}
		}
		tags = append(tags, object.Name)
		manager.Users.Update(
			bson.M{"_id": user.ID},
			bson.M{"$set": bson.M{"tags": tags}},
		)
	}

	return &newTag, nil
}



// GetTags returns an array of tags
func GetTags() ([]model.Tag, *AppError) {
	manager := GetMongoManager()
	defer manager.Close()

	//set created-date for the new ones
	all := make([]model.Tag, 0)
	if err := manager.Tags.Find(nil).All(&all); err != nil {
		message := fmt.Sprintf("Error fetching tags")
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	return all, nil
}

// DeleteTag deletes a tag
func DeleteTag(object *model.Tag) *AppError {
	manager := GetMongoManager()
	defer manager.Close()

	//we must have the tag name for deletion
	if object.Name == "" {
		message := fmt.Sprintf("Cannot delete an invalid tag [%v]", object)
		return &AppError{nil, message, http.StatusInternalServerError}
	}

	//delete
	if err := manager.Tags.RemoveId(object.Name); err != nil {
		message := fmt.Sprintf("Error deleting tag [%v]", object)
		return &AppError{err, message, http.StatusInternalServerError}
	}

	return nil
}
