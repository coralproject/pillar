package main

import (
	"fmt"
	"time"

	"github.com/ardanlabs/kit/log"

	"github.com/coralproject/pillar/server/model"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Timeslice struct {
	Id        bson.ObjectId `json:"id" bson:"_id"`
	Start     int64         `json:"s" bson:"start"`
	Duration  string        `json:"d" bson:"duration"`
	Target    string        `json:"t" bson:"target"`
	TargetDoc interface{}   `json:"tdoc" bson:"target_doc"`
	Data      interface{}   `json:"data" bson:"data"`
}

func getRange(c *mgo.Collection) (time.Time, time.Time) {

	// get start time

	var first model.Comment
	err := c.Find(nil).Limit(1).Sort("date_created").One(&first)
	if err != nil {
		// no documents in this collection
		log.Dev(context, "buildTimeseries", "Collection %v is empty", c)

	}

	start := first.DateCreated

	// get end time

	var last model.Comment
	err = c.Find(nil).Limit(1).Sort("-date_created").One(&last)
	if err != nil {

	}

	end := last.DateCreated

	return start, end

}

func calc(cs []*model.Comment) map[string]interface{} {

	d := make(map[string]interface{})

	d["comments"] = len(cs)

	return d

}

func buildTimeseries(cs map[string]*mgo.Collection, ds map[string]time.Duration) {

	log.User(context, "buildTimeseries", "Building timeseries: collections %v, durations %v", cs, ds)

	limiter := 0
	limit := 1000

	// let's start fresh here
	tc := db.C("comment_timeseries")
	tc.DropCollection()

	// for each collection
	for k, c := range cs {

		// for each duration
		for dk, d := range ds {

			log.Dev(context, "buildTimeseries", "Building %v from %v", k, c)

			start, end := getRange(c)
			log.Dev(context, "buildTimeseries", "From %+v to %+v", start, end)

			// range over timeseries

			t := start
			for t.Before(end) {

				comments := []*model.Comment{}
				err := c.Find(bson.M{"date_created": bson.M{"$gte": t, "$lt": t.Add(d)}}).All(&comments)
				if err != nil {

				}

				t = t.Add(d)

				data := calc(comments)
				if data["comments"] == 0 {
					continue
				}

				timeslice := Timeslice{
					Id:       bson.NewObjectId(),
					Start:    t.Unix(),
					Duration: dk,
					Target:   "total",
					Data:     data,
				}

				err = tc.Insert(timeslice)
				if err != nil {
					log.Error(context, "Build", err, "Error inserting %+v\n", timeslice)
				}

				fmt.Printf("%v %+v\n", limiter, timeslice)

				// user map
				us := make(map[bson.ObjectId][]*model.Comment)
				// assets map
				as := make(map[bson.ObjectId][]*model.Comment)
				for _, comment := range comments {
					us[comment.UserID] = append(us[comment.UserID], comment)
					as[comment.AssetID] = append(as[comment.AssetID], comment)
				}

				// range through the user targets and insert
				for tid, tcomments := range us {
					var user model.User
					db.C("user").Find(bson.M{"_id": tid}).One(&user)
					timeslice = Timeslice{
						Id:        bson.NewObjectId(),
						Start:     t.Unix(),
						Duration:  dk,
						Target:    "user",
						TargetDoc: user,
						Data:      calc(tcomments),
					}
					err = tc.Insert(timeslice)
					if err != nil {
						log.Error(context, "Build", err, "Error inserting %+v\n", timeslice)
					}

				}

				// range through the asset targets and insert
				for tid, tcomments := range as {
					var asset model.Asset
					db.C("asset").Find(bson.M{"_id": tid}).One(&asset)
					timeslice = Timeslice{
						Id:        bson.NewObjectId(),
						Start:     t.Unix(),
						Duration:  dk,
						Target:    "asset",
						TargetDoc: asset,
						Data:      calc(tcomments),
					}
					err = tc.Insert(timeslice)
					if err != nil {
						log.Error(context, "Build", err, "Error inserting %+v\n", timeslice)
					}

				}

				// sanity
				limiter++
				if limiter >= limit {
					//break

				}

			}
		}

	}

}
