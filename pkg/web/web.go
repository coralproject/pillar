package web

import (
	"bytes"
	"encoding/json"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

type HandlerFunc func(c *AppContext)

func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c := NewContext(rw, r)
	defer c.Close()

	//delegates to the actual handler code
	h(c)
}

// AppContext encapsulates application specific runtime information
type AppContext struct {
	Writer http.ResponseWriter
	Header http.Header
	Body   io.ReadCloser
	Vars   map[string]string
	DB     *db.MongoDB
}

func (c *AppContext) Close() {
	c.DB.Close()
}

func (c *AppContext) GetValue(key string) string {
	return c.Vars[key]
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

func NewContext(rw http.ResponseWriter, r *http.Request) *AppContext {

	if r == nil {
		return &AppContext{rw, nil, nil, mux.Vars(r), db.NewMongoDB()}
	}

	return &AppContext{rw, r.Header, r.Body, mux.Vars(r), db.NewMongoDB()}
}

// AppError encapsulates application specific error
type AppError struct {
	Error   error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}
