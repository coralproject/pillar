package worker

import (
	"bytes"
	"encoding/json"
	"github.com/coralproject/pillar/pkg/web"
	"github.com/coralproject/pillar/pkg/model"
	"log"
	"os"
	"reflect"
)

//various constants
const (
	//end-point names
	EPUser    string = "user"
	EPAsset   string = "asset"
	EPAction  string = "action"
	EPComment string = "comment"
)

func getESEndPoint(suffix string) string {
	return os.Getenv("ES_URL") + suffix
}

func isModel(i interface{}) bool {
	if i == nil {
		return false
	}

	v := reflect.ValueOf(i)
	_, isModel := v.Interface().(model.Model)
	if !isModel {
		return false
	}
	return true
}

func SyncElasticSearch(i interface{}) {

	if !isModel(i) {
		log.Printf("Error - Invlaid model found.\n")
		return
	}

	object := i.(model.Model)
	switch object.(type) {
	case *model.User:
		pushToES(getESEndPoint(EPUser)+"/"+object.Id(), object)
	case *model.Asset:
		pushToES(getESEndPoint(EPAsset)+"/"+object.Id(), object)
	case *model.Comment:
		pushToES(getESEndPoint(EPComment)+"/"+object.Id(), object)
	case *model.Action:
		pushToES(getESEndPoint(EPAction)+"/"+object.Id(), object)
	default:
	}
}

func pushToES(url string, object interface{}) {
	b, err := json.Marshal(object)
	if err != nil {
		log.Printf("Error marshalling object [%v]", object)
		return
	}

	if _, err := web.Request(web.POST, url, nil, bytes.NewBuffer(b)); err != nil {
		log.Printf("Error sending data to ES [%v]", err)
	}
}
