package config

import (
	"log"
	"os"
)

var (
	config *Config
)

//Config encapsulates application specific configuration
type Config struct {
	Home     string
	MongoURL string
	Logger   *log.Logger
	Address  string
}

//Address returns the address where the app is running
func Address() string {
	return config.Address
}

//Logger returns the logger for this app
func Logger() *log.Logger {
	return config.Logger
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

	logger := log.New(os.Stdout, "Pillar: ", log.LstdFlags|log.Llongfile|log.Ldate|log.Ltime)
	config = &Config{Home: home, Address: address, MongoURL: url, Logger: logger}
}
