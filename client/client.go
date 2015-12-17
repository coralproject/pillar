package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/coralproject/pillar/server/model"
)

const methodGet string = "GET"
const methodPost string = "POST"

const url string = "http://localhost:8080/api/import/"
const urlUser string = url + "user"
const urlAsset string = url + "asset"
const urlComment string = url + "comment"

const dataUsers = "data/users.json"
const dataAssets = "data/assets.json"
const dataComments = "data/comments.json"

type restResponse struct {
	Status     string
	Header     http.Header
	Payload    string
	StatusCode int
}

func main() {

	//	//insert assets
	//	addAssets()
	//
	//	//insert users
	//	addUsers()
	//
	//	//insert comments
	//	addComments()
	wapoFiddler()
}

func addAssets() {
	file, err := os.Open(dataAssets)
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
		doRequest(methodPost, urlAsset, bytes.NewBuffer(data))
	}
}

func addUsers() {
	file, err := os.Open(dataUsers)
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
		doRequest(methodPost, urlUser, bytes.NewBuffer(data))
	}
}

func addComments() {
	file, err := os.Open(dataComments)
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
		doRequest(methodPost, urlComment, bytes.NewBuffer(data))
	}
}

func doRequest(method string, urlStr string, payload io.Reader) restResponse {

	request, err := http.NewRequest(method, urlStr, payload)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error in processing request [%s]", err.Error())
	}
	defer response.Body.Close()

	resBody, _ := ioutil.ReadAll(response.Body)

	rest := restResponse{
		response.Status,
		response.Header,
		string(resBody),
		response.StatusCode,
	}

	return rest
}
