package test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"servicehub/common/batch_handler/test/dto"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSReadHandlers(t *testing.T) {
	cache := make(map[string]*dto.Dog)
	handler := NewSReadHpHandler(SReadHpHandlerConfig{
		CacheSet: func(key string, item *dto.Dog) error {
			fmt.Printf("CacheSet: key: %s\n", key)
			cache[key] = item
			return nil
		},
		CacheGet: func(key string) (*dto.Dog, error) {
			fmt.Printf("CacheGet: key: %s\n", key)
			if item, ok := cache[key]; ok {
				return item, nil
			}
			return nil, errors.New("no item in cache")
		},
		Do: func(key string, params dto.DogPo) (*dto.Dog, error) {
			fmt.Printf("Do => key: %s\n", key)
			return &dto.Dog{
				Id:   key,
				Name: params.Name,
				Age:  params.Age,
			}, nil
		},
		Wait:         wait,
		MaxBatchSize: maxBatchSize,
	})

	var wg sync.WaitGroup

	n := 6
	m := 7
	wg.Add(m * n)

	for i := 1; i <= n; i++ {
		idx := i
		go func() {
			defer wg.Done()
			dog, err := handler.QueryFuture(strconv.Itoa(idx), dto.DogPo{Name: fmt.Sprintf("dog%d", idx), Age: idx})()
			assert.NoError(t, err)
			assert.Equal(t, &dto.Dog{
				Id:   fmt.Sprintf("%d", idx),
				Name: fmt.Sprintf("dog%d", idx),
				Age:  idx,
			}, dog)
		}()
	}
	time.Sleep(time.Millisecond * 200)
	fmt.Printf("cache:%v\n", cache)
	for i := 1; i <= (m-1)*n; i++ {
		idx := (i % n) + 1
		go func() {
			defer wg.Done()
			dog, err := handler.QueryFuture(strconv.Itoa(idx), dto.DogPo{Name: fmt.Sprintf("dog%d", idx), Age: idx})()
			assert.NoError(t, err)
			assert.Equal(t, &dto.Dog{
				Id:   fmt.Sprintf("%d", idx),
				Name: fmt.Sprintf("dog%d", idx),
				Age:  idx,
			}, dog)
		}()
	}
}
