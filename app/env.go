package app

import (
	env "github.com/joho/godotenv"
)

func LoadEnv() {
	err := env.Load(".env")
	if err != nil {
		panic(err)
	}
}
