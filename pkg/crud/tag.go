package crud

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
	"net/http"
)

// CreateTagTargets creates TagTarget entries for various tags on an entity
func CreateTagTargets(manager *MongoManager, tags []string, tt *TagTarget) error {

	for _, name := range tags {

		tt.ID = bson.NewObjectId()
		tt.Name = name
		tt.DateCreated = time.Now()

		//skip the same entry, if exists
		dbEntity := TagTarget{}
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

// UpsertTag adds/updates tags to the master list
func UpsertTag(object *Tag) (*Tag, *AppError) {
	manager := GetMongoManager()
	defer manager.Close()

	//set created-date for the new ones
	var dbEntity Tag
	if manager.Tags.FindId(object.Name).One(&dbEntity); dbEntity.Name == "" {
		object.DateCreated = time.Now()
	}

	object.DateUpdated = time.Now()
	_, err := manager.Tags.UpsertId(object.Name, object)
	if err != nil {
		message := fmt.Sprintf("Error creating tag [%+v]", object)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}
	fmt.Printf("Tag: %+v\n\n", object)

	return object, nil
}

// GetTags returns an array of tags
func GetTags() ([]Tag, *AppError) {
	manager := GetMongoManager()
	defer manager.Close()

	//set created-date for the new ones
	var all []Tag
	if err := manager.Tags.Find(nil).All(&all); err != nil {
		message := fmt.Sprintf("Error fetching tags")
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	return all, nil
}

// DeleteTag deletes a tag
func DeleteTag(object *Tag) *AppError {
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

