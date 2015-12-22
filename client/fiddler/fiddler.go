package fiddler
import (
	"encoding/json"
	"bytes"
)


func getBuffer(object interface{}) *bytes.Buffer {
	b, _ := json.Marshal(object)
	return bytes.NewBuffer(b)
}
