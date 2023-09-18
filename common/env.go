package common

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Prod Environment = "prod"
	Dev  Environment = "dev"
	Test Environment = "test"
)

func Env(basePath ...string) {
	var currentEnv Environment
	var envFile string

	currentEnv = Environment(os.Getenv("ENV"))

	switch currentEnv {
	case Prod:
		envFile = ".env.prod"
	case Test:
		envFile = ".env.test"
	default:
		envFile = ".env.dev"
	}

	if len(basePath) == 0 {
		basePath = append(basePath, "")
	}

	envFile = path.Join(basePath[0], envFile)

	fmt.Println("Reading env variables from", envFile)

	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("Could not load .env file")
	}
}
