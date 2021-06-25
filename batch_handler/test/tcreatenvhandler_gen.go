// Code generated by xvv, DO NOT EDIT.
// HandlerType: Create

package test

import (
	"sync"
	"time"

	"github.com/duchiporexia/goutils/batch_handler/test/dto"
)

type TCreateNvHandlerConfig struct {
	Do           func(paramsList []dto.DogPo) []error
	Wait         time.Duration
	MaxBatchSize int
}

type TCreateNvHandler struct {
	Do           func(paramsList []dto.DogPo) []error
	Wait         time.Duration
	MaxBatchSize int
	batch        *tCreateNvHandlerBatch
	m            sync.Mutex
}

func NewTCreateNvHandler(config TCreateNvHandlerConfig) *TCreateNvHandler {
	return &TCreateNvHandler{
		Do:           config.Do,
		Wait:         config.Wait,
		MaxBatchSize: config.MaxBatchSize,
	}
}

func (s *TCreateNvHandler) Query(params dto.DogPo) error {
	return s.QueryFuture(params)()
}

func (s *TCreateNvHandler) QueryFuture(params dto.DogPo) func() error {
	s.m.Lock()
	if s.batch == nil {
		s.batch = &tCreateNvHandlerBatch{done: make(chan struct{})}
	}
	batch := s.batch
	pos := batch.getPosition(s, params)
	s.m.Unlock()

	return func() error {
		<-batch.done
		var err error
		// its convenient to be able to return a single error for everything
		if len(batch.error) == 1 {
			err = batch.error[0]
		} else if batch.error != nil {
			err = batch.error[pos]
		}

		return err
	}
}

func (s *TCreateNvHandler) QueryAll(paramsList []dto.DogPo) []error {
	return s.QueryAllFuture(paramsList)()
}

func (s *TCreateNvHandler) QueryAllFuture(paramsList []dto.DogPo) func() []error {
	thunks := make([]func() error, len(paramsList))
	for i, item := range paramsList {
		thunks[i] = s.QueryFuture(item)
	}

	return func() []error {
		errors := make([]error, len(paramsList))
		for i, thunk := range thunks {
			errors[i] = thunk()
		}
		return errors
	}
}

type tCreateNvHandlerBatch struct {
	paramsList []dto.DogPo
	error      []error
	closing    bool
	done       chan struct{}
}

func (s *tCreateNvHandlerBatch) getPosition(handler *TCreateNvHandler, params dto.DogPo) int {
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

func (s *tCreateNvHandlerBatch) waiting(handler *TCreateNvHandler) {
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

func (s *tCreateNvHandlerBatch) finish(handler *TCreateNvHandler) {
	s.error = handler.Do(s.paramsList)
	close(s.done)

}
