package service

import (
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/model"
	"os"
	"testing"
)

func TestCreateUser(t *testing.T) {
	fmt.Println("In CreateUsers")
	file, err := os.Open("users.json")
	if err != nil {
		fmt.Println("opening config file", err.Error())
	}

	users := []model.User{}
	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&users); err != nil {
		fmt.Println("Error reading user data", err.Error())
	}

	for _, one := range users {
		_, err := CreateUser(one)
		if err != nil {
			t.Fail()
		}
	}
}
