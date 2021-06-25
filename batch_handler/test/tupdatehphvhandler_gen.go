// Code generated by xvv, DO NOT EDIT.
// HandlerType: Update

package test

import (
	"sync"
	"time"

	"github.com/duchiporexia/goutils/batch_handler/test/dto"
)

type TUpdateHpHvHandlerConfig struct {
	Do           func(keys []int, paramsList []dto.DogPo) ([]*dto.Dog, []error)
	Wait         time.Duration
	MaxBatchSize int
	CacheDel     func(keys ...int) error
}

type TUpdateHpHvHandler struct {
	Do           func(keys []int, paramsList []dto.DogPo) ([]*dto.Dog, []error)
	Wait         time.Duration
	MaxBatchSize int
	CacheDel     func(keys ...int) error
	batch        *tUpdateHpHvHandlerBatch
	m            sync.Mutex
}

func NewTUpdateHpHvHandler(config TUpdateHpHvHandlerConfig) *TUpdateHpHvHandler {
	return &TUpdateHpHvHandler{
		Do:           config.Do,
		Wait:         config.Wait,
		MaxBatchSize: config.MaxBatchSize,
		CacheDel:     config.CacheDel,
	}
}

func (s *TUpdateHpHvHandler) Query(key int, params dto.DogPo) (*dto.Dog, error) {
	return s.QueryFuture(key, params)()
}

func (s *TUpdateHpHvHandler) QueryFuture(key int, params dto.DogPo) func() (*dto.Dog, error) {
	s.m.Lock()
	if s.batch == nil {
		s.batch = &tUpdateHpHvHandlerBatch{done: make(chan struct{})}
	}
	batch := s.batch
	pos := batch.getPosition(s, key, params)
	s.m.Unlock()

	return func() (*dto.Dog, error) {
		<-batch.done
		var data *dto.Dog
		if pos < len(batch.data) {
			data = batch.data[pos]
		}
		var err error
		// its convenient to be able to return a single error for everything
		if len(batch.error) == 1 {
			err = batch.error[0]
		} else if batch.error != nil {
			err = batch.error[pos]
		}

		return data, err
	}
}

func (s *TUpdateHpHvHandler) QueryAll(keys []int, paramsList []dto.DogPo) ([]*dto.Dog, []error) {
	return s.QueryAllFuture(keys, paramsList)()
}

func (s *TUpdateHpHvHandler) QueryAllFuture(keys []int, paramsList []dto.DogPo) func() ([]*dto.Dog, []error) {
	thunks := make([]func() (*dto.Dog, error), len(keys))
	for i, item := range keys {
		thunks[i] = s.QueryFuture(item, paramsList[i])
	}

	return func() ([]*dto.Dog, []error) {
		errors := make([]error, len(keys))
		results := make([]*dto.Dog, len(keys))
		for i, thunk := range thunks {
			results[i], errors[i] = thunk()
		}
		return results, errors
	}
}

type tUpdateHpHvHandlerBatch struct {
	keys       []int
	paramsList []dto.DogPo
	data       []*dto.Dog
	error      []error
	closing    bool
	done       chan struct{}
}

func (s *tUpdateHpHvHandlerBatch) getPosition(handler *TUpdateHpHvHandler, key int, params dto.DogPo) int {
	for i, existingKey := range s.keys {
		if key == existingKey {
			return i
		}
	}
	pos := len(s.keys)
	s.keys = append(s.keys, key)
	s.paramsList = append(s.paramsList, params)
	if pos == 0 {
		go s.waiting(handler)
	}

	if handler.MaxBatchSize != 0 && pos >= handler.MaxBatchSize-1 {
		if !s.closing {
			s.closing = true
			handler.batch = nil
			go s.finish(handler)
		}
	}

	return pos
}

func (s *tUpdateHpHvHandlerBatch) waiting(handler *TUpdateHpHvHandler) {
	time.Sleep(handler.Wait)
	handler.m.Lock()
	if s.closing {
		handler.m.Unlock()
		return
	}
	handler.batch = nil
	handler.m.Unlock()

	s.finish(handler)
}

func (s *tUpdateHpHvHandlerBatch) finish(handler *TUpdateHpHvHandler) {
	s.data, s.error = handler.Do(s.keys, s.paramsList)
	close(s.done)

	if handler.CacheDel != nil {
		_ = handler.CacheDel(s.keys...)
	}
}
