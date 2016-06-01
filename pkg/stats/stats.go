package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/pborman/uuid"
	"golang.org/x/net/context"

	ca "github.com/coralproject/pillar/app/stats/calc"
	"github.com/coralproject/pillar/pkg/backend"
	"github.com/coralproject/pillar/pkg/backend/mongodb"
)

const (
	cfgLoggingLevel = "LOGGING_LEVEL"
)

var (
	config struct {
		mongodb struct {
			addrs, database, username, password, passwordFile string
			ssl                                               bool
		}
	}

	uid string
)

func init() {

	// Initialize logging
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel, log.USER)

	// Generate UUID to use with the logs
	uid = uuid.New()

	// Flag information and defaults.
	flag.StringVar(&config.mongodb.addrs, "mongodb-address", "127.0.0.1:27017", "comma-seperated list of mongodb host:port pairs")
	flag.StringVar(&config.mongodb.username, "mongodb-username", "", "mongodb username")
	flag.StringVar(&config.mongodb.password, "mongodb-password", "", "mongodb password (defaults to MONGODB_PASSWORD)")
	flag.StringVar(&config.mongodb.passwordFile, "mongodb-password-file", "", "mongodb password file")
	flag.StringVar(&config.mongodb.database, "mongodb-database", "coral", "mongodb database")
	flag.BoolVar(&config.mongodb.ssl, "mongodb-ssl", false, "use TLS for mongodb connections")
}

func main() {
	flag.Parse()

	// Parse the MongoDB address list.
	addrs := strings.Split(config.mongodb.addrs, ",")

	// Check if a password file was provided.
	if config.mongodb.passwordFile != "" {
		log.User(uid, "main", "Reading MongoDB password from %s", config.mongodb.passwordFile)
		passwordBytes, err := ioutil.ReadFile(config.mongodb.passwordFile)
		if err != nil {
			log.Error(uid, "main", err, "Unable to read password's file.")
		}

		config.mongodb.password = string(passwordBytes)
	}

	// Set a the environment variable, MONGODB_PASSWORD, as a default value
	// for password.
	if config.mongodb.password == "" {
		config.mongodb.password = os.Getenv("MONGODB_PASSWORD")
		if config.mongodb.username != "" && config.mongodb.password == "" {
			log.Dev(uid, "main", "Warning: a username is in use without a password")
		}
	}

	if config.mongodb.ssl {
		log.User(uid, "main", "Using TLS for MongoDB connections")
	}

	log.User(uid, "main", "Connecting to MongoDB at %s/%s", addrs, config.mongodb.database)
	b, err := mongodb.NewMongoDBBackend(addrs, config.mongodb.database, config.mongodb.username, config.mongodb.password, config.mongodb.ssl)
	if err != nil {
		log.Error(uid, "main", err, "Failed on connecting to MongoDB.")
	}

	log.User(uid, "main", "Calculating stats")
	ctx := backend.NewBackendContext(context.Background(), backend.NewIdentityMap(b))
	if err := ca.CalculateUserStatistics(ctx); err != nil {
		log.Error(uid, "main", err, "Failed on calculating stats.")
	}
}
