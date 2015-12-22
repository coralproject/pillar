package rest

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const MethodGet string = "GET"
const MethodPost string = "POST"

const BaseUrl string = "http://localhost:8080/api/import/"
const UrlUser string = BaseUrl + "user"
const UrlAsset string = BaseUrl + "asset"
const UrlComment string = BaseUrl + "comment"

const dataUsers = "data/users.json"
const dataAssets = "data/assets.json"
const dataComments = "data/comments.json"

type RestResponse struct {
	Status     string
	Header     http.Header
	Payload    string
	StatusCode int
}

func Request(method string, urlStr string, payload io.Reader) RestResponse {

	request, err := http.NewRequest(method, urlStr, payload)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error in processing request [%s]", err.Error())
	}
	defer response.Body.Close()

	resBody, _ := ioutil.ReadAll(response.Body)

	rest := RestResponse{
		response.Status,
		response.Header,
		string(resBody),
		response.StatusCode,
	}

	return rest
}
