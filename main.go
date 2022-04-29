package main

import (
	"course-work/app/nats"
	"course-work/app/redis"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	if err := nats.Init(); err != nil {
		panic(err)
	}

	defer nats.Close()

	if err := redis.Init(); err != nil {
		panic(err)
	}

	defer redis.Close()
}
