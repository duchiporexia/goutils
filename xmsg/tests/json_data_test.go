package tests

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"
	"time"
)

//"createTime": "2021-02-08T12:03:20.598Z",
//"updateTime": 1612785800598,
var jsonString = `{
	"id": "",
	"name": "Tom",
	"age": 30,
	"createTime": "2021-02-08T12:03:20.598Z",
	"updateTime": 1612785800598,
	"grade_info": {
		"school_name": "xxx school",
		"grade": 1
	},
	"students": [
		{"name": "Jake", "age": 8, "sex": 0},
		{"name": "Lindy", "age": 9, "sex": 1},
		{"name": "LiLei", "age": 7, "sex": 0}
	]
}`

func TestGrade_EasyJSON(t *testing.T) {
	runtime.GOMAXPROCS(1)

	jsonB := []byte(jsonString)

	//// 500000
	//for i := 0; i < 670000; i++ {
	//	teacher := new(Teacher)
	//	err := teacher.UnmarshalJSON(jsonB)
	//	if err != nil {
	//		fmt.Println("err UnmarshalJSON:", err)
	//	} else {
	//		//fmt.Println("teacher:", teacher)
	//	}
	//	time.Unix(1612783727, 712000)
	//}
	teacher := new(Teacher)
	teacher.UnmarshalJSON(jsonB)

	ret := make([]byte, 0, 10000)
	// 820000
	for i := 0; i < 1; i++ {
		ret = ret[:0]
		jsonStr, err := teacher.MarshalJSON()
		if err != nil {
			fmt.Println("err MarshalJSON:", err)
		} else {
			ret = append(ret, jsonStr...)
			fmt.Println("jsonStr:", string(jsonStr))
		}
	}

}

func TestStandardJSON(t *testing.T) {
	jsonB := []byte(jsonString)
	for i := 0; i < 660000; i++ {
		teacher := new(Teacher)
		err := json.Unmarshal(jsonB, teacher)
		if err != nil {
			fmt.Println("err UnmarshalJSON:", err)
		} else {
			//fmt.Println("teacher:", teacher)
		}
		time.Unix(1612783727, 712000)
	}
}
