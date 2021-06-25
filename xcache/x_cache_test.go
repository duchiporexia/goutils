package xcache

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type Dog struct {
	Name string
	Age  int
}

func TestXCache(t *testing.T) {
	cache := NewXCache(&XCacheConfig{
		ShardCount:   4,
		EvictionTime: time.Second * 12,
		CleanPeriod:  time.Second * 3,
		CleanMaxSize: 1000000,
	})
	k1 := "key1"
	cache.Set(k1, []byte("this is the value"))
	intf, ok := cache.Get(k1)
	val := intf.([]byte)
	if ok {
		fmt.Printf("%v\n", string(val))
	}
	cache.Set(k1, []byte("this has been changed"))
	intf, ok = cache.Get(k1)
	val = intf.([]byte)
	if ok {
		fmt.Printf("%v\n", string(val))
	}
	for i := 0; i < 100; i++ {
		go func() {
			v, _ := cache.GetOrFetch("kef", func(key string) (interface{}, error) {
				fmt.Printf("trying to load key:%v ......\n", key)
				time.Sleep(time.Second * 2)
				return key + "1", nil
			})
			fmt.Printf("v:%v\n", v)
		}()
	}
	valF, ok := cache.GetOrFetch("kef", nil)
	fmt.Printf("valf:%v\n", valF)
	fmt.Printf("start loop .....\n")
	for {
		time.Sleep(time.Second * 2)
		intf, ok = cache.Get(k1)
		val, _ = intf.([]byte)
		if ok {
			fmt.Printf("%v\n", string(val))
		} else {
			fmt.Printf("deleted\n")
		}
	}
	cache.Close()
}

func TestXCacheLoad(t *testing.T) {
	cache := NewXCache(&XCacheConfig{})
	N := 20000000
	insertMaxValue := time.Duration(0)
	insertChan := make(chan time.Duration, 20)

	deleteMaxValue := time.Duration(0)
	deleteChan := make(chan time.Duration, 20)

	getMaxValue := time.Duration(0)
	getChan := make(chan time.Duration, 20)

	fmt.Printf("load data ....\n")
	for i := 0; i < N; i++ {
		key := fmt.Sprintf("key:%v", i)
		cache.Set(key, &Dog{Name: "dog", Age: 33})
	}
	fmt.Printf("load data: done.\n")
	for i := 0; i < 200; i++ {
		go func() {
			for true {
				name := fmt.Sprintf("dog:%v", rand.Intn(99999999))
				start := time.Now()
				cache.Set(name, &Dog{Name: "dog", Age: 33})
				elapsed := time.Since(start)
				if elapsed > insertMaxValue {
					insertChan <- elapsed
				}
				/////////////////////////////////////////////
				start = time.Now()
				cache.Delete(name)
				elapsed = time.Since(start)
				if elapsed > deleteMaxValue {
					deleteChan <- elapsed
				}
			}
		}()
	}
	for i := 0; i < 1000; i++ {
		go func() {
			for true {
				name := fmt.Sprintf("dog:%v", rand.Intn(99999999))
				start := time.Now()
				_, ok := cache.Get(name)
				if ok {
				}
				elapsed := time.Since(start)
				if elapsed > getMaxValue {
					getChan <- elapsed
				}
			}
		}()
	}

	go func() {
		for true {
			select {
			case v := <-getChan:
				if v > getMaxValue {
					getMaxValue = v
				}
			case v := <-deleteChan:
				if v > deleteMaxValue {
					deleteMaxValue = v
				}
			case v := <-insertChan:
				if v > insertMaxValue {
					insertMaxValue = v
				}
			}
		}
	}()

	for true {
		time.Sleep(time.Second)
		fmt.Printf("insertMaxValue:%v, getMaxValue:%v, deleteMaxValue:%v\n", insertMaxValue, getMaxValue, deleteMaxValue)
		name := fmt.Sprintf("dog:%v", rand.Intn(99999999))
		start := time.Now()
		getN := 1000000
		for i := 0; i < getN; i++ {
			_, ok := cache.Get(name)
			if ok {
			}
		}

		elapsed := time.Since(start)
		fmt.Printf("getMaxValue (%v):%v\n", getN, elapsed)
	}
}
