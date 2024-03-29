package environment

import "time"

type setting struct {
	Application struct {
		ContextRequest time.Duration `envconfig:"CONTEXT_REQUEST" default:"2.1s"`
	}

	Server struct {
		Context      string        `envconfig:"SERVER_CONTEXT" default:"petshop-api"`
		Port         string        `envconfig:"PORT" default:"5001" required:"true" ignored:"false"`
		ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`
		WriteTimeout time.Duration `envconfig:"READ_TIMEOUT" default:"10s"`
	}

	Redis struct {
		Addr        string        `envconfig:"REDIS_ADDR" default:"localhost:6379"`
		Password    string        `envconfig:"REDIS_PASSWORD"`
		DB          int           `envconfig:"REDIS_DB" default:"0"`
		PoolSize    int           `envconfig:"POOL_SIZE" default:"100"`
		ReadTimeout time.Duration `envconfig:"READ_TIMEOUT" default:"2s"`
	}

	Postgres struct {
		DBUser     string `envconfig:"DB_USER" default:"petshop-system"`
		DBPassword string `envconfig:"DB_PASSWORD" default:"test1234"`
		DBName     string `envconfig:"DB_NAME" default:"petshop-system"`
		DBHost     string `envconfig:"DB_HOST" default:"localhost"`
		DBPort     string `envconfig:"DB_PORT" default:"5432"`
		DBType     string `envconfig:"DB_TYPE" default:"postgres"`
	}

	Kafka struct {
		Schedule struct {
			BootstrapServer string `envconfig:"KAFKA_SCHEDULE_BOOTSTRAP_SERVER" default:"localhost:29092"`
			GroupID         string `envconfig:"KAFKA_SCHEDULE_GROUPID" default:"kafka_schedule"`
			AutoOffsetReset string `envconfig:"KAFKA_SCHEDULE_AUTO_OFFSET_RESET" default:"earliest"`
			Topic           string `envconfig:"KAFKA_SCHEDULE_TOPIC" default:"schedule"`
		}
	}
}

var Setting setting
