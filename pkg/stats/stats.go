package stats

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	"github.com/pborman/uuid"

	"github.com/coralproject/pillar/pkg/backend"
)

const (
	cfgLoggingLevel = "LOGGING_LEVEL"
	// VersionNumber is the version for sponge
	VersionNumber = 0.1
)

var (
	uid string
)

var (
	config struct {
		mongodb struct {
			addrs, database, username, password, passwordFile string
			ssl                                               bool
		}
	}
)

// Init initialize log and get configuration
func Init() {

	// Initialize logging
	logLevel := func() int {
		ll, err := cfg.Int(cfgLoggingLevel)
		if err != nil {
			return log.USER
		}
		return ll
	}
	log.Init(os.Stderr, logLevel, log.Ldefault)

	// Generate UUID to use with the logs
	uid = uuid.New()

	config.mongodb.addrs = os.Getenv("MONGODB_ADDRESS")
	config.mongodb.username = os.Getenv("MONGODB_USERNAME")
	config.mongodb.password = os.Getenv("MONGODB_PASSWORD")
	config.mongodb.passwordFile = os.Getenv("MONGODB_PASSWORD_FILE")
	config.mongodb.database = os.Getenv("MONGODB_DATABASE")
	config.mongodb.ssl, _ = strconv.ParseBool(os.Getenv("MONGODB_SSL"))

	if config.mongodb.addrs == "" {
		log.Fatal(uid, "Init", "Remember to set up environment variables to connect to MongoDB server.")
	}

	// Check if a password file was provided.
	if config.mongodb.passwordFile != "" {
		log.User(uid, "Init", "Reading MongoDB password from %s", config.mongodb.passwordFile)
		passwordBytes, err := ioutil.ReadFile(config.mongodb.passwordFile)
		if err != nil {
			log.Error(uid, "Init", err, "Error reading password file.")
		}

		config.mongodb.password = string(passwordBytes)
	}

	// Set a the environment variable, MONGODB_PASSWORD, as a default value
	// for password.
	if config.mongodb.password == "" {
		config.mongodb.password = os.Getenv("MONGODB_PASSWORD")
		if config.mongodb.username != "" && config.mongodb.password == "" {
			log.User(uid, "Init", "Warning: a username is in use without a password")
		}
	}

	if config.mongodb.ssl {
		log.User(uid, "Init", "Using TLS for MongoDB connections")
	}
}

// Calculate does all the calculation
func Calculate() {

	log.User(uid, "main", "Connecting to MongoDB at %s", config.mongodb.database)
	b, err := backend.NewMongoDBBackend(strings.Split(config.mongodb.addrs, ","), config.mongodb.database, config.mongodb.username, config.mongodb.password, config.mongodb.ssl)
	if err != nil {
		log.Fatal(uid, "main", "Failed on connecting with MongoDB. Error: %v", err)
	}

	defer b.Close()

	log.User(uid, "main", "Calculating stats")
	ctx := backend.NewBackendContext(context.Background(), backend.NewIdentityMap(b))
	if err := CalculateUserStatistics(ctx); err != nil {
		log.Error(uid, "main", err, "Failed on calculating stats.")
	}

}
