package web

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/coralproject/pillar/pkg/amqp"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/statsd"
	"github.com/gorilla/mux"
)

// overwrites the handlerfunc
type HandlerFunc func(c *AppContext)

func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c := NewContext(rw, r)
	defer c.Close()

	//delegates to the actual handler code
	h(c)
}

//AppContext encapsulates application specific runtime information
type AppContext struct {
	Writer http.ResponseWriter
	Header http.Header
	Body   io.ReadCloser
	Vars   map[string]string
	MDB    *db.MongoDB
	RMQ    *amqp.MQ
	SD     *statsd.SD
	Event  string
}

func (c *AppContext) Close() {
	c.MDB.Close()
	c.RMQ.Close()
	c.SD.Close()
}

func (c *AppContext) GetValue(key string) string {
	return c.Vars[key]
}

func (c *AppContext) SetValue(k string, v string) string {

	c.Vars[k] = v
	return c.Vars[k]
}

//Returns a cloned context with db and mq resources
//A cloned context must not be closed
func (c *AppContext) Clone() *AppContext {
	var ac AppContext
	ac.MDB = c.MDB
	ac.RMQ = c.RMQ
	ac.SD = c.SD
	return &ac
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

	var c AppContext
	c.Writer = rw
	c.MDB = db.NewMongoDB(os.Getenv("MONGODB_URL"))
	c.RMQ = amqp.NewMQ(os.Getenv("AMQP_URL"), os.Getenv("AMQP_EXCHANGE"))
	c.SD = statsd.NewSD(os.Getenv("STATSD_URL"))

	if r != nil {
		c.Header = r.Header
		c.Body = r.Body
		c.Vars = mux.Vars(r)

		// append the query vars to the vars array
		for k, v := range r.URL.Query() {
			c.Vars[k] = v[0]
		}

	} else {
		c.Vars = make(map[string]string)
	}

	return &c
}

// AppError encapsulates application specific error
type AppError struct {
	Error   error  `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}
