package tests

import (
	"github.com/duchiporexia/goutils/xmsg"
)

//go:generate go run github.com/duchiporexia/goutils/xgormgen -entity User -entity Dog -output model_gorm_gen.go

////go:generate go run github.com/duchiporexia/goutils/xgormgen -entity User
type User struct {
	Id         xmsg.LongID `gorm:"primary_key"`
	Name       string      `gorm:"column:user_name"`
	Age        int
	CreateTime *xmsg.DateTime
}

type Dog struct {
	Id   xmsg.LongID `gorm:"primary_key"`
	Name string
	Age  int
}
