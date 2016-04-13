package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
	"log"
	"os"
	"testing"
)

//various constants
const (
	MethodGet    string = "GET"
	MethodPost   string = "POST"
	MethodOption string = "OPTIONS"

	URLUser    string = "api/import/user"
	URLAsset   string = "api/import/asset"
	URLAction  string = "api/import/action"
	URLComment string = "api/import/comment"
	URLTag     string = "api/tag"
	URLTags    string = "api/tags"

	DataUsers    = "fixtures/users.json"
	DataAssets   = "fixtures/assets.json"
	DataComments = "fixtures/comments.json"
	DataActions  = "fixtures/actions.json"
	DataTags     = "fixtures/tags.json"
)

var baseURL string

func getBaseURL() string {
	return baseURL
}

func getHeader() {
	m := make(map[string]string)
	m["Content-Type"] = "application/json"

	return m
}

func init() {
	baseURL = os.Getenv("PILLAR_URL")
	if baseURL == "" {
		log.Fatalf("Error connecting to Pillar: PILLAR_URL not found.")
	}

	db := db.NewMongoDB()
	defer db.Close()

	//Empty all test data
	db.TagTargets.RemoveAll(nil)
	db.Tags.RemoveAll(nil)
	db.Actions.RemoveAll(nil)
	db.Comments.RemoveAll(nil)
	db.Users.RemoveAll(nil)
	db.Assets.RemoveAll(nil)
	db.CayUserActions.RemoveAll(nil)
}

func TestCORS(t *testing.T) {
	if r, _ := web.Request(MethodOption, getBaseURL()+URLTag,
		getHeader(), nil); r.StatusCode != 200 {
		t.Fail()
	}

	if r, _ := web.Request(MethodOption, getBaseURL()+URLTags,
		getHeader(), nil); r.StatusCode != 200 {
		t.Fail()
	}

	if r, _ := web.Request(MethodOption, getBaseURL()+URLAsset,
		getHeader(), nil); r.StatusCode != 200 {
		t.Fail()
	}

	if r, _ := web.Request(MethodOption, getBaseURL()+URLUser,
		getHeader(), nil); r.StatusCode != 200 {
		t.Fail()
	}

	if r, _ := web.Request(MethodOption, getBaseURL()+URLComment,
		getHeader(), nil); r.StatusCode != 200 {
		t.Fail()
	}
}

func TestCreateAssets(t *testing.T) {
	file, err := os.Open(DataAssets)
	if err != nil {
		fmt.Printf("Error reading asset data [%s]", err.Error())
	}

	var objects []model.Asset
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading asset data", err.Error())
	}

	for _, one := range objects {
		data, _ := json.Marshal(one)
		if r, _ := web.Request(MethodPost, getBaseURL()+URLAsset,
			getHeader(), bytes.NewBuffer(data)); r.StatusCode != 200 {
			t.Fail()
		}
	}
}

func TestCreateUsers(t *testing.T) {
	file, err := os.Open(DataUsers)
	if err != nil {
		fmt.Printf("Error reading user data [%s]", err.Error())
	}

	objects := []model.User{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading user data", err.Error())
	}

	for _, one := range objects {
		data, _ := json.Marshal(one)
		if r, _ := web.Request(MethodPost, getBaseURL()+URLUser,
			getHeader(), bytes.NewBuffer(data)); r.StatusCode != 200 {
			t.Fail()
		}
	}
}

func TestCreateComments(t *testing.T) {
	file, err := os.Open(DataComments)
	if err != nil {
		fmt.Printf("Error reading comment data [%s]", err.Error())
	}

	objects := []model.Comment{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading comment data", err.Error())
	}

	for _, one := range objects {
		data, _ := json.Marshal(one)
		if r, _ := web.Request(MethodPost, getBaseURL()+URLComment,
			getHeader(), bytes.NewBuffer(data)); r.StatusCode != 200 {
			t.Fail()
		}
	}
}

func TestCreateActions(t *testing.T) {
	file, err := os.Open(DataActions)
	if err != nil {
		fmt.Printf("Error reading action data [%s]", err.Error())
	}

	objects := []model.Action{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading action data", err.Error())
	}

	for _, one := range objects {
		data, _ := json.Marshal(one)
		if r, _ := web.Request(MethodPost, getBaseURL()+URLAction,
			getHeader(), bytes.NewBuffer(data)); r.StatusCode != 200 {
			t.Fail()
		}
	}
}

func TestCreateTags(t *testing.T) {
	file, err := os.Open(DataTags)
	if err != nil {
		fmt.Printf("Error reading tag data [%s]", err.Error())
	}

	objects := []model.Tag{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading action data", err.Error())
	}

	for _, one := range objects {
		data, _ := json.Marshal(one)
		if r, _ := web.Request(MethodPost, getBaseURL()+URLTag,
			getHeader(), bytes.NewBuffer(data)); r.StatusCode != 200 {
			t.Fail()
		}
	}
}
