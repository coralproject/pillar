package crud

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"reflect"
)

// CreateUser creates a new user resource
func CreateUser(object *User) (*User, *AppError) {

	// get a mongo connection
	manager := GetMongoManager()
	defer manager.Close()

	dbEntity := User{}

	//return, if exists
	manager.Users.FindId(object.ID).One(&dbEntity)
	if dbEntity.ID != "" {
		message := fmt.Sprintf("%s exists with ID [%s]\n", reflect.TypeOf(object).Name(), object.ID)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	//upsert if entity exists with same source.id
	manager.Users.Find(bson.M{"source.id": object.Source.ID}).One(&dbEntity)
	if dbEntity.ID != "" {
		object.ID = dbEntity.ID
		_, err := manager.Users.UpsertId(dbEntity.ID, object)
		if err != nil {
			message := fmt.Sprintf("Error updating existing User [%s], %s", object.Source.ID, err)
			return nil, &AppError{err, message, http.StatusInternalServerError}
		}
		return object, nil
	}

	object.ID = bson.NewObjectId()
	err := manager.Users.Insert(object)
	if err != nil {
		message := fmt.Sprintf("Error creating user [%s]", err)
		return nil, &AppError{err, message, http.StatusInternalServerError}
	}

	err = CreateTagTargets(manager, object.Tags, &TagTarget{Target:Users, TargetID:object.ID})
	if err != nil {
		message := fmt.Sprintf("Error creating TagStat [%s]", err)
		return nil, &AppError{nil, message, http.StatusInternalServerError}
	}

	return object, nil
}

//append action to user's actions array and update stats
func updateUserOnAction(user *User, object *Action, manager *MongoManager) {
	actions := append(user.Actions, object.ID)
	if user.Stats[object.Type] == nil {
		user.Stats[object.Type] = 0
	}

	user.Stats[object.Type] = user.Stats[object.Type].(int) + 1
	manager.Comments.Update(
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"actions": actions, "stats": user.Stats}},
	)
}

//update stats on this user for #comments
func updateUserOnComment(user *User, manager *MongoManager) {
	if user.Stats == nil {
		user.Stats = make(map[string]interface{})
	}

	if user.Stats[StatsComments] == nil {
		user.Stats[StatsComments] = 0
	}

	user.Stats[StatsComments] = user.Stats[StatsComments].(int) + 1
	manager.Users.Update(
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{"stats": user.Stats}},
	)
}
