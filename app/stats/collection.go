package main

import (
	//"fmt"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/pillar/server/model"

	"github.com/stretchr/stew/objects"
	"gopkg.in/mgo.v2/bson"
)

func calcCollectionStats(c CollectionStats) {

	var docs []map[string]interface{} // slice of documents we're working with, in generic form

	//	err := c.Collection.Find(bson.M{"stats": bson.M{"$exists": true}}).Sort("-_id").Limit(100).All(&docs)
	err := c.Collection.Find(nil).Sort("-_id").Limit(100).All(&docs)
	if err != nil {
		log.Error("stats", "request", err, "Could not load assets")
	}

	for _, doc := range docs {

		comments := make([]model.Comment, 0)

		id := objects.Map(doc).Get("_id")

		err := db.C("comment").Find(bson.M{c.ForeignKeyField: id}).Sort("-_id").All(&comments)
		if err != nil {
			log.Error("stats", "request", err, "Could not load user comments")
		}

		stats := calc(comments)

		c.Collection.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"stats": stats}})

		//fmt.Printf("\n\n%+v\n\n%+v", doc, stats)

	}

}
