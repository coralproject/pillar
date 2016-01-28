package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

//Context information for pillar server - Read-Only
type Context struct {
	Home     string `json:"home" bson:"home"`
	MongoURL string `json:"mongo_url" bson:"mongo_url"` //mongodb://localhost:27017/coral
}

var context Context

//GetContext returns the context information for pillar
func GetContext() Context {
	return context
}

//Single logger for the application
var (
	Logger *log.Logger
)

//export PILLAR_HOME=path to pillar home
func init() {
	home := os.Getenv("PILLAR_HOME")
	if home == "" {
		log.Fatal("Error initializing Server: PILLAR_HOME not found.")
	}
	//set pillar home
	context.Home = home

	url := os.Getenv("MONGODB_URL")
	if home == "" {
		log.Fatal("Error initializing Server: MONGODB_URL not found.")
	}
	context.MongoURL = url

	logFile := getLogFile(home)
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file [%s], %s", logFile, err)
	}
	fmt.Printf("Pillar log file: %s\n\n", logFile)

	Logger = log.New(file, "Pillar: ", log.LstdFlags|log.Llongfile|log.Ldate|log.Ltime)
}

func isDir(path string) bool {
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return true
	}

	return false
}

func exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}

	return false
}

func getLogFile(pillarHome string) string {
	if pillarHome == "" {
		pillarHome = strings.TrimSuffix(os.TempDir(), "/") + "/pillar"
	}

	logPath := pillarHome + "/logs"
	if !exists(logPath) {
		os.MkdirAll(logPath, 0700)
	}

	return logPath + "/pillar.log"
}
