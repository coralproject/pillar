package main

import (
	"bytes"
	"encoding/json"
	"github.com/coralproject/pillar/server/model"
	"github.com/stretchr/stew/objects"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
	"time"
)

func wapoFiddler() {

	story := "http://washingtonpost.com/opinions/reformers-want-to-erase-confuciuss-influence-in-asia-thats-a-mistake/2015/05/28/529c1d3a-042e-11e5-a428-c984eb077d4e_story.html"

	manager := getMongoManager()
	defer manager.close()

	all := make([]interface{}, 10)

	manager.comments.Find(bson.M{"object.context.uri": story}).Sort("postedTime").All(&all)

	for _, one := range all {
		data, _ := json.Marshal(one)

		comment := map[string]interface{}{}
		json.Unmarshal(data, &comment)

		asset := getAsset(comment)
		doRequest(methodPost, urlAsset, getBuffer(asset))
		doRequest(methodPost, urlUser, getBuffer(getUser(comment)))
		doRequest(methodPost, urlComment, getBuffer(getComment(comment, asset.URL)))
	}
}

func getBuffer(object interface{}) *bytes.Buffer {
	b, _ := json.Marshal(object)
	return bytes.NewBuffer(b)
}

func getAsset(m objects.Map) model.Asset {
	asset := model.Asset{}
	url := getOne(m.Get("object.context")).GetString("uri")
	asset.URL = url
	asset.SourceID = url
	return asset
}

func getUser(m objects.Map) model.User {
	user := model.User{}

	user.SourceID = m.GetString("actor.id")
	user.UserName = m.GetString("actor.title")
	user.Status = m.GetString("actor.status")
	user.Avatar = m.GetString("actor.avatar")

	return user
}

func getComment(m objects.Map, url string) model.Comment {
	comment := model.Comment{}

	comment.Body = m.GetString("object.content")
	comment.Status = m.GetString("object.status")

	t, _ := time.Parse(time.RFC3339, m.GetString("postedTime"))
	comment.DateCreated = t

	t, _ = time.Parse(time.RFC3339, m.GetString("updated"))
	comment.DateUpdated = t

	comment.Source.ID = findCommentId(m.GetString("id"))
	comment.Source.AssetID = url
	comment.Source.UserID = m.GetString("actor.id")
	parentID := getOne(m.Get("targets")).GetString("conversationID")
	comment.Source.ParentID = findCommentID(parentID)

	return comment
}

//returns the first one from an array of json documents
func getOne(list interface{}) objects.Map {

	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice:
		slice := reflect.ValueOf(list)

		//must convert the Interface to map[string]interface{}
		//so that it can be converted to an objects.Map
		var m map[string]interface{}
		m = slice.Index(0).Interface().(map[string]interface{})

		return objects.Map(m)
	}
	return nil
}

//str := "http://washingtonpost.com/ECHO/item/c5c8f176-3f27-4228-94d1-8dffc73028ac"
//str = "http://js-kit.com/activities/post/c5c8f176-3f27-4228-94d1-8dffc73028ac"

func findCommentID(url string) string {
	var s []string
	s = strings.Split(url, "/")
	return s[(len(s) - 1)]
}
