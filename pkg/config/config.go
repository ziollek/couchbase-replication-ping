package config

import (
	"github.com/spf13/viper"
	"time"
)

type Couchbase struct {
	URI      string `yaml:"uri"`
	Bucket   string `yaml:"bucket"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Generator struct {
	TTL  time.Duration `yaml:"ttl"`
	Size int           `yaml:"size"`
}

type Config struct {
	Source      *Couchbase `yaml:"source"`
	Destination *Couchbase `yaml:"destination"`
	Generator   *Generator `yaml:"generator"`
	Key         string     `yaml:"key"`
}

func GetConfig() (*Config, error) {
	var config Config
	err := viper.Unmarshal(&config)
	return &config, err
}
