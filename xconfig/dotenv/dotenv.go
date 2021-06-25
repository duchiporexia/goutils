package dotenv

import (
	"fmt"
	"github.com/duchiporexia/goutils/xlog"
	"github.com/duchiporexia/goutils/xutils"
	"github.com/joho/godotenv"
	"os"
)

func init() {
	env := os.Getenv("APP_ENV")
	if "" == env {
		env = "dev"
	}

	envDirPath := os.Getenv("APP_CONF_PATH")
	if envDirPath == "" {
		envDirPath = "conf"
	}

	envFileName := envDirPath + "/.env." + env

	if xutils.FileExists(envFileName) {
		err := godotenv.Load(envFileName)
		if err != nil {
			xlog.Fatal(fmt.Sprintf("Error loading %v file", envFileName))
		}
	} else {
		xlog.Warn(fmt.Sprintf("No env file: %v", envFileName))
	}

	err := godotenv.Load(envDirPath + "/.env")
	if err != nil {
		xlog.Fatal("Error loading .env file")
	}
}
