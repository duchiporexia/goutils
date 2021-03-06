// Code generated by xvv, DO NOT EDIT.
// HandlerType: Read

package test

import (
	"sync"
	"time"

	"github.com/duchiporexia/goutils/batch_handler/test/dto"
)

type TReadNpHandlerConfig struct {
	Do            func(keys []string) ([]*dto.Dog, []error)
	Wait          time.Duration
	MaxBatchSize  int
	CacheBatchGet func(keys []string) ([]*dto.Dog, error)
	CacheBatchSet func(keys []string, items []*dto.Dog) error
}

type TReadNpHandler struct {
	Do            func(keys []string) ([]*dto.Dog, []error)
	Wait          time.Duration
	MaxBatchSize  int
	CacheBatchGet func(keys []string) ([]*dto.Dog, error)
	CacheBatchSet func(keys []string, items []*dto.Dog) error
	batch         *tReadNpHandlerBatch
	m             sync.Mutex
}

func NewTReadNpHandler(config TReadNpHandlerConfig) *TReadNpHandler {
	return &TReadNpHandler{
		Do:            config.Do,
		Wait:          config.Wait,
		MaxBatchSize:  config.MaxBatchSize,
		CacheBatchGet: config.CacheBatchGet,
		CacheBatchSet: config.CacheBatchSet,
	}
}

func (s *TReadNpHandler) Query(key string) (*dto.Dog, error) {
	return s.QueryFuture(key)()
}

func (s *TReadNpHandler) QueryFuture(key string) func() (*dto.Dog, error) {
	s.m.Lock()
	if s.batch == nil {
		s.batch = &tReadNpHandlerBatch{done: make(chan struct{})}
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

func (s *TReadNpHandler) QueryAll(keys []string) ([]*dto.Dog, []error) {
	return s.QueryAllFuture(keys)()
}

func (s *TReadNpHandler) QueryAllFuture(keys []string) func() ([]*dto.Dog, []error) {
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

type tReadNpHandlerBatch struct {
	keys    []string
	data    []*dto.Dog
	error   []error
	closing bool
	done    chan struct{}
}

func (s *tReadNpHandlerBatch) getPosition(handler *TReadNpHandler, key string) int {
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

func (s *tReadNpHandlerBatch) waiting(handler *TReadNpHandler) {
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

func (s *tReadNpHandlerBatch) finish(handler *TReadNpHandler) {
	if s.fetchItemsWithCache(handler) {
		return
	}
	s.data, s.error = handler.Do(s.keys)
	close(s.done)

	if handler.CacheBatchSet != nil {
		if len(s.error) == 0 {
			_ = handler.CacheBatchSet(s.keys, s.data)
			return
		}
		var keys []string
		var data []*dto.Dog

		for i, err := range s.error {
			if err == nil {
				keys = append(keys, s.keys[i])
				data = append(data, s.data[i])
			}
		}

		if len(keys) > 0 {
			_ = handler.CacheBatchSet(keys, data)
			return
		}
	}
}

func (s *tReadNpHandlerBatch) fetchItemsWithCache(handler *TReadNpHandler) bool {
	if handler.CacheBatchGet == nil {
		return false
	}
	items, err := handler.CacheBatchGet(s.keys)
	if err != nil {
		return false
	}
	// collect nil items and fetch data with handle.Do
	newKeys := make([]string, 0, 8)
	for i, item := range items {
		if item == nil {
			newKeys = append(newKeys, s.keys[i])
		}
	}

	if len(newKeys) == 0 {
		s.data, s.error = items, nil
		close(s.done)
		return true
	}

	finalErrs := make([]error, len(items))
	newItems, errs := handler.Do(newKeys)
	idx := 0
	for i, item := range items {
		if item == nil {
			items[i] = newItems[idx]
			if len(errs) == 1 {
				finalErrs[i] = errs[0]
			} else if errs != nil {
				finalErrs[i] = errs[idx]
			}
			idx++
		}
	}
	// set and notify
	s.data, s.error = items, errs
	close(s.done)

	if errs == nil && handler.CacheBatchSet != nil {
		_ = handler.CacheBatchSet(newKeys, newItems)
	}
	return true
}
