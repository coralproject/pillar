package stats

import (
	//"errors"
	"fmt"

	"github.com/coralproject/pillar/server/model"
	"github.com/coralproject/pillar/server/service"

	//"github.com/ardanlabs/kit/log"
	"gopkg.in/mgo.v2/bson"
)

func onCreateComment(p map[string]string) error {

	db := service.GetMongoManager()
	defer db.Close()

	var c model.Comment

	//return, if exists
	db.Comments.FindId(bson.ObjectIdHex("567b0850e19ac8852dd2bb5c")).One(&c)

	//	db.DB("coral").C("comment").FindId(bson.ObjectId("567b0850e19ac8852dd2bb5c")).One(&c)

	fmt.Printf("TEST %#v\n\n", c)

	return nil

}

/*
func updateCommentOnAssetCount(_id string) error {

	return nil

}
*/
