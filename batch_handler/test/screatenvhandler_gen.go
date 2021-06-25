// Code generated by xvv, DO NOT EDIT.
// HandlerType: Create

package test

import (
	"sync"
	"time"

	"github.com/duchiporexia/goutils/batch_handler/test/dto"
)

type SCreateNvHandlerConfig struct {
	Do           func(params dto.DogPo) error
	Wait         time.Duration
	MaxBatchSize int
}

type SCreateNvHandler struct {
	Do           func(params dto.DogPo) error
	Wait         time.Duration
	MaxBatchSize int
	batch        *sCreateNvHandlerBatch
	m            sync.Mutex
}

func NewSCreateNvHandler(config SCreateNvHandlerConfig) *SCreateNvHandler {
	return &SCreateNvHandler{
		Do:           config.Do,
		Wait:         config.Wait,
		MaxBatchSize: config.MaxBatchSize,
	}
}

func (s *SCreateNvHandler) Query(params dto.DogPo) error {
	return s.QueryFuture(params)()
}

func (s *SCreateNvHandler) QueryFuture(params dto.DogPo) func() error {
	s.m.Lock()
	if s.batch == nil {
		s.batch = &sCreateNvHandlerBatch{done: make(chan struct{})}
	}
	singleFunc := s.batch.getSingleFunc(s, params)
	s.m.Unlock()
	return singleFunc
}

type sCreateNvHandlerBatch struct {
	paramsList   []dto.DogPo
	error        []error
	funcs        []func() error
	funcDoneList []chan struct{}
	closing      bool
	done         chan struct{}
}

func (s *sCreateNvHandlerBatch) getSingleFunc(handler *SCreateNvHandler, params dto.DogPo) func() error {
	pos := len(s.funcDoneList)
	funcDone := make(chan struct{})
	s.funcDoneList = append(s.funcDoneList, funcDone)
	s.error = append(s.error, nil)

	singleFunc := func() error {
		<-s.done

		var err error

		err = handler.Do(params)
		s.error[pos] = err
		close(funcDone)

		return err
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

func (s *sCreateNvHandlerBatch) waiting(handler *SCreateNvHandler) {
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
