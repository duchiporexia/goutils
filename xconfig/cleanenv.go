package xconfig

import (
	"fmt"
	"github.com/duchiporexia/goutils/xutils"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

// {APP_CONF_PATH}/app_{env}.yaml
func getEnvFileName() string {
	env := Get("APP_ENV", "dev")
	confPath := Get("APP_CONF_PATH", "conf")
	return confPath + "/app_" + env + ".yaml"
}

func LoadConfig(fileName string, cfg interface{}) {
	if fileName == "" {
		fileName = getEnvFileName()
	}
	if xutils.FileExists(fileName) {
		err := cleanenv.ReadConfig(fileName, cfg)
		if err != nil {
			log.Fatal(fmt.Sprintf("Error loading %v file: %v", fileName, err))
		}
	} else {
		log.Printf(fmt.Sprintf("No env file: %v", fileName))
	}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatal(fmt.Sprintf("Error ReadEnv %v ", err))
	}
}
