package service

import (
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/server/model"
	"os"
	"testing"
)

const dataUsers = "../../data/users.json"
const dataAssets = "../../data/assets.json"
const dataComments = "../../data/comments.json"

func TestCreateAsset(t *testing.T) {
	file, err := os.Open(dataAssets)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []model.Asset{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading asset data", err.Error())
	}

	for _, one := range objects {
		_, err := CreateAsset(one)
		if err != nil {
			t.Fail()
		}
	}
}

func TestCreateUser(t *testing.T) {
	file, err := os.Open(dataUsers)
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	objects := []model.User{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&objects); err != nil {
		fmt.Println("Error reading user data", err.Error())
	}

	for _, one := range objects {
		_, err := CreateUser(one)
		if err != nil {
			t.Fail()
		}
	}
}
