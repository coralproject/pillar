package model

import (
	"fmt"
)

var (
	objectTypes = map[string]struct{}{
		"comments":   struct{}{},
		"assets":     struct{}{},
		"users":      struct{}{},
		"actions":    struct{}{},
		"tags":       struct{}{},
		"dimensions": struct{}{},
	}
)

func ObjectTypes() []string {
	objectTypeList := make([]string, 0, len(objectTypes))
	for objectType, _ := range objectTypes {
		objectTypeList = append(objectTypeList, objectType)
	}
	return objectTypeList
}

func ValidateObjectType(objectType string) error {
	if _, ok := objectTypes[objectType]; !ok {
		return fmt.Errorf("unkown object type: %s", objectType)
	}
	return nil
}

// ObjectTypeInstance is a factory method for known valid types. Any other
// input value will produce an insance of `map[string]interface{}`.
func ObjectTypeInstance(objectType string) interface{} {
	switch objectType {
	case "actions":
		return &Action{}
	case "assets":
		return &Asset{}
	case "comments":
		return &Comment{}
	case "users":
		return &User{}
	case "tags":
		return Tag{}
	case "dimensions":
		return Dimension{}
	}

	return make(map[string]interface{})
}
