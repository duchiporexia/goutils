package autoload

import (
	"github.com/duchiporexia/goutils/xutils"
	"os"
)

func init() {
	os.Setenv("APP_CONF_PATH", xutils.AbsPath("common/test"))
	os.Setenv("APP_ENV", "test")
}
