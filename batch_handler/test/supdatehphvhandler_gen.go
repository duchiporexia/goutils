// Code generated by xvv, DO NOT EDIT.
// HandlerType: Update

package test

import (
	"sync"
	"time"

	"github.com/duchiporexia/goutils/batch_handler/test/dto"
)

type SUpdateHpHvHandlerConfig struct {
	Do           func(key int, params dto.DogPo) (*dto.Dog, error)
	Wait         time.Duration
	MaxBatchSize int
	CacheDel     func(keys ...int) error
}

type SUpdateHpHvHandler struct {
	Do           func(key int, params dto.DogPo) (*dto.Dog, error)
	Wait         time.Duration
	MaxBatchSize int
	CacheDel     func(keys ...int) error
	batch        *sUpdateHpHvHandlerBatch
	m            sync.Mutex
}

func NewSUpdateHpHvHandler(config SUpdateHpHvHandlerConfig) *SUpdateHpHvHandler {
	return &SUpdateHpHvHandler{
		Do:           config.Do,
		Wait:         config.Wait,
		MaxBatchSize: config.MaxBatchSize,
		CacheDel:     config.CacheDel,
	}
}

func (s *SUpdateHpHvHandler) Query(key int, params dto.DogPo) (*dto.Dog, error) {
	return s.QueryFuture(key, params)()
}

func (s *SUpdateHpHvHandler) QueryFuture(key int, params dto.DogPo) func() (*dto.Dog, error) {
	s.m.Lock()
	if s.batch == nil {
		s.batch = &sUpdateHpHvHandlerBatch{done: make(chan struct{})}
	}
	singleFunc := s.batch.getSingleFunc(s, key, params)
	s.m.Unlock()
	return singleFunc
}

type sUpdateHpHvHandlerBatch struct {
	keys         []int
	paramsList   []dto.DogPo
	data         []*dto.Dog
	error        []error
	funcs        []func() (*dto.Dog, error)
	funcDoneList []chan struct{}
	closing      bool
	done         chan struct{}
}

func (s *sUpdateHpHvHandlerBatch) getSingleFunc(handler *SUpdateHpHvHandler, key int, params dto.DogPo) func() (*dto.Dog, error) {
	for i, existingKey := range s.keys {
		if key == existingKey {
			pos := i
			return func() (*dto.Dog, error) {
				<-s.funcDoneList[pos]
				return s.data[pos], s.error[pos]
			}
		}
	}
	pos := len(s.keys)
	s.keys = append(s.keys, key)
	funcDone := make(chan struct{})
	s.funcDoneList = append(s.funcDoneList, funcDone)
	s.data = append(s.data, nil)
	s.error = append(s.error, nil)

	singleFunc := func() (*dto.Dog, error) {
		<-s.done
		var item *dto.Dog
		var err error

		item, err = handler.Do(key, params)
		s.data[pos] = item
		s.error[pos] = err
		close(funcDone)

		if handler.CacheDel != nil {
			_ = handler.CacheDel(key)
		}
		return item, err
	}

	s.funcs = append(s.funcs, singleFunc)

	if pos == 0 {
		go s.waiting(handler)
	}

	if handler.MaxBatchSize != 0 && pos >= handler.MaxBatchSize-1 {
		if !s.closing {
			s.closing = true
			handler.batch = nil
			close(s.done)
		}
	}

	return singleFunc
}

func (s *sUpdateHpHvHandlerBatch) waiting(handler *SUpdateHpHvHandler) {
	time.Sleep(handler.Wait)
	handler.m.Lock()
	if s.closing {
		handler.m.Unlock()
		return
	}
	handler.batch = nil
	handler.m.Unlock()

	close(s.done)
}
