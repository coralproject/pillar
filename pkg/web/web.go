package web

import (
	"bytes"
	"encoding/json"
	"github.com/coralproject/pillar/pkg/db"
	"io"
	"io/ioutil"
	"net/http"
)

type HandlerFunc func(rw http.ResponseWriter, r *http.Request, c *AppContext)

func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c := NewContext(r.Body)
	defer c.Close()
	h(rw, r, c)
}

// AppContext encapsulates application specific runtime information
type AppContext struct {
	DB     *db.MongoDB
	Body   io.ReadCloser
	Header http.Header
}

func (c *AppContext) Close() {
	c.DB.Close()
}

//Unmarshall unmarshalls the Body to the passed object
func (c *AppContext) Unmarshall(input interface{}) error {
	bytez, _ := ioutil.ReadAll(c.Body)
	if err := json.Unmarshal(bytez, input); err != nil {
		return err
	}
	return nil
}

//Marshall marshalls an incoming object and sets it to the Body
func (c *AppContext) Marshall(j interface{}) {
	bytez, _ := json.Marshal(j)
	c.Body = ioutil.NopCloser(bytes.NewReader(bytez))
}

func NewContext(body io.ReadCloser) *AppContext {
	return &AppContext{db.NewMongoDB(), body, nil}
}

// AppError encapsulates application specific error
type AppError struct {
	Error   error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}
