package types

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type RedisSession struct {
	Client *redis.Client
	Ctx    context.Context
}
