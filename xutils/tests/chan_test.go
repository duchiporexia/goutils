package tests

import (
	"encoding/json"
	"fmt"
	"github.com/duchiporexia/goutils/batch_handler/test/dto"
	"testing"
	"time"
	"unsafe"
)

type Dog struct {
	ID   uint `gorm:"primary_key"`
	Name string
	Age  int
}

func TestChan(t *testing.T) {
	funcDone := make(chan struct{})
	go func() {
		fmt.Printf("sleep 1\n")
		//funcDone <- struct{}{}
		time.Sleep(time.Second)
		close(funcDone)
		fmt.Printf("chan is closed\n")
	}()
	go func() {
		fmt.Printf("sleep ...\n")
		time.Sleep(time.Second * 2)
		fmt.Printf("listen chan\n")
		_, ok := <-funcDone
		fmt.Printf("chan is done %v\n", ok)
	}()
	<-funcDone
	fmt.Printf("program sleep 5\n")
	time.Sleep(time.Second * 5)
	<-funcDone
	fmt.Printf("program is done\n")
}

func getFuncs() []func() (int, int) {
	keys := []int{10, 20, 30, 40, 50}
	var funcs []func() (int, int)
	for i, existingKey := range keys {
		val := existingKey
		funcs = append(funcs, func() (int, int) {
			return i, val
		})
	}
	return funcs
}
func TestClosure(t *testing.T) {
	funcs := getFuncs()
	for _, fun := range funcs {
		i, existingKey := fun()
		fmt.Printf("i:%v, existingKey:%v\n", i, existingKey)
	}
}

type Integer int64

func TestTypes(t *testing.T) {
	nums := []int64{3, 4, 5}
	x, ok := (interface{}(nums)).([]Integer)
	fmt.Printf("x:%v, ok:%v\n", x, ok)
	y, ok := (interface{}(nums)).([]interface{})
	fmt.Printf("y:%v, ok:%v\n", y, ok)
	x = *(*[]Integer)(unsafe.Pointer(&nums))
	fmt.Printf("x:%v\n", x)
	//i := *(*[]interface{})(unsafe.Pointer(&nums))
	//fmt.Printf("i:%v\n", i)
	//for ii, v := range i {
	//	fmt.Printf("ii:%v, v:%v\n", ii, v)
	//}
}

func TestArray(t *testing.T) {
	var paramsList []int
	var data []*Dog
	paramsList = append(paramsList, 0)
	data = append(data, &Dog{Name: "dogx"})
	fmt.Printf("paramsList:%v\n", paramsList)
	x, _ := json.Marshal(data)
	fmt.Printf("data:%v\n", string(x))

	keys := make([]string, 0, 3)
	params := make([]*dto.Dog, 0, 3)
	for i, item := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		if item%2 == 0 {
			keys = append(keys, fmt.Sprintf("key:%d", i))
			params = append(params, &dto.Dog{Name: "dogx", Age: item})
		}
	}

	fmt.Printf("keys:%v\n", keys)
	newKeys := make([]string, 0, 3)
	fmt.Printf("newKeys len:%v cap:%v\n", len(newKeys), cap(newKeys))
	for i, item := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		if item != 0 {
			newKeys = append(newKeys, fmt.Sprintf("key:%d", i))
			fmt.Printf("newKeys len:%v cap:%v\n", len(newKeys), cap(newKeys))
		}
	}

	fmt.Printf("newKeys len:%v cap:%v\n", len(newKeys), cap(newKeys))
}
