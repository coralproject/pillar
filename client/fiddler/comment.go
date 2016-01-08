package fiddler

import (
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/client/db"
	"github.com/coralproject/pillar/client/rest"
	"github.com/coralproject/pillar/server/model"
	"github.com/stretchr/stew/objects"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
	"time"
)

//LoadComments imports comments
func LoadComments() {

	//story := "http://washingtonpost.com/posteverything/wp/2015/05/20/feminists-want-us-to-define-these-ugly-sexual-encounters-as-rape-dont-let-them/"
	story := "http://washingtonpost.com/opinions/reformers-want-to-erase-confuciuss-influence-in-asia-thats-a-mistake/2015/05/28/529c1d3a-042e-11e5-a428-c984eb077d4e_story.html"
	//story := "http://washingtonpost.com/world/europe/european-leaders-seek-last-ditch-offer-to-bring-greece-from-brink-of-default/2015/06/30/960aded8-1ea2-11e5-a135-935065bc30d0_story.html"

	manager := db.GetMongoManager()
	defer manager.Close()

	all := make([]interface{}, 10)

	manager.Comments.Find(bson.M{"object.context.uri": story}).Sort("postedTime").All(&all)

	fmt.Printf("Found %d comments\n", len(all))
	fmt.Printf("Import in progress...\n")
	var nActions, nAssets, nUsers, nSuccess, nFailure int
	for _, one := range all {
		data, _ := json.Marshal(one)

		commentJson := map[string]interface{}{}
		json.Unmarshal(data, &commentJson)

		asset := getAsset(commentJson)
		if response := rest.Request(rest.MethodPost, rest.URLAsset, getBuffer(asset)); response.StatusCode == 200 {
			nAssets++
		}

		if response := rest.Request(rest.MethodPost, rest.URLUser, getBuffer(getUser(commentJson))); response.StatusCode == 200 {
			nUsers++
		}

		users := getAllUsers(commentJson)
		for i := 0; i < len(users); i++ {
			if response := rest.Request(rest.MethodPost, rest.URLUser, getBuffer(users[i])); response.StatusCode == 200 {
				nUsers++
			}
		}

		comment := getComment(commentJson, asset.URL)
		if response := rest.Request(rest.MethodPost, rest.URLComment, getBuffer(comment)); response.StatusCode == 200 {
			nSuccess++
		} else {
			nFailure++
		}

		actions := getAllActions(commentJson, &comment)
		for i := 0; i < len(actions); i++ {
			action := actions[i]
			if response := rest.Request(rest.MethodPost, rest.URLAction, getBuffer(action)); response.StatusCode == 200 {
				nActions++
			}
		}
	}
	fmt.Printf("Finished importing comments, suceess: [%d] failure: [%d]\n\n\n", nSuccess, nFailure)
}

func getAsset(m objects.Map) model.Asset {
	asset := model.Asset{}
	url := getArray(m.Get("object.context"))[0].GetString("uri")
	asset.URL = url
	asset.SourceID = url
	return asset
}

func getUser(m objects.Map) model.User {
	user := model.User{}

	user.SourceID = m.GetStringOrEmpty("actor.id")
	user.UserName = m.GetStringOrEmpty("actor.title")
	user.Status = m.GetStringOrEmpty("actor.status")
	user.Avatar = m.GetStringOrEmpty("actor.avatar")

	return user
}

func getAllUsers(m objects.Map) []model.User {
	users := []model.User{}

	maps := getArray(m.Get("object.likes"))
	for _, one := range maps {
		m := one.Get("actor")
		if m != nil {
			users = append(users, getUser(m.(map[string]interface{})))
		}
	}

	maps = getArray(m.Get("object.flags"))
	for _, one := range maps {
		m := one.Get("actor")
		if m != nil {
			users = append(users, getUser(m.(map[string]interface{})))
		}
	}

	return users
}

func getComment(m objects.Map, url string) model.Comment {
	comment := model.Comment{}

	comment.Body = m.GetString("object.content")
	comment.Status = m.GetString("object.status")

	t, _ := time.Parse(time.RFC3339, m.GetString("postedTime"))
	comment.DateCreated = t

	t, _ = time.Parse(time.RFC3339, m.GetString("updated"))
	comment.DateUpdated = t

	comment.Source.ID = getShortCommentID(m.GetString("id"))
	comment.Source.AssetID = url
	comment.Source.UserID = m.GetString("actor.id")

	target := getArray(m.Get("targets"))[0]
	targetID := getShortCommentID(target.GetString("id"))
	comment.Source.ParentID = getShortCommentID(targetID)

	//fmt.Printf("Comment: %s\n\n", comment.Source.ID)

//	//get likes and flags as actions
//	populateActions(m.Get("object.likes"), model.ActionTypeLikes, &comment)
//	//	fmt.Printf("Getting Flags....\n\n")
//	populateActions(m.Get("object.flags"), model.ActionTypeFlags, &comment)

	//fmt.Printf("Actions....%d\n\n", len(comment.Actions))

	return comment
}

//returns the first one from an array of json documents
func getOne(list interface{}) objects.Map {

	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice:
		slice := reflect.ValueOf(list)

		//must convert the Interface to map[string]interface{}
		//so that it can be converted to an objects.Map
		//var m map[string]interface{}
		m := slice.Index(0).Interface().(map[string]interface{})

		return objects.Map(m)
	}
	return nil
}

//str := "http://washingtonpost.com/ECHO/item/c5c8f176-3f27-4228-94d1-8dffc73028ac"
//str = "http://js-kit.com/activities/post/c5c8f176-3f27-4228-94d1-8dffc73028ac"
func getShortCommentID(url string) string {
	var s []string
	s = strings.Split(url, "/")
	return s[(len(s) - 1)]
}

func getAllActions(m objects.Map, comment *model.Comment) []model.Action {

	actions := []model.Action{}

	//Add all likes
	mapArray := getArray(m.Get("object.likes"))
	for i := 0; i < len(mapArray); i++ {
		if mapArray[i] == nil {
			continue
		}
		actions = append(actions, getOneAction(mapArray[i], comment, model.ActionTypeLikes))
	}

	//Add all flags
	mapArray = getArray(m.Get("object.flags"))
	for i := 0; i < len(mapArray); i++ {
		if mapArray[i] == nil {
			continue
		}
		actions = append(actions, getOneAction(mapArray[i], comment, model.ActionTypeFlags))
	}

	return actions
}

func getOneAction(m objects.Map, comment *model.Comment, actionType string) model.Action {
	action := model.Action{}

	t, _ := time.Parse(time.RFC3339, m.GetString("published"))
	action.Date = t
	action.Type = actionType
	action.TargetType = model.TargetTypeComment
	action.Source.UserID = m.GetString("actor.id")
	action.Source.TargetID = comment.Source.ID

	return action
}

//func populateActions(list interface{}, actionType string, comment *model.Comment) {
//
//	array := getArray(list)
//	if len(array) == 0 {
//		return
//	}
//
//	for i := 0; i < len(array); i++ {
//		action := model.Action{}
//		if array[i] == nil {
//			continue
//		}
//		//fmt.Printf("Item: %s\n\n\n", array[i])
//
//		t, _ := time.Parse(time.RFC3339, array[i].GetString("published"))
//		//		time.Parse(time.RFC3339, m.GetString("updated"))
//		//		t, _ := time.Parse(shortForm, array[i].GetString("published"))
//		action.SourceUserID = array[i].GetString("actor.id")
//		action.Date = t
//		action.Type = actionType
//		//fmt.Printf("Action: %+v\n", action)
//		comment.Actions = append(comment.Actions, action)
//	}
//}

//when the item is an array, we must convert it to a slice
func getArray(list interface{}) []objects.Map {

	var resultArray []objects.Map
	if list == nil {
		return resultArray
	}

	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice:
		slice := reflect.ValueOf(list)

		//must convert the Interface to map[string]interface{}
		//so that it can be converted to an objects.Map
		//fmt.Printf("Size of slice: %d\n\n", slice.Len())
		for i := 0; i < slice.Len(); i++ {
			//var m map[string]interface{}
			//fmt.Printf("Item: %s\n\n", slice.Index(i))
			resultArray = append(resultArray, slice.Index(i).Interface().(map[string]interface{}))
		}

		return resultArray
	}

	return nil
}
