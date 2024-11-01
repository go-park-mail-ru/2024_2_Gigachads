package redis

import (
	"fmt"
	"time"
	"context"
	"mail/config"
	"github.com/redis/go-redis/v9"
)

func Init(cfg *config.Config) (*redis.Client, error) {
	cfgRedis := cfg.Redis
	redisAddress := fmt.Sprintf("%s:%s", cfgRedis.IP, cfgRedis.Port)
	client := redis.NewClient(&redis.Options{
		DB:   cfgRedis.DBnum,
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