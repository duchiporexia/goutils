package testutils

import (
	_ "github.com/duchiporexia/goutils/xutils/testutils/autoload"
	"gorm.io/driver/mysql"
)

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-txdb"
	"github.com/duchiporexia/goutils/xlog"
	"gorm.io/gorm"
)

func init() {
	fmt.Printf("In Testing !!!!!!\n")
}

func RegisterDBForTest() {
	txdb.Register("txdb", "mysql", fmt.Sprintf("root:rootpwd@tcp(127.0.0.1:3306)/shub_dev?charset=utf8mb4&parseTime=True"))
}

func InitIdGeneratorTable(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE db_id_generator")
	if err != nil {
		return err
	}
	_, err = db.Query(`INSERT INTO db_id_generator 
    (id, version, name, create_time, expire_time) values 
     (1, 1, 'name1', '2020-01-01','2020-01-01' ),
     (2, 1, 'name2', '2020-01-01','2020-01-01' )`)
	return err
}

func NewGormDBForTesting() *gorm.DB {
	dbg, err := gorm.Open(mysql.New(mysql.Config{DriverName: "txdb"}), &gorm.Config{
		Logger: xlog.NewGormLogger(),
	})
	if err != nil {
		panic(err)
	}
	return dbg
}
