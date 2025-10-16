package config

import "github.com/joho/godotenv"

func LoadEnv(filenames ...string) {
	if err := godotenv.Load(filenames...); err != nil {
		panic(err)
	}
}
