package test

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

//various constants
const (
	MethodGet    string = "GET"
	MethodPost   string = "POST"
	MethodOption string = "OPTIONS"

	BaseURL    string = "http://localhost:8080/api/"
	URLUser    string = BaseURL + "import/user"
	URLAsset   string = BaseURL + "import/asset"
	URLAction  string = BaseURL + "import/action"
	URLComment string = BaseURL + "import/comment"
	URLTag     string = BaseURL + "tag"
	URLTags    string = BaseURL + "tags"

	DataUsers    = "fixtures/users.json"
	DataAssets   = "fixtures/assets.json"
	DataComments = "fixtures/comments.json"
	DataActions  = "fixtures/actions.json"
	DataTags     = "fixtures/tags.json"
)

//Response encapsulates a REST response
type Response struct {
	Status     string
	Header     http.Header
	Payload    string
	StatusCode int
}

//Request is a common method for a REST call, returns a Response
func request(method string, url string, payload io.Reader) (*Response, error) {

	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	request.Header.Set("Content-Type", "application/json")

	var client http.Client
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	resBody, _ := ioutil.ReadAll(response.Body)

	rest := Response{
		response.Status,
		response.Header,
		string(resBody),
		response.StatusCode,
	}

	return &rest, nil
}
