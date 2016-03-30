package test

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

//various constants
const (
	MethodGet    string = "GET"
	MethodPost   string = "POST"
	MethodOption string = "OPTIONS"

	URLUser    string = "api/import/user"
	URLAsset   string = "api/import/asset"
	URLAction  string = "api/import/action"
	URLComment string = "api/import/comment"
	URLTag     string = "api/tag"
	URLTags    string = "api/tags"

	DataUsers    = "fixtures/users.json"
	DataAssets   = "fixtures/assets.json"
	DataComments = "fixtures/comments.json"
	DataActions  = "fixtures/actions.json"
	DataTags     = "fixtures/tags.json"
)

var baseURL string

func init() {
	baseURL = os.Getenv("PILLAR_URL")
	if baseURL == "" {
		log.Fatalf("Error connecting to Pillar: PILLAR_URL not found.")
	}
	log.Printf("BaseURL: %s\n\n", getBaseURL())
}

func getBaseURL() string {
	return baseURL
}

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
