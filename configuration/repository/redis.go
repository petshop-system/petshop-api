package repository

import (
	//glog "bitbucket.org/maironmscosta/golang-log/v1"
	"context"
	redis "github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"time"
)

var ctx = context.Background()

func NewRedisClient(addr string, db int, password string, poolSize int, readTimeout time.Duration, loggerSugar *zap.SugaredLogger) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password, // no password set
		DB:          db,       // use default DB
		PoolSize:    poolSize,
		ReadTimeout: readTimeout,
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		loggerSugar.Errorw("error to start redis", "pong", pong, "err", err.Error(),
			"Addr", addr)
		panic(err)
	}

	return rdb
}
