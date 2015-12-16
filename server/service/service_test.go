package service

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"github.com/coralproject/pillar/server/model"
)

func TestCreateAsset(t *testing.T) {
	file, err := os.Open("../../data/assets.json")
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
	file, err := os.Open("../../data/users.json")
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
