package config

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	err := godotenv.Load(".env")
	return err
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
