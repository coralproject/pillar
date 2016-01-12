package cfg

import (
	"encoding/json"
	"os"
	"log"
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

	file, err := os.Open(home+"/pillar.json")
	if err != nil {
		log.Fatal("Error initializing Server: $PILLAR_HOME/pillar.json not found.")
	}

	jsonParser := json.NewDecoder(file)
	if err = jsonParser.Decode(&context); err != nil {
		log.Fatal("Error initializing Server: invalid pillar.json.")
	}

	//set pillar home
	context.Home = home
}

func GetContext() Context {
	return context
}
