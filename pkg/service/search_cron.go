package service

import (
	"log"

	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
	"gopkg.in/mgo.v2/bson"
)

// func init() {
// 	log.Printf("Xenia URL: %s\n", os.Getenv("XENIA_URL"))
// 	log.Printf("Xenia Auth: %s\n", os.Getenv("XENIA_AUTH"))
// }

func UpdateSearch() {

	log.Printf("New scheduled job - UpdateSearch!\n")
	c := web.NewContext(nil, nil)
	defer c.Close()

	searches := []model.Search{}
	c.MDB.DB.C(model.Searches).Find(nil).All(&searches)

	for _, one := range searches {
		doUpdateSearch(c, one)
	}
}

func doUpdateSearch(c *web.AppContext, search model.Search) {

	log.Printf("Starting UpdateSearch: %s", search.Name)
	//map of new users from search
	m, a, err := getNewUsers(c.MDB, search)
	if err != nil {
		log.Printf("UpdateSearch failed: %s", search.Name)
		return

	}

	//remove tag when user from old list are no longer in new list
	for _, one := range search.Result.Users {
		if _, ok := m[one.ID]; !ok {
			if user := removeTag(c, one.ID, search.Tag); user != nil {
				p := model.Event{model.EventTagRemoved, model.PayloadTag{search.Tag, *user}}
				PublishEvent(c, nil, p)
			}
		}
	}

	// add tags to users who are returned by the search
	for _, value := range m {
		if user := addTag(c, value.ID, search.Tag); user != nil {
			p := model.Event{model.EventTagAdded, model.PayloadTag{search.Tag, *user}}
			PublishEvent(c, nil, p)
		}
	}

	//save new users to search.results
	r := model.SearchResult{Count: len(m), Users: a}
	c.MDB.DB.C(model.Searches).UpdateId(search.ID, bson.M{"$set": bson.M{"result": r}})
	log.Printf("UpdateSearch successful [query: %v, count %d]\n", search.Query, len(m))
	c.SD.Client.Inc("Update_Search_Successful", 1, 1.0)
}

func addTag(c *web.AppContext, id bson.ObjectId, tag string) *model.User {
	var user model.User
	if err := c.MDB.DB.C(model.Users).FindId(id).One(&user); err != nil {
		return nil
	}

	for _, one := range user.Tags {
		if one == tag {
			return nil //return if the tag exists
		}
	}

	//add the new tag
	log.Printf("Adding tag [%s] to user [%s]", tag, user.ID)
	tags := append(user.Tags, tag)
	c.MDB.DB.C(model.Users).UpdateId(id, bson.M{"$set": bson.M{"tags": tags}})
	c.SD.Client.Inc("New_Tag_Added", 1, 1.0)
	return &user
}

func removeTag(c *web.AppContext, id bson.ObjectId, tag string) *model.User {
	var user model.User
	if err := c.MDB.DB.C(model.Users).FindId(id).One(&user); err != nil {
		return nil
	}

	var tags []string
	for _, one := range user.Tags {
		if one == tag {
			continue //skip the one already present
		}

		tags = append(tags, one)
	}

	log.Printf("Removing tag [%s] from user [%s]", tag, user.ID)
	c.MDB.DB.C(model.Users).UpdateId(id, bson.M{"$set": bson.M{"tags": tags}})
	c.SD.Client.Inc("Tag_Removed", 1, 1.0)
	return &user
}

//returns new sets of users from the search
func getNewUsers(db *db.MongoDB, search model.Search) (map[bson.ObjectId]model.User, []model.User, error) {

	ids, err := getUserIds(search)
	if err != nil {
		return nil, nil, err
	}

	m := make(map[bson.ObjectId]model.User, len(ids))
	a := make([]model.User, len(ids))

	for i := 0; i < len(ids); i++ {
		var user model.User
		key := bson.ObjectIdHex(ids[i])
		db.DB.C(model.Users).FindId(key).One(&user)
		m[key] = user
		a[i] = user
	}

	return m, a, nil
}
