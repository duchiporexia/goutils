package tests

import (
	"github.com/duchiporexia/goutils/xmsg"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestMsgPDemo(t *testing.T) {
	foo := &Foo{
		Id:         xmsg.LongID(9223372036854775807),
		Name:       "xvv",
		Age:        18,
		Map:        map[string]string{"k1": "v1"},
		Map2:       map[string]Dog{"d1": Dog{Name: "dog1", Age: 2}},
		Arr:        []string{"k1", "k2"},
		Arr2:       []Dog{Dog{Name: "dog1", Age: 2}, Dog{Name: "dog2", Age: 22}},
		Dog:        Dog{Name: "dog1", Age: 2},
		Dog2:       &Dog{Name: "dog12", Age: 22},
		CreateTime: xmsg.DateTime(time.Now()),
	}
	fooBytes, _ := foo.MarshalMsg(nil)
	decodedFoo := &Foo{}
	decodedFoo.UnmarshalMsg(fooBytes)

	assert.Equal(t, decodedFoo.Id, foo.Id)
	assert.Equal(t, decodedFoo.Name, foo.Name)
	assert.Equal(t, decodedFoo.Age, foo.Age)
	assert.True(t, reflect.DeepEqual(decodedFoo.Map, foo.Map))
	assert.True(t, reflect.DeepEqual(decodedFoo.Map2, foo.Map2))
	assert.True(t, reflect.DeepEqual(decodedFoo.Arr, foo.Arr))
	assert.True(t, reflect.DeepEqual(decodedFoo.Arr2, foo.Arr2))
	assert.True(t, reflect.DeepEqual(decodedFoo.Dog, foo.Dog))
	assert.True(t, reflect.DeepEqual(decodedFoo.Dog2, foo.Dog2))
	assert.Equal(t, decodedFoo.CreateTime.UnixNano(), foo.CreateTime.UnixNano())
}

func TestMsgPDemoTeachers(t *testing.T) {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2014-11-12T11:45:26.371Z"
	tx, _ := time.Parse(layout, str)
	teacher := MTeacher{
		Id:          xmsg.LongID(33),
		Name:        "Tom",
		Age:         30,
		CreateTime:  xmsg.DateTime(tx),
		Update1Time: xmsg.DateTime(tx),
		Update2Time: xmsg.DateTime(tx),
		Update3Time: xmsg.DateTime(tx),
		Update4Time: xmsg.DateTime(tx),
		//Update5Time: tx,
		//Update6Time: tx,
		//Update7Time: tx,
		//Update8Time: tx,
		Grade: MGrade{
			School: "xxx school",
			Grade:  1,
		},
		Students: []MStudent{
			MStudent{Name: "Jake", Age: 8, Sex: 0},
			MStudent{Name: "Lindy", Age: 8, Sex: 0},
			MStudent{Name: "LiLei", Age: 8, Sex: 0},
		},
	}
	//fooBytes, _ := teacher.MarshalMsg(nil)
	//for i := 0; i < 5000000; i++ {
	//	decoded := &MTeacher{}
	//	decoded.UnmarshalMsg(fooBytes)
	//}
	//
	for i := 0; i < 5000000; i++ {
		teacher.MarshalMsg(nil)
	}

	//fmt.Println(decoded)

}
