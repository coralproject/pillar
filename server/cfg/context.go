package cfg

import (
	"os"
	"log"
	"fmt"
)

//Context information for pillar server
//Server expects a pillar.json under PILLAR_HOME to bootstrap
type Context struct {
	Home      string `json:"home" bson:"home"`
	MongoURL  string `json:"mongo_url" bson:"mongo_url"` //mongodb://localhost:27017/coral
}

var context Context

//export PILLAR_HOME=path to pillar home
func init() {
	home := os.Getenv("PILLAR_HOME")
	if home == "" {
		log.Fatal("Error initializing Server: PILLAR_HOME not found.")
	}
	//set pillar home
	context.Home = home

	url := os.Getenv("MONGODB_URL")
	if home == "" {
		log.Fatal("Error initializing Server: MONGODB_URL not found.")
	}
	context.MongoURL = url

	fmt.Printf("Context: %+v\n\n", context)
}

func GetContext() Context {
	return context
}
