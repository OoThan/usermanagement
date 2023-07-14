package ds

import (
	"context"
	"net"
	"time"

	"github.com/OoThan/usermanagement/config"
	"github.com/OoThan/usermanagement/pkg/logger"
	"github.com/go-redis/redis/v8"
)

func LoadRDB() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(config.Redis().Host, config.Redis().Port),
		Username: config.Redis().Username,
		Password: config.Redis().Password,
		DB:       0,
	})

	logger.Sugar.Info("Successfully connected to redis")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		logger.Sugar.Error(err)
		panic(err)
	}

	return rdb, nil
}
