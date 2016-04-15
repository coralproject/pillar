package service

import (
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
	"gopkg.in/mgo.v2/bson"
	"log"
)

func UpdateSearch() {

	c := web.NewContext(nil, nil)
	defer c.Close()

	searches := []model.Search{}
	c.DB.Searches.Find(nil).All(&searches)

	for _, one := range searches {
		doUpdateSearch(c, one)
	}
}

func doUpdateSearch(c *web.AppContext, search model.Search) {
	//map of new users from search
	m, a := getNewUsers(c.DB, search)
	if m == nil || len(m) == 0 {
		log.Printf("UpdateSearch - no new users, skipping [%s]\n", search.Query)
		return
	}

	//remove tag when user from old list are no longer in new list
	for _, one := range search.Result.Users {
		if _, ok := m[one.ID]; !ok {
			if user := removeTag(c.DB, one.ID, search.Tag); user != nil {
				p := model.PayloadTag{model.EventTagRemoved, search.Tag, *user}
				PublishEvent(c, nil, p)
			}
		}
	}

	for _, value := range m {
		if user := addTag(c.DB, value.ID, search.Tag); user != nil {
			p := model.PayloadTag{model.EventTagAdded, search.Tag, *user}
			PublishEvent(c, nil, p)
		}
	}

	//save new users to search.results
	r := model.SearchResult{Count: len(m), Users: a}
	c.DB.Searches.UpdateId(search.ID, bson.M{"$set": bson.M{"result": r}})
	log.Printf("UpdateSearch successful [query: %v, count %d]\n", search.Query, len(m))
}

func addTag(db *db.MongoDB, id bson.ObjectId, tag string) *model.User {
	var user model.User
	if err := db.Users.FindId(id).One(&user); err != nil {
		return nil
	}

	for _, one := range user.Tags {
		if one == tag {
			return nil //return if the tag exists
		}
	}

	//add the new tag
	tags := append(user.Tags, tag)
	db.Users.UpdateId(id, bson.M{"$set": bson.M{"tags": tags}})
	return &user
}

func removeTag(db *db.MongoDB, id bson.ObjectId, tag string) *model.User {
	var user model.User
	if err := db.Users.FindId(id).One(&user); err != nil {
		return nil
	}

	var tags []string
	for _, one := range user.Tags {
		if one == tag {
			continue //skip the one already present
		}

		tags = append(tags, one)
	}

	db.Users.UpdateId(id, bson.M{"$set": bson.M{"tags": tags}})
	return &user
}

//returns new sets of users from the search
func getNewUsers(db *db.MongoDB, search model.Search) (map[bson.ObjectId]model.User, []model.User) {

	ids := getUserIds(search)
	if len(ids) == 0 {
		return nil, nil
	}

	m := make(map[bson.ObjectId]model.User, len(ids))
	a := make([]model.User, len(ids))
	for i := 0; i < len(ids); i++ {
		var user model.User
		key := bson.ObjectIdHex(ids[i])
		db.Users.FindId(key).One(&user)
		m[key] = user
		a[i] = user
	}

	return m, a
}

