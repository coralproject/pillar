package main

import (
	"errors"
	"os"
	"time"

	//"github.com/coralproject/pillar/server/model"

	"gopkg.in/mgo.v2"

	"github.com/ardanlabs/kit/log"
)

var (
	context int32
	db      *mgo.Session
)

//export MONGODB_URL=mongodb://localhost:27017/coral
func initDb() *mgo.Session {
	uri := os.Getenv("MONGODB_URL")
	if uri == "" {
		log.Error("start", "init", errors.New("Error connecting - MONGODB_URL not found!"), "Getting MONGODB_URL env variable.")
		os.Exit(1)
	}

	db, err := mgo.Dial(uri)
	if err != nil {
		log.Error("start", "init", err, "Connecting to mongo")
		panic(err) // no, not really <--- do we really need to panic?
	}

	return db

}

func init() {

	logLevel := func() int {
		return log.DEV
	}

	log.Init(os.Stdout, logLevel)

	context = int32(time.Now().Unix())

	log.User(context, "init", "Initializing")

	db = initDb()

}

func main() {

	log.User(context, "main", "Beginning main %+v", db)

	initDurations()

}
