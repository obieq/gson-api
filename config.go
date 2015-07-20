package gsonapi

import (
	"log"

	"github.com/joho/godotenv"
)

// GsonAPIConfig => global variable that stores the config values
var GsonAPIConfig Config

func init() {
	GsonAPIConfig = Config{}
	// parse .env file
	GsonAPIConfig.ParseEnvs(nil)
}

// Config => stores a map[string]string of all values parse from an .env file
type Config struct {
	Envs map[string]string
}

// ConfigTester => sole use is to increase test coverage
type ConfigTester struct {
	Error error
	Envs  map[string]string
}

// URL => returns the url config value
func (c *Config) URL() string {
	// validate
	var url = c.Envs["GSON_API_URL"]
	if url == "" {
		log.Panic("the following env value cannot be blank: GSON_API_URL")
	}

	return url
}

// ParseEnvs => reads the config values from a .env file
func (c *Config) ParseEnvs(ct *ConfigTester) (err error) {
	var envs map[string]string

	if ct == nil {
		envs, err = godotenv.Read()
	} else {
		envs = ct.Envs
		err = ct.Error
	}

	if err != nil {
		// log.Fatal("Error loading GsonAPI .env file")
		log.Panicln("Error loading GsonAPI .env file:", err)
	}

	// set envs
	c.Envs = envs

	// verify required config values were supplied via the .env file
	c.URL()

	return err
}
