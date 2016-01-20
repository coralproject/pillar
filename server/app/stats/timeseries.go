package main

import (
	"errors"
	"time"

	//"github.com/coralproject/pillar/server/model"

	//	"gopkg.in/mgo.v2"

	"github.com/ardanlabs/kit/log"
)

var (
	DURATIONS = make(map[string]time.Duration)
)

func initDurations() {

	// master duration strings, to be moved to config
	durationStrings := map[string]string{
		"hour": "1h",
		"day":  "24h",
		"week": "168h",
	}

	// parse strings into durations, build master DURATIONS map
	for k, v := range durationStrings {

		d, err := time.ParseDuration(v)
		if err != nil {
			log.Error(context, "initDurations", errors.New("Could not parse Duration"), "Could not parse duration %+v", d)
		} else {

			DURATIONS[k] = d
		}

	}

	log.User(context, "initDurations", "%+v", DURATIONS)

}

func buildTimeseries() {

}
