package main

import (
	"errors"
	"time"

	//"github.com/coralproject/pillar/server/model"

	"github.com/ardanlabs/kit/log"
)

func getDurations() map[string]time.Duration {

	ds := make(map[string]time.Duration)

	// master duration strings, to be moved to config
	durationStrings := map[string]string{
		//		"hour":  "1h",
		"day":   "24h",
		"week":  "168h",
		"month": "720h",
	}

	// parse strings into durations, build master DURATIONS map
	for k, v := range durationStrings {

		d, err := time.ParseDuration(v)
		if err != nil {
			log.Error(uid, "initDurations", errors.New("Could not parse Duration"), "Could not parse duration %+v", d)
		} else {

			ds[k] = d
		}

	}

	return ds

}

// func getCollections() map[string]*mgo.Collection {
//
// 	collections := map[string]*mgo.Collection{
// 		"comment": db.C("comment"),
// 	}
//
// 	return collections
//
// }
