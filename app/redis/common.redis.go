package redis

import (
	"context"
	"course-work/app/types"
	"course-work/app/vocabulary"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var session *types.RedisSession

func Init() error {
	var ctx context.Context = context.Background()

	url := os.Getenv("REDIS_URL")

	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf(vocabulary.REDIS_CONNECTION_PROBLEM, err.Error())
	}

	session = &types.RedisSession{
		Ctx:    ctx,
		Client: client,
	}

	return nil
}

func Close() {
	session.Client.Close()
}
