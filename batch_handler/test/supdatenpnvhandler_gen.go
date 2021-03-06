// Code generated by xvv, DO NOT EDIT.
// HandlerType: Update

package test

import (
	"sync"
	"time"
)

type SUpdateNpNvHandlerConfig struct {
	Do           func(key int) error
	Wait         time.Duration
	MaxBatchSize int
	CacheDel     func(keys ...int) error
}

type SUpdateNpNvHandler struct {
	Do           func(key int) error
	Wait         time.Duration
	MaxBatchSize int
	CacheDel     func(keys ...int) error
	batch        *sUpdateNpNvHandlerBatch
	m            sync.Mutex
}

func NewSUpdateNpNvHandler(config SUpdateNpNvHandlerConfig) *SUpdateNpNvHandler {
	return &SUpdateNpNvHandler{
		Do:           config.Do,
		Wait:         config.Wait,
		MaxBatchSize: config.MaxBatchSize,
		CacheDel:     config.CacheDel,
	}
}

func (s *SUpdateNpNvHandler) Query(key int) error {
	return s.QueryFuture(key)()
}

func (s *SUpdateNpNvHandler) QueryFuture(key int) func() error {
	s.m.Lock()
	if s.batch == nil {
		s.batch = &sUpdateNpNvHandlerBatch{done: make(chan struct{})}
	}
	singleFunc := s.batch.getSingleFunc(s, key)
	s.m.Unlock()
	return singleFunc
}

type sUpdateNpNvHandlerBatch struct {
	keys         []int
	error        []error
	funcs        []func() error
	funcDoneList []chan struct{}
	closing      bool
	done         chan struct{}
}

func (s *sUpdateNpNvHandlerBatch) getSingleFunc(handler *SUpdateNpNvHandler, key int) func() error {
	for i, existingKey := range s.keys {
		if key == existingKey {
			pos := i
			return func() error {
				<-s.funcDoneList[pos]
				return s.error[pos]
			}
		}
	}
	pos := len(s.keys)
	s.keys = append(s.keys, key)
	funcDone := make(chan struct{})
	s.funcDoneList = append(s.funcDoneList, funcDone)
	s.error = append(s.error, nil)

	singleFunc := func() error {
		<-s.done

		var err error

		err = handler.Do(key)
		s.error[pos] = err
		close(funcDone)

		if handler.CacheDel != nil {
			_ = handler.CacheDel(key)
		}
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

func (s *sUpdateNpNvHandlerBatch) waiting(handler *SUpdateNpNvHandler) {
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
