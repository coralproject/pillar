package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"github.com/coralproject/pillar/server/model"
	"encoding/json"
	"bytes"
)

const methodGet string = "GET"
const methodPost string = "POST"

const url string = "http://localhost:8080/api/import/"
const urlUser string = url + "user"
const urlAsset string = url + "asset"
const urlComment string = url + "comment"

const dataUsers  = "../src/github.com/coralproject/pillar/data/users.json"
const dataAssets  = "../src/github.com/coralproject/pillar/data/assets.json"
const dataComments  = "../src/github.com/coralproject/pillar/data/comments.json"

type restResponse struct {
	status string
	header http.Header
	payload string
}

func main() {

	//insert users
	//doRequest(methodPost, urlUser, dataUsers)

	//insert assets
	addAssets()
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

func doRequest(method string, urlStr string, payload *bytes.Buffer) {

	request, err := http.NewRequest(method, urlStr, payload)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error in processing request [%s]", err.Error())
	}
	defer response.Body.Close()

	resBody, _ := ioutil.ReadAll(response.Body)

	rest := restResponse {
		response.Status,
		response.Header,
		string(resBody),
	}

	fmt.Printf("%+v\n\n", rest)
}
