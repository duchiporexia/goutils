// Code generated by xvv, DO NOT EDIT.
// HandlerType: Update

package test

import (
	"sync"
	"time"

	"github.com/duchiporexia/goutils/batch_handler/test/dto"
)

type TUpdateNpHvHandlerConfig struct {
	Do           func(keys []int) ([]*dto.Dog, []error)
	Wait         time.Duration
	MaxBatchSize int
	CacheDel     func(keys ...int) error
}

type TUpdateNpHvHandler struct {
	Do           func(keys []int) ([]*dto.Dog, []error)
	Wait         time.Duration
	MaxBatchSize int
	CacheDel     func(keys ...int) error
	batch        *tUpdateNpHvHandlerBatch
	m            sync.Mutex
}

func NewTUpdateNpHvHandler(config TUpdateNpHvHandlerConfig) *TUpdateNpHvHandler {
	return &TUpdateNpHvHandler{
		Do:           config.Do,
		Wait:         config.Wait,
		MaxBatchSize: config.MaxBatchSize,
		CacheDel:     config.CacheDel,
	}
}

func (s *TUpdateNpHvHandler) Query(key int) (*dto.Dog, error) {
	return s.QueryFuture(key)()
}

func (s *TUpdateNpHvHandler) QueryFuture(key int) func() (*dto.Dog, error) {
	s.m.Lock()
	if s.batch == nil {
		s.batch = &tUpdateNpHvHandlerBatch{done: make(chan struct{})}
	}
	batch := s.batch
	pos := batch.getPosition(s, key)
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

func (s *TUpdateNpHvHandler) QueryAll(keys []int) ([]*dto.Dog, []error) {
	return s.QueryAllFuture(keys)()
}

func (s *TUpdateNpHvHandler) QueryAllFuture(keys []int) func() ([]*dto.Dog, []error) {
	thunks := make([]func() (*dto.Dog, error), len(keys))
	for i, item := range keys {
		thunks[i] = s.QueryFuture(item)
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

type tUpdateNpHvHandlerBatch struct {
	keys    []int
	data    []*dto.Dog
	error   []error
	closing bool
	done    chan struct{}
}

func (s *tUpdateNpHvHandlerBatch) getPosition(handler *TUpdateNpHvHandler, key int) int {
	for i, existingKey := range s.keys {
		if key == existingKey {
			return i
		}
	}
	pos := len(s.keys)
	s.keys = append(s.keys, key)
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

func (s *tUpdateNpHvHandlerBatch) waiting(handler *TUpdateNpHvHandler) {
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

func (s *tUpdateNpHvHandlerBatch) finish(handler *TUpdateNpHvHandler) {
	s.data, s.error = handler.Do(s.keys)
	close(s.done)

	if handler.CacheDel != nil {
		_ = handler.CacheDel(s.keys...)
	}
}
