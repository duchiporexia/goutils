// Code generated by xvv, DO NOT EDIT.
// HandlerType: Create

package test

import (
	"servicehub/common/batch_handler/test/dto"
	"sync"
	"time"
)

type TCreateHvHandlerConfig struct {
	Do           func(paramsList []dto.DogPo) ([]*dto.Dog, []error)
	Wait         time.Duration
	MaxBatchSize int
}

type TCreateHvHandler struct {
	Do           func(paramsList []dto.DogPo) ([]*dto.Dog, []error)
	Wait         time.Duration
	MaxBatchSize int
	batch        *tCreateHvHandlerBatch
	m            sync.Mutex
}

func NewTCreateHvHandler(config TCreateHvHandlerConfig) *TCreateHvHandler {
	return &TCreateHvHandler{
		Do:           config.Do,
		Wait:         config.Wait,
		MaxBatchSize: config.MaxBatchSize,
	}
}

func (s *TCreateHvHandler) Query(params dto.DogPo) (*dto.Dog, error) {
	return s.QueryFuture(params)()
}

func (s *TCreateHvHandler) QueryFuture(params dto.DogPo) func() (*dto.Dog, error) {
	s.m.Lock()
	if s.batch == nil {
		s.batch = &tCreateHvHandlerBatch{done: make(chan struct{})}
	}
	batch := s.batch
	pos := batch.getPosition(s, params)
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

func (s *TCreateHvHandler) QueryAll(paramsList []dto.DogPo) ([]*dto.Dog, []error) {
	return s.QueryAllFuture(paramsList)()
}

func (s *TCreateHvHandler) QueryAllFuture(paramsList []dto.DogPo) func() ([]*dto.Dog, []error) {
	thunks := make([]func() (*dto.Dog, error), len(paramsList))
	for i, item := range paramsList {
		thunks[i] = s.QueryFuture(item)
	}

	return func() ([]*dto.Dog, []error) {
		errors := make([]error, len(paramsList))
		results := make([]*dto.Dog, len(paramsList))
		for i, thunk := range thunks {
			results[i], errors[i] = thunk()
		}
		return results, errors
	}
}

type tCreateHvHandlerBatch struct {
	paramsList []dto.DogPo
	data       []*dto.Dog
	error      []error
	closing    bool
	done       chan struct{}
}

func (s *tCreateHvHandlerBatch) getPosition(handler *TCreateHvHandler, params dto.DogPo) int {
	pos := len(s.paramsList)
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

func (s *tCreateHvHandlerBatch) waiting(handler *TCreateHvHandler) {
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

func (s *tCreateHvHandlerBatch) finish(handler *TCreateHvHandler) {
	s.data, s.error = handler.Do(s.paramsList)
	close(s.done)

}
