package xconfig

import (
	"fmt"
	"testing"
)

type databaseConfig struct {
	Name     string `yaml:"name" env:"NAME" env-default:"app_dev"`
	User     string `yaml:"user" env:"USER" env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"postgrespwd"`
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
}

type appConfig struct {
	Test1 string `yaml:"test1" env:"test1" env-default:"testvalue1"`
	Test2 string `yaml:"test2" env:"test2" env-default:"testvalue2"`

	Env      string `yaml:"env" env:"APP_ENV" env-default:"dev"`
	ConfPath string `yaml:"confPath" env:"APP_CONF_PATH" env-default:"conf"`

	Database databaseConfig `yaml:"database" env-prefix:"POSTGRES_"`
}

var cfg appConfig

func TestCleanEnvLoad(t *testing.T) {
	LoadConfig("", &cfg)
	fmt.Printf("%v\n", cfg)
}
