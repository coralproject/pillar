package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
)

const (
	MaxResults int = 1000
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

func getUserIds(search model.Search) ([]string, error) {

	mu := MaxResults
	var err error
	mus := os.Getenv("PILLAR_CRON_SEARCH_MAX_USERS")
	if mus != "" {
		mu64, err := strconv.ParseInt(mus, 10, 64)
		mu = int(mu64)
		if err != nil {
			log.Printf("Unrecognized value PILLAR_CRON_SEARCH_MAX_USERS, expecting int")
		}
	}

	url := os.Getenv("XENIA_URL") + search.Query + "?limit=" + strconv.FormatInt(int64(mu), 10)

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = os.Getenv("XENIA_AUTH")

	log.Printf("Xenia URL [%v]\n", url)
	response, err := web.Request(web.GET, url, header, nil)
	if err != nil {
		log.Printf("Error getting response from Xenia [%v]\n", err)
		return nil, err
	}

	if response.StatusCode == 404 {
		fmt.Printf("Response from Xenia: 404 [%v]\n", response.Body)
		return nil, errors.New("Response from Xenia: 404")
	}

	var b body
	jsonParser := json.NewDecoder(strings.NewReader(response.Body))
	if err := jsonParser.Decode(&b); err != nil {
		log.Printf("Error decoding data [%s]", err.Error())
		return nil, err
	}

	docs := b.Results[0].Docs
	ids := make([]string, 0)
	for i := 0; i < len(b.Results[0].Docs); i++ {

		ids = append(ids, docs[i].ID)

		if i == mu-1 {
			break
		}
	}

	return ids, nil
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
