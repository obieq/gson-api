package gsonapi

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config => global variable that stores the config values
var Config *config

func init() {
	Config = newConfig()
}

// config => stores a map[string]string of all values parse from an .env file
type config struct {
	URL string
}

func newConfig() *config {
	c := &config{}
	c.ParseConfigFile("config")
	c.Validate()

	return c
}

// ParseConfigFile => reads the config values from the config.json file
func (c *config) ParseConfigFile(configFileName string) (err error) {
	viper.Reset()
	viper.SetConfigName(configFileName)
	if err = viper.ReadInConfig(); err != nil {
		log.Panicln("Fatal error reading gson-api viper config file:", err)
	}

	// get api base url
	c.URL = getString("gson_api_url")

	return err
}

func (c *config) Validate() {
	if c.URL == "" {
		log.Panicln("gson api config error: URL cannot be blank")
	}
}

func getString(propKey string) string {
	v := viper.GetString(propKey)
	if isEnv(v) {
		return getEnv(v[4 : len(v)-1])
	}
	return v
}

func isEnv(configValue string) bool {
	return strings.HasPrefix(configValue, "ENV[")
}

func getEnv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		log.Panicln("env value cannot be blank:", name)
	}
	return v
}
