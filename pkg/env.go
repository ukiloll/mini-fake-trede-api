package pkg

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}
	return nil
}