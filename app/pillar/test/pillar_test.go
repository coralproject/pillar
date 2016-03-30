package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"os"
	"testing"
)

func init() {
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
	if r, _ := request(MethodOption, getBaseURL()+URLTag, nil); r.StatusCode != 200 {
		t.Fail()
	}

	if r, _ := request(MethodOption, getBaseURL()+URLTags, nil); r.StatusCode != 200 {
		t.Fail()
	}

	if r, _ := request(MethodOption, getBaseURL()+URLAsset, nil); r.StatusCode != 200 {
		t.Fail()
	}

	if r, _ := request(MethodOption, getBaseURL()+URLUser, nil); r.StatusCode != 200 {
		t.Fail()
	}

	if r, _ := request(MethodOption, getBaseURL()+URLComment, nil); r.StatusCode != 200 {
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
		if r, _ := request(MethodPost, getBaseURL()+URLAsset, bytes.NewBuffer(data)); r.StatusCode != 200 {
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
		if r, _ := request(MethodPost, getBaseURL()+URLUser, bytes.NewBuffer(data)); r.StatusCode != 200 {
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
		if r, _ := request(MethodPost, getBaseURL()+URLComment, bytes.NewBuffer(data)); r.StatusCode != 200 {
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
		if r, _ := request(MethodPost, getBaseURL()+URLAction, bytes.NewBuffer(data)); r.StatusCode != 200 {
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
		if r, _ := request(MethodPost, getBaseURL()+URLTag, bytes.NewBuffer(data)); r.StatusCode != 200 {
			t.Fail()
		}
	}
}
