package service

import (
	"fmt"
	"github.com/coralproject/pillar/pkg/db"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
	"github.com/stretchr/stew/objects"
	"gopkg.in/mgo.v2/bson"
	"os"
	"reflect"
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
	if m == nil {
		fmt.Printf("Skipping this search [%s] - no new users!!!\n", search.Query)
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

func getUserIds(search model.Search) []string {
	url := os.Getenv("XENIA_URL") + search.Query

	header := make(map[string]string)
	header["Content-Type"] = "application/json"
	header["Authorization"] = os.Getenv("XENIA_AUTH")

	response, _ := web.Request(web.GET, url, header, nil)
	if response.StatusCode != 200 {
		//fmt.Printf("Error in xenia call %v", response)
		return nil
	}

	m, err := objects.NewMapFromJSON(response.Body)
	if err != nil {
		//fmt.Printf("Error in call")
		return nil
	}
	//get all items from Docs array as an array of objects.Map
	d := getArray(m.Get("results"))[0].Get("Docs")
	stats := getArray(d)
	ids := make([]string, len(stats))
	for i := 0; i < len(stats); i++ {
		ids[i] = stats[i].Get("_id").(string)
	}

	return ids
}

//when the item is an array, we must convert it to a slice
func getArray(list interface{}) []objects.Map {

	var resultArray []objects.Map
	if list == nil {
		return resultArray
	}

	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice:
		slice := reflect.ValueOf(list)

		//must convert the Interface to map[string]interface{}
		//so that it can be converted to an objects.Map
		//fmt.Printf("Size of slice: %d\n\n", slice.Len())
		for i := 0; i < slice.Len(); i++ {
			//var m map[string]interface{}
			//fmt.Printf("Item: %s\n\n", slice.Index(i))
			resultArray = append(resultArray, slice.Index(i).Interface().(map[string]interface{}))
		}

		return resultArray
	}

	return nil
}
