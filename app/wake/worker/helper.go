package worker

import (
	"encoding/json"
)

//MnU marshalls and unmarshalls to a concrete structure
//You must pass a pointer to the concrete structure
func MnU(payload interface{}, object interface{}) error {

	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, object)
	if err != nil {
		return err
	}

	return nil
}
