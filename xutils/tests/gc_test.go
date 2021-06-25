package tests

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

type gCDog struct {
	Id  int64
	Age int64
}

type gCDogValue struct {
	//key []byte : can't use this here.
	gCDog     gCDog
	collision int
}

func TestBit(t *testing.T) {
	fmt.Printf("%v", 1<<5-1)
	m := make(map[int][]byte, 8)
	v, ok := m[22]
	if v == nil {
		fmt.Printf("v:is nil\n")
	}
	fmt.Printf("v:%v, ok:%v\n", v, ok)
}
func TestSlice(t *testing.T) {
	answer1 := make([][]int, 5)
	fmt.Printf("%v %v\n", len(answer1), answer1[1])
	for i := 0; i < 5; i++ {
		//item := make([]int, 6)
		answer1[i] = nil
	}
	fmt.Printf("%v %v\n", len(answer1), answer1[0])
}
func TestMap(t *testing.T) {
	m := make(map[string]gCDog)
	name := "dog1"
	dog, ok := m[name]
	fmt.Printf("ok:%v, dog:%v\n", ok, dog)
	dog.Id = 11111
	dog2, ok := m[name]
	fmt.Printf("ok:%v, dog2:%v\n", ok, dog2)

	dog = gCDog{Id: 1, Age: 2}
	m[name] = dog
	dog3, ok := m[name]
	fmt.Printf("ok:%v, dog3:%v\n", ok, dog3)
	dog3.Id = 11
	dog4, ok := m[name]
	fmt.Printf("ok:%v, dog4:%v\n", ok, dog4)
	fmt.Printf("map:%v\n", m)
	delete(m, name)
	dog5, ok := m[name]
	fmt.Printf("ok:%v, dog5:%v\n", ok, dog5)
	dog.Id = 12
	dog6, ok := m[name]
	fmt.Printf("ok:%v, dog6:%v\n", ok, dog6)

	fmt.Printf("map2:%v\n", m)

	m2 := make(map[string]*gCDog)
	fmt.Printf("map3:%v\n", m2)
	m2[name] = &dog
	fmt.Printf("map4:%v\n", m2)
}

func TestDemo(t *testing.T) {
	//demoGC("1")
	//demoGC("2")
	//demoGC("3")
	//demoGC("4")
	//demoGC("5")
	//demoGC("5.0")
	//demoGC("5.0.1")
	//demoGC("5.1")
	//demoGC("5.2")
	//demoGC("6")
	demoGC("7")
}
func demoGC(n string) {
	var N int = 5e7 // 5000w
	var ShardN int = 100
	type Item struct {
		value interface{}
	}
	switch n {
	case "1": // Big map with a pointer in the value
		m := make(map[int32]*int32)
		for i := 0; i < N; i++ {
			n := int32(i)
			m[n] = &n
		}
		runtime.GC()
		fmt.Printf("With %T, GC took %s\n", m, timeGC())
		m = nil
	case "2":
		m := make(map[int32]int32)
		for i := 0; i < N; i++ {
			n := int32(i)
			m[n] = n
		}
		runtime.GC()
		fmt.Printf("With %T, GC took %s\n", m, timeGC())
		m = nil
	case "3":
		shards := make([]map[int32]*int32, ShardN)
		for i := range shards {
			shards[i] = make(map[int32]*int32)
		}
		for i := 0; i < N; i++ {
			n := int32(i)
			shards[i%ShardN][n] = &n
		}
		runtime.GC()
		fmt.Printf("With map shards (%T), GC took %s\n", shards, timeGC())
		shards = nil
	case "4":
		shards := make([]map[int32]int32, ShardN)
		for i := range shards {
			shards[i] = make(map[int32]int32)
		}
		for i := 0; i < N; i++ {
			n := int32(i)
			shards[i%ShardN][n] = n
		}
		runtime.GC()
		fmt.Printf("With map shards (%T), GC took %s\n", shards, timeGC())
		shards = nil
	case "5":
		m := make(map[int32]Item)
		for i := 0; i < N; i++ {
			n := int32(i)
			m[n] = Item{value: gCDog{Id: 333, Age: 3}} // gCDog{Id: 33333, Age: 3}
		}
		runtime.GC()
		fmt.Printf("With map shards (%T), GC took %s\n", m, timeGC())
		m = nil
	case "5.0":
		m := make(map[int32]string)
		for i := 0; i < N; i++ {
			n := int32(i)
			m[n] = fmt.Sprintf("%v", n)
		}
		runtime.GC()
		fmt.Printf("With map (%T), GC took %s\n", m, timeGC())
		m = nil
	case "5.0.1":
		m := make(map[string]int32)
		for i := 0; i < N; i++ {
			n := int32(i)
			m[fmt.Sprintf("%v", n)] = n
		}
		runtime.GC()
		fmt.Printf("With map (%T), GC took %s\n", m, timeGC())
		m = nil
	case "5.1":
		m := make(map[int32]gCDog)
		for i := 0; i < N; i++ {
			n := int32(i)
			m[n] = gCDog{Id: 33333, Age: 3}
		}
		runtime.GC()
		fmt.Printf("With map (%T), GC took %s\n", m, timeGC())
		m = nil

		m = nil
	case "5.2":
		shards := make([]map[int32]gCDog, ShardN)
		for i := range shards {
			shards[i] = make(map[int32]gCDog)
		}
		for i := 0; i < N; i++ {
			n := int32(i)
			shards[i%ShardN][n] = gCDog{Id: 33333, Age: 3}
		}
		runtime.GC()
		fmt.Printf("With map shards (%T), GC took %s\n", shards, timeGC())
		shards = nil
	case "6":
		shards := make([]map[int32]*gCDog, ShardN)
		for i := range shards {
			shards[i] = make(map[int32]*gCDog)
		}
		for i := 0; i < N; i++ {
			n := int32(i)
			shards[i%ShardN][n] = &gCDog{Id: 666666, Age: 3}
		}
		runtime.GC()
		fmt.Printf("With map shards (%T), GC took %s\n", shards, timeGC())
		shards = nil
	case "7":
		shards := make([]gCDogValue, 8)
		for i := 0; i < N; i++ {
			shards = append(shards, gCDogValue{collision: 20, gCDog: gCDog{Id: 22423, Age: 3}})
		}
		runtime.GC()
		fmt.Printf("With slice (%T), GC took %s\n", shards, timeGC())
		shards = nil
	}
}
func timeGC() time.Duration {
	start := time.Now()
	runtime.GC()
	return time.Since(start)
}
