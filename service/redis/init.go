package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"mail/config"
	"time"
)

func Init(cfg *config.Config, DBnum int) (*redis.Client, error) {
	cfgRedis := cfg.Redis
	redisAddress := fmt.Sprintf("%s:%s", cfgRedis.IP, cfgRedis.Port)
	client := redis.NewClient(&redis.Options{
		DB:   DBnum,
		Addr: redisAddress,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}
