package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"servicehub/common/batch_handler/test/dto"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestTReadHandlers(t *testing.T) {
	cache := make(map[string]*dto.Dog)
	handler := NewTReadHpHandler(TReadHpHandlerConfig{
		CacheBatchSet: func(keys []string, items []*dto.Dog) error {
			fmt.Printf("CacheBatchSet: keys: %s\n", keys)
			for i, key := range keys {
				cache[key] = items[i]
			}
			return nil
		},
		CacheBatchGet: func(keys []string) ([]*dto.Dog, error) {
			fmt.Printf("CacheBatchGet: keys: %s\n", keys)
			values := make([]*dto.Dog, len(keys))
			for i, key := range keys {
				values[i] = cache[key]
			}
			return values, nil
		},
		Do: func(keys []string, paramsList []dto.DogPo) ([]*dto.Dog, []error) {
			fmt.Printf("Do: keys: %s\n", keys)
			var values []*dto.Dog
			for i, dogP := range paramsList {
				values = append(values, &dto.Dog{
					Id:   keys[i],
					Name: dogP.Name,
					Age:  dogP.Age,
				})
			}
			return values, nil
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

	wg.Wait()
}
