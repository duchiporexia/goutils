package tests

import (
	. "github.com/duchiporexia/goutils/xmsg"
	"time"
)

//easyjson -all -omit_empty ./json_data.go

type Teacher struct {
	Id         LongID     `json:"id"`
	Name       string     `json:"name"`
	Age        int        `json:"age"`
	CreateTime *time.Time `json:"createTime"`
	UpdateTime *DateTime  `json:"updateTime"`
	Grade      Grade      `json:"grade_info"`
	Students   []Student  `json:"students"`
}

type Grade struct {
	School string `json:"school_name"`
	Grade  int    `json:"grade"`
}

type Student struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  uint8  `json:"sex"`
}
