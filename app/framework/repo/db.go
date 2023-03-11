package repo

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"

	"go.uber.org/zap"
)

func NewRedis(c DBConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: "",
		DB:       0,
	})

	for i := 0; i < 10; i++ {
		err := client.Ping().Err()
		if err == nil {
			return client, nil
		}
		zap.S().Info("retrying db ping: ", err)
		time.Sleep(time.Second)
	}

	return nil, fmt.Errorf("could not connect to redis")
}

func NewDBConfig() DBConfig {
	return DBConfig{}
}

type DBConfig struct {
	Host string `env:"REDIS_HOST"`
	Port int    `env:"REDIS_PORT"`
}
