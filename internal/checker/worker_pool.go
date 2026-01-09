package checker

import (
	"context"
	"sync"
)

type Job struct {
	ID   int
	Type string // "whois" or "ipquality"
	Data interface{}
}

type WorkerPool struct {
	poolSize int
	jobs     chan Job
	results  chan interface{}
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewWorkerPool(size int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		poolSize: size,
		jobs:     make(chan Job, 1000),
		results:  make(chan interface{}, 1000),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (wp *WorkerPool) Start(workerFunc func(Job) interface{}) {
	for i := 0; i < wp.poolSize; i++ {
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			for {
				select {
				case <-wp.ctx.Done():
					return
				case job, ok := <-wp.jobs:
					if !ok {
						return
					}
					result := workerFunc(job)
					wp.results <- result
				}
			}
		}()
	}
}

func (wp *WorkerPool) AddJob(job Job) {
	wp.jobs <- job
}

func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
}

func (wp *WorkerPool) Results() <-chan interface{} {
	return wp.results
}
