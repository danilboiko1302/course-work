package main

import (
	"course-work/app/controller"
	"course-work/app/nats"
	"course-work/app/redis"
	"math/rand"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	rand.Seed(time.Now().UnixNano())

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

	if err := controller.GetRoute().Run(":8000"); err != nil {
		panic(err)
	}
}
