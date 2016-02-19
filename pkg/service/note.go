package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/model"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// CreateNote creates a new note resource
func CreateNote(object *model.Note) (*model.Note, *AppError) {

	// Insert Comment
	manager := GetMongoManager()
	defer manager.Close()

	if object.UserID == "" {
		//find user using source information and set the reference
		var user model.User
		manager.Users.Find(bson.M{"source.id": object.Source.UserID}).One(&user)
		if user.ID == "" {
			message := fmt.Sprintf("Invalid user with source ID [%s]\n", object.Source.UserID)
			return nil, &AppError{nil, message, http.StatusInternalServerError}
		}
		object.UserID = user.ID
	}

	//find target and set the reference
	switch object.Target {
	case model.Users:
		addNoteToUser(object, manager)
		break

	case model.Comments:
		addNoteToComment(object, manager)
		break
	}

	return object, nil
}

func addNoteToComment(object *model.Note, manager *MongoManager) (*model.Note, *AppError) {
	var dbEntity model.Comment

	if object.TargetID != "" {
		if manager.Comments.FindId(object.TargetID).One(&dbEntity); dbEntity.ID == "" {
			message := fmt.Sprintf("Invalid comment ID [%s]\n", object.TargetID)
			return nil, &AppError{nil, message, http.StatusInternalServerError}
		}
	} else {
		if manager.Comments.Find(bson.M{"source.id": object.Source.TargetID}).One(&dbEntity); dbEntity.ID == "" {
			message := fmt.Sprintf("Invalid comment source ID [%s]\n", object.Source.TargetID)
			return nil, &AppError{nil, message, http.StatusInternalServerError}
		}
	}

	//append this note to comment's notes array
	notes := append(dbEntity.Notes, *object)
	manager.Comments.Update(bson.M{"_id": dbEntity.ID},
		bson.M{"$set": bson.M{"notes": notes}})

	return object, nil
}

func addNoteToUser(object *model.Note, manager *MongoManager) (*model.Note, *AppError) {
	var dbEntity model.User

	if object.TargetID != "" {
		if manager.Users.FindId(object.TargetID).One(&dbEntity); dbEntity.ID == "" {
			message := fmt.Sprintf("Invalid user ID [%s]\n", object.TargetID)
			return nil, &AppError{nil, message, http.StatusInternalServerError}
		}
	} else {
		if manager.Users.Find(bson.M{"source.id": object.Source.TargetID}).One(&dbEntity); dbEntity.ID == "" {
			message := fmt.Sprintf("Invalid user source ID [%s]\n", object.Source.TargetID)
			return nil, &AppError{nil, message, http.StatusInternalServerError}
		}
	}

	//append this note to comment's notes array
	notes := append(dbEntity.Notes, *object)
	manager.Comments.Update(bson.M{"_id": dbEntity.ID},
		bson.M{"$set": bson.M{"notes": notes}})

	return object, nil
}
