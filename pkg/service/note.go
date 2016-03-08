package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// CreateNote creates a new note resource
func CreateNote(context *web.AppContext) (*model.Note, *web.AppError) {

	var input model.Note
	context.Unmarshall(&input)

	// Insert Comment
	if input.UserID == "" {
		//find user using source information and set the reference
		var user model.User
		context.DB.Users.Find(bson.M{"source.id": input.Source.UserID}).One(&user)
		if user.ID == "" {
			message := fmt.Sprintf("Invalid user with source ID [%s]\n", input.Source.UserID)
			return nil, &web.AppError{nil, message, http.StatusInternalServerError}
		}
		input.UserID = user.ID
	}

	//find target and set the reference
	switch input.Target {
	case model.Users:
		addNoteToUser(context.DB, &input)
		break

	case model.Comments:
		addNoteToComment(context.DB, &input)
		break
	}

	return &input, nil
}

func addNoteToComment(db *db.MongoDB, object *model.Note) (*model.Note, *web.AppError) {
	var dbEntity model.Comment

	if object.TargetID != "" {
		if db.Comments.FindId(object.TargetID).One(&dbEntity); dbEntity.ID == "" {
			message := fmt.Sprintf("Invalid comment ID [%s]\n", object.TargetID)
			return nil, &web.AppError{nil, message, http.StatusInternalServerError}
		}
	} else {
		if db.Comments.Find(bson.M{"source.id": object.Source.TargetID}).One(&dbEntity); dbEntity.ID == "" {
			message := fmt.Sprintf("Invalid comment source ID [%s]\n", object.Source.TargetID)
			return nil, &web.AppError{nil, message, http.StatusInternalServerError}
		}
	}

	//append this note to comment's notes array
	notes := append(dbEntity.Notes, *object)
	db.Comments.Update(bson.M{"_id": dbEntity.ID},
		bson.M{"$set": bson.M{"notes": notes}})

	return object, nil
}

func addNoteToUser(db *db.MongoDB, object *model.Note) (*model.Note, *web.AppError) {
	var dbEntity model.User

	if object.TargetID != "" {
		if db.Users.FindId(object.TargetID).One(&dbEntity); dbEntity.ID == "" {
			message := fmt.Sprintf("Invalid user ID [%s]\n", object.TargetID)
			return nil, &web.AppError{nil, message, http.StatusInternalServerError}
		}
	} else {
		if db.Users.Find(bson.M{"source.id": object.Source.TargetID}).One(&dbEntity); dbEntity.ID == "" {
			message := fmt.Sprintf("Invalid user source ID [%s]\n", object.Source.TargetID)
			return nil, &web.AppError{nil, message, http.StatusInternalServerError}
		}
	}

	//append this note to comment's notes array
	notes := append(dbEntity.Notes, *object)
	db.Comments.Update(bson.M{"_id": dbEntity.ID},
		bson.M{"$set": bson.M{"notes": notes}})

	return object, nil
}
