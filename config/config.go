package config

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"

	"github.com/robvdl/gcms/util"
)

// AppConfig struct is for storing application configuration.
// Note that the underscores are only needed due to a envconfig issue
// see bug: https://github.com/kelseyhightower/envconfig/issues/26
type AppConfig struct {
	// general app settings
	Debug bool   `default:"true"`
	Port  string `default:"8080"`

	// database connection string
	Database_URL string

	// cache and session
	Cache_URL      string
	Session_Secret string

	// password encryption settings
	Password_Algorithm  string `default:"pbkdf2-sha256"`
	Password_Iterations int    `default:"12000"` // used by pbkdf2
	Password_Salt_Size  int    `default:"6"`     // salt size used by pbkdf2
	Password_Cost       int    `default:"10"`    // cost is used by bcrypt
}

// Config stores the global application configuration instance
var Config AppConfig

// LoadAppConfig will try to load the config file /etc/default/appname first
// if that exists, otherwise it will try .env in the current directory.
// If neither was found we rely entiry on environment variables (12-factor).
func LoadAppConfig(project string) {
	loadEnvConfig("/etc/default/"+project, ".env")

	// envconfig then loads environment variables into the Config struct
	err := envconfig.Process(project, &Config)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Set Gin release or debug mode based on Config.Debug
	if Config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

// loadEnvConfig takes a list of config file paths, it will try to load
// the first one it can find and returns either true if a config file was
// loaded or false if none of the config files existed.  It will only ever
// load one config file and return.  A config file that cannot be read
// properly will throw an error.
func loadEnvConfig(filenames ...string) bool {
	for _, filename := range filenames {
		if util.Exists(filename) {
			err := godotenv.Load(filename)

			// if the config file cannot be read we want to know about it
			if err != nil {
				log.Fatal(err.Error())
			} else {
				return true
			}
		}
	}
	return false
}
