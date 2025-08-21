package envi

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetEnvVariable(envName string) string {

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	returnedEnv := os.Getenv(envName)

	return returnedEnv
}
