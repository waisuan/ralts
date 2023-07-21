package config

import (
	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
)

type Config struct {
	DatabaseConn     string `env:"DATABASE_URL,notEmpty"`
	AuthToken        string `env:"AUTH_TOKEN" envDefault:"token"`
	MaxConnCount     int    `env:"MAX_CONN_COUNT" envDefault:"50"`
	RedisConn        string `env:"REDIS_URL,notEmpty"`
	MaxSentMsgPerDay int    `env:"MAX_SENT_MSG_PER_DAY" envDefault:"1000"`
	ServiceUsername  string `env:"SERVICE_USERNAME"`
	ServicePassword  string `env:"SERVICE_PASSWORD"`
	AuthEnabled      bool   `env:"AUTH_ENABLED" envDefault:"true"`
}

const projectDirName = "ralts"

func NewConfig(testing bool) *Config {
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	var err error
	if testing {
		err = godotenv.Load(string(rootPath) + `/.env.test`)
	} else {
		err = godotenv.Load(string(rootPath) + `/.env`)
	}
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	cfg := Config{}
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}

	return &cfg
}
