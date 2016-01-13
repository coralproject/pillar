package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/coralproject/pillar/server/model"
	"github.com/coralproject/pillar/client/rest"
	"testing"
)

func TestCreateAssets(t *testing.T) {
	file, err := os.Open(rest.DataAssets)
	if err != nil {
		fmt.Printf("Error reading asset data [%s]", err.Error())
	}

	objects := []model.Asset{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading asset data", err.Error())
	}

	for _, one := range objects {
		data, _ := json.Marshal(one)
		rest.Request(rest.MethodPost, rest.URLAsset, bytes.NewBuffer(data))
	}
}

func TestCreateUsers(t *testing.T) {
	file, err := os.Open(rest.DataUsers)
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
		rest.Request(rest.MethodPost, rest.URLUser, bytes.NewBuffer(data))
	}
}

func TestCreateComments(t *testing.T) {
	file, err := os.Open(rest.DataComments)
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
		rest.Request(rest.MethodPost, rest.URLComment, bytes.NewBuffer(data))
	}
}

func TestCreateActions(t *testing.T) {
	file, err := os.Open(rest.DataActions)
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
		rest.Request(rest.MethodPost, rest.URLAction, bytes.NewBuffer(data))
	}
}
