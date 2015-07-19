package gsonapi

import (
	"log"

	"github.com/joho/godotenv"
)

var GsonApiConfig Config

func init() {
	env, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading GsonApi .env file")
	}

	// map env values to config
	url := env["GSON_API_URL"]

	// validate
	if url == "" {
		log.Fatal("the following env value cannot be blank: GSON_API_URL")
	}

	// create config
	GsonApiConfig = Config{URL: url}
}

type Config struct {
	URL string
}
