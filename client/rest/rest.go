package rest

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

//various constants
const (
	MethodGet  string = "GET"
	MethodPost string = "POST"

	BaseURL    string = "http://localhost:8080/api/import/"
	URLUser    string = BaseURL + "user"
	URLAsset   string = BaseURL + "asset"
	URLAction    string = BaseURL + "action"
	URLComment string = BaseURL + "comment"

	dataUsers    = "data/users.json"
	dataAssets   = "data/assets.json"
	dataComments = "data/comments.json"
)

//Response encapsulates a REST response
type Response struct {
	Status     string
	Header     http.Header
	Payload    string
	StatusCode int
}

//Request is a common method for a REST call, returns a Response
func Request(method string, urlStr string, payload io.Reader) Response {

	request, err := http.NewRequest(method, urlStr, payload)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Error in processing request [%s]", err.Error())
	}
	defer response.Body.Close()

	resBody, _ := ioutil.ReadAll(response.Body)

	rest := Response{
		response.Status,
		response.Header,
		string(resBody),
		response.StatusCode,
	}

	return rest
}
