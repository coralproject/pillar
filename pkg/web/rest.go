package web

import (
	"io"
	"io/ioutil"
	"net/http"
)

const (
	GET    string = "GET"
	POST   string = "POST"
	PUT    string = "PUT"
	DELETE string = "DELETE"
)

//Response encapsulates a http response
type Response struct {
	Status     string
	Header     http.Header
	Body       string
	StatusCode int
}

//Request is a common method for a REST call, returns a Response
func Request(method string, url string, header map[string]string, payload io.Reader) (*Response, error) {

	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	//set headers for the request
	for k, v := range header {
		request.Header.Add(k, v)
	}

	var client http.Client
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	resBody, _ := ioutil.ReadAll(response.Body)

	return &Response{
		response.Status,
		response.Header,
		string(resBody),
		response.StatusCode,
	}, nil
}
