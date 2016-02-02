package config

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	LogFile string = "pillar.log"
)

var (
	Logger *log.Logger
	config *Config
)

type Config struct {
	Home     string
	LogFile  string
	MongoURL string
	Address  string
}

func init() {

	home := os.Getenv("PILLAR_HOME")
	if home == "" {
		log.Fatal("Error initializing Pillar: PILLAR_HOME not found.")
	}

	address := os.Getenv("PILLAR_ADDRESS")
	if address == "" {
		log.Fatal("Error initializing Pillar: PILLAR_ADDRESS not found.")
	}

	url := os.Getenv("MONGODB_URL")
	if url == "" {
		log.Fatal("Error initializing Pillar: MONGODB_URL not found.")
	}

	logFile := getLogFile(home)
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file [%s], %s", logFile, err)
	}
	fmt.Printf("Pillar log file: %s\n\n", logFile)

	Logger = log.New(file, "Pillar: ", log.LstdFlags|log.Llongfile|log.Ldate|log.Ltime)
	config = &Config{Home: home, Address: address, MongoURL: url, LogFile: logFile}
}

func getLogFile(pillarHome string) string {
	if pillarHome == "" {
		pillarHome = strings.TrimSuffix(os.TempDir(), "/") + "/pillar"
	}

	logPath := pillarHome + "/logs"
	if !exists(logPath) {
		os.MkdirAll(logPath, 0700)
	}

	return logPath + "/" + LogFile
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
