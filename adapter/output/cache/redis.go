package cache

import (
	"github.com/go-redis/redis/v8"
	"github.com/petshop-system/petshop-api/application/domain"
	"github.com/petshop-system/petshop-api/configuration/environment"
	"github.com/petshop-system/petshop-api/configuration/repository"
	"go.uber.org/zap"
	"time"
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
		r.LoggerSugar.Errorw("It was not possible insert value in redis", "err", err.Error())
		return err
	}

	return nil
}

func (r *Redis) Get(ctx domain.ContextControl, key string) (string, error) {

	value, err := r.RedisClient.Get(ctx.Context, key).Result()
	if err != nil {
		r.LoggerSugar.Warnw("It was not possible get value in redis", "err", err.Error())
		return "", err
	}

	return value, err
}

func (r *Redis) Delete(ctx domain.ContextControl, key string) error {

	if _, err := r.RedisClient.Del(ctx.Context, key).Result(); err != nil {
		r.LoggerSugar.Warnw("It was not possible delete value in redis", "err", err.Error())
		return err
	}

	return nil
}
