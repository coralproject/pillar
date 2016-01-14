package fiddler

import (
	"encoding/json"
	"fmt"
	"github.com/coralproject/pillar/client/db"
	"github.com/coralproject/pillar/client/rest"
	"github.com/coralproject/pillar/server/model"
	"github.com/stretchr/stew/objects"
)

//LoadActors imports actors
func LoadActors() {

	manager := db.GetMongoManager()
	defer manager.Close()

	all := make([]interface{}, 10)

	manager.Actors.Find(nil).All(&all)

	fmt.Printf("Found %d Actors\n", len(all))
	fmt.Printf("Import in progress...\n")
	var nSuccess, nFailure int
	for _, one := range all {
		data, _ := json.Marshal(one)

		user := map[string]interface{}{}
		json.Unmarshal(data, &user)
		m := objects.Map(user)
		if response := rest.Request(rest.MethodPost, rest.URLUser, getBuffer(getActor(m))); response.StatusCode == 200 {
			nSuccess++
		} else {
			nFailure++
		}
	}
	fmt.Printf("Finished importing actors, suceess: [%d] failure: [%d]\n\n\n", nSuccess, nFailure)
}

func getActor(m objects.Map) model.User {
	user := model.User{}

	user.Source.ID = m.GetStringOrEmpty("_id")
	user.Name = m.GetStringOrEmpty("title")
	user.Status = m.GetStringOrEmpty("status")
	user.Avatar = m.GetStringOrEmpty("avatar")

	return user
}
