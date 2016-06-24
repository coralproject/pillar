package service

import (
	"encoding/json"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
	"log"
	"os"
	"strings"
)

const (
	MaxResults int = 20
)

type body struct {
	Results []result `json:"results"`
}

type result struct {
	Name string `json:"Name"`
	Docs []doc  `json:"Docs"`
}

type doc struct {
	ID string `json:"_id"`
}

func getUserIds(search model.Search) []string {
	url := os.Getenv("XENIA_URL") + search.Query

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = os.Getenv("XENIA_AUTH")

	log.Printf("Xenia URL [%v]\n", url)
	response, err := web.Request(web.GET, url, header, nil)
	if err != nil {
		log.Printf("Error getting response from Xenia [%v]\n", err)
		return nil
	}

	var b body
	jsonParser := json.NewDecoder(strings.NewReader(response.Body))
	if err := jsonParser.Decode(&b); err != nil {
		log.Printf("Error decoding data [%s]", err.Error())
		return nil
	}

	docs := b.Results[0].Docs
	ids := make([]string, len(docs))
	for i := 0; i < len(b.Results[0].Docs); i++ {
		ids[i] = docs[i].ID
		if i == MaxResults-1 {
			break
		}
	}

	return ids
}

////when the item is an array, we must convert it to a slice
//func getArray(list interface{}) []objects.Map {
//
//	var resultArray []objects.Map
//	if list == nil {
//		return resultArray
//	}
//
//	switch reflect.TypeOf(list).Kind() {
//	case reflect.Slice:
//		slice := reflect.ValueOf(list)
//
//		//must convert the Interface to map[string]interface{}
//		//so that it can be converted to an objects.Map
//		//fmt.Printf("Size of slice: %d\n\n", slice.Len())
//		for i := 0; i < slice.Len(); i++ {
//			//var m map[string]interface{}
//			//fmt.Printf("Item: %s\n\n", slice.Index(i))
//			resultArray = append(resultArray, slice.Index(i).Interface().(map[string]interface{}))
//		}
//
//		return resultArray
//	}
//
//	return nil
//}
