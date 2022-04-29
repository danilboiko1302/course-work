package main

import (
	"course-work/app/controller"
	"course-work/app/nats"

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

	// if err := redis.Init(); err != nil {
	// 	panic(err)
	// }

	// defer redis.Close()

	if err := controller.GetRoute().Run(":8000"); err != nil {
		panic(err)
	}
}
