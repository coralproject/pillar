package fiddler

import (
	"bytes"
	"encoding/json"
)

func getBuffer(object interface{}) *bytes.Buffer {
	b, _ := json.Marshal(object)
	return bytes.NewBuffer(b)
}
