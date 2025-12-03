package cache

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/configuration/environment"
	"github.com/petshop-system/petshop-api/configuration/repository"
	"go.uber.org/zap"
)

const (
	ErrorToInsertValueInRedis = "Failed to insert value in Redis"
	ErrorToGetInRedis         = "Failed to get value from Redis"
	ErrorToDeleteInRedis      = "Failed to delete value in Redis"
)

type Redis struct {
	RedisClient *redis.Client
	LoggerSugar *zap.SugaredLogger
}

func NewRedis(loggerSugar *zap.SugaredLogger) Redis {

	redisClient := repository.NewRedisClient(environment.Setting.Redis.Addr, environment.Setting.Redis.DB,
		environment.Setting.Redis.Password, environment.Setting.Redis.PoolSize,
		environment.Setting.Redis.ReadTimeout, loggerSugar)

	return Redis{
		RedisClient: redisClient,
		LoggerSugar: loggerSugar,
	}
}

func (r *Redis) Set(ctx domain.ContextControl, key string, payload string, expirationTime time.Duration) error {

	if _, err := r.RedisClient.Set(ctx.Context, key, payload, expirationTime).Result(); err != nil {
		r.LoggerSugar.Errorw(ErrorToInsertValueInRedis, "err", err.Error())
		return err
	}

	return nil
}

func (r *Redis) Get(ctx domain.ContextControl, key string) (string, error) {

	value, err := r.RedisClient.Get(ctx.Context, key).Result()
	if err != nil {
		r.LoggerSugar.Warnw(ErrorToGetInRedis, "err", err.Error())
		return "", err
	}

	return value, err
}

func (r *Redis) Delete(ctx domain.ContextControl, key string) error {

	if _, err := r.RedisClient.Del(ctx.Context, key).Result(); err != nil {
		r.LoggerSugar.Warnw(ErrorToDeleteInRedis, "err", err.Error())
		return err
	}

	return nil
}
