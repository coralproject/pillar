package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/server/model"
	"github.com/stretchr/stew/objects"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"strings"
	"time"
)

func wapoFiddler() {

	//story := "http://washingtonpost.com/posteverything/wp/2015/05/20/feminists-want-us-to-define-these-ugly-sexual-encounters-as-rape-dont-let-them/"
	//story := "http://washingtonpost.com/opinions/reformers-want-to-erase-confuciuss-influence-in-asia-thats-a-mistake/2015/05/28/529c1d3a-042e-11e5-a428-c984eb077d4e_story.html"
	story := "http://washingtonpost.com/world/europe/european-leaders-seek-last-ditch-offer-to-bring-greece-from-brink-of-default/2015/06/30/960aded8-1ea2-11e5-a135-935065bc30d0_story.html"

	manager := getMongoManager()
	defer manager.close()

	all := make([]interface{}, 10)

	manager.comments.Find(bson.M{"object.context.uri": story}).Sort("postedTime").All(&all)

	fmt.Printf("Found %d comments\n", len(all))
	fmt.Printf("Import in progress...\n")
	var nAssets, nUsers, nComments int
	for _, one := range all {
		data, _ := json.Marshal(one)

		comment := map[string]interface{}{}
		json.Unmarshal(data, &comment)

		asset := getAsset(comment)
		if response := doRequest(methodPost, urlAsset, getBuffer(asset)); response.StatusCode == 200 {
			nAssets++
		}

		if response := doRequest(methodPost, urlUser, getBuffer(getUser(comment))); response.StatusCode == 200 {
			nUsers++
		}

		if response := doRequest(methodPost, urlComment, getBuffer(getComment(comment, asset.URL))); response.StatusCode == 200 {
			nComments++
		}
	}
	fmt.Printf("Finished importing: Comments[%d]\n\n\n", nComments)
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

	user.SourceID = m.GetStringOrEmpty("actor.id")
	user.UserName = m.GetStringOrEmpty("actor.title")
	user.Status = m.GetStringOrEmpty("actor.status")
	user.Avatar = m.GetStringOrEmpty("actor.avatar")

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

	comment.Source.ID = getShortCommentID(m.GetString("id"))
	comment.Source.AssetID = url
	comment.Source.UserID = m.GetString("actor.id")

	target := getOne(m.Get("targets"))
	targetID := getShortCommentID(target.GetString("id"))
	//	targetConversationID := getShortCommentID(target.GetString("conversationID"))
	//	fmt.Printf("ID: %s\n", comment.Source.ID)
	//	fmt.Printf("targetID: %s\n", targetID)
	//	fmt.Printf("targetConversationID: %s\n\n\n", targetConversationID)

	comment.Source.ParentID = getShortCommentID(targetID)

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
func getShortCommentID(url string) string {
	var s []string
	s = strings.Split(url, "/")
	return s[(len(s) - 1)]
}
