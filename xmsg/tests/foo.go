package tests

import (
	. "github.com/duchiporexia/goutils/xmsg"
)

//go:generate msgp --tests=false

type Foo struct {
	Id         LongID            `msg:"id"`
	Name       string            `msg:"name"`
	Age        int               `msg:"age"`
	Map        map[string]string `msg:"map"`
	Map2       map[string]Dog    `msg:"map2"`
	Arr        []string          `msg:"arr"`
	Arr2       []Dog             `msg:"arr2"`
	Dog        Dog               `msg:"dog"`
	Dog2       *Dog              `msg:"dog2"`
	CreateTime DateTime          `msg:"createTime,extension"`
	UpdateTime *DateTime         `msg:"updateTime,extension"`
}

type Dog struct {
	Name string `msg:"name"`
	Age  int    `msg:"age"`
}
type MTeacher struct {
	Id          LongID   `msg:"id"`
	Name        string   `msg:"name"`
	Age         int      `msg:"age"`
	CreateTime  DateTime `msg:"createTime,extension"`
	Update1Time DateTime `msg:"update1Time,extension"`
	Update2Time DateTime `msg:"update2Time,extension"`
	Update3Time DateTime `msg:"update3Time,extension"`
	Update4Time DateTime `msg:"update4Time,extension"`
	//Update5Time time.Time   `msg:"update5Time"`
	//Update6Time time.Time   `msg:"update6Time"`
	//Update7Time time.Time   `msg:"update7Time"`
	//Update8Time time.Time   `msg:"update8Time"`
	Grade    MGrade     `msg:"grade_info"`
	Students []MStudent `msg:"students"`
}

type MGrade struct {
	School string `msg:"school_name"`
	Grade  int    `msg:"grade"`
}

type MStudent struct {
	Name string `msg:"name"`
	Age  int    `msg:"age"`
	Sex  uint8  `msg:"sex"`
}
