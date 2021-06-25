package test

import (
	"fmt"
	"servicehub/common/batch_handler/test/dto"
	"testing"
	"time"
)

const (
	wait         = time.Millisecond * 2
	maxBatchSize = 5
)

func TestTCreateHandlers(t *testing.T) {
	handler := NewTCreateNvHandler(TCreateNvHandlerConfig{
		Do: func(paramsList []dto.DogPo) []error {
			fmt.Printf("create paramsList: %v\n", paramsList)
			return nil
		},
		Wait:         wait,
		MaxBatchSize: maxBatchSize,
	})
	thunk1 := handler.QueryAllFuture([]dto.DogPo{{Name: "dog1", Age: 1}, {Name: "dog2", Age: 2}})
	thunk2 := handler.QueryFuture(dto.DogPo{Name: "dog3", Age: 3})
	thunk3 := handler.QueryFuture(dto.DogPo{Name: "dog4", Age: 4})
	_, _, _ = thunk1(), thunk2(), thunk3()
}
