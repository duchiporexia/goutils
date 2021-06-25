package xgorm

import (
	"errors"
	"github.com/duchiporexia/goutils/xerr"
	"github.com/duchiporexia/goutils/xlog"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBConfig struct {
	DbType string `yaml:"dbType" env:"DB_TYPE" env-default:"postgres"`
	DbUrl  string `yaml:"dbUrl" env:"DB_URL" env-default:"postgres://postgres:postgrespwd@localhost:5432/app_dev?sslmode=disable"`
}

var DB *gorm.DB

func Init(cfg *DBConfig) {
	if cfg.DbUrl == "" {
		xlog.Fatal("db url is empty")
	}
	var dialector gorm.Dialector
	switch cfg.DbType {
	case "postgres":
		dialector = postgres.Open(cfg.DbUrl)
	case "mysql":
		dialector = mysql.Open(cfg.DbUrl)
	default:
		xlog.Fatal("unknown db type:" + cfg.DbType)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: xlog.NewGormLogger(),
	})
	if err != nil {
		panic(err)
	}
	DB = db
}

func HandleNoRows(db *gorm.DB) *gorm.DB {
	if db.Error != nil && errors.Is(db.Error, gorm.ErrRecordNotFound) {
		db.Error = xerr.ErrNoRows
	}
	return db
}
