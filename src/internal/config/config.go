package config

import (
	"flag"
	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Listen struct {
		Port string `yaml:"port" env-default:"8000"`
	} `yaml:"listen"`

	Database struct {
		URL string `env:"url" env-required:"true"`
	} `yaml:"database"`

	Hash struct {
		Salt string `yaml:"salt" env-required:"true"`
	} `yaml:"hash"`

	Auth struct {
		SecretKey string `yaml:"secret_key" env-required:"true"`
	} `yaml:"auth"`

	CORS struct {
		Origins []string `yaml:"origins"`
	} `yaml:"cors"`
}

var cfg *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		cfg = &Config{}

		path := getConfigPath()

		if err := cleanenv.ReadConfig(path, cfg); err != nil {
			panic(err)
		}
	})

	return cfg
}

func getConfigPath() string {
	var path string
	flag.StringVar(&path, "config", "../config/default.yaml", "set config file")

	envPath := os.Getenv("CONFIG_PATH")

	if len(envPath) > 0 {
		path = envPath
	}

	return path
}
