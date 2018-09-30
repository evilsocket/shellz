package queue

import (
	"runtime"
	"sync"

	"github.com/evilsocket/shellz/log"
)

type Job interface{}
type Logic func(arg Job)

type workQueue struct {
	workers  int
	jobChan  chan Job
	stopChan chan struct{}
	jobs     sync.WaitGroup
	done     sync.WaitGroup
	logic    Logic
}

func New(workers int, logic Logic) *workQueue {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}
	log.Debug("starting job queue with %d workers ...", workers)

	wq := &workQueue{
		workers:  workers,
		jobChan:  make(chan Job),
		stopChan: make(chan struct{}),
		jobs:     sync.WaitGroup{},
		done:     sync.WaitGroup{},
		logic:    logic,
	}

	for i := 0; i < workers; i++ {
		wq.done.Add(1)
		go wq.worker(i)
	}

	return wq
}

func (wq *workQueue) worker(id int) {
	defer wq.done.Done()

	log.Debug("job queue worker %d started", id)

	for {
		select {
		case <-wq.stopChan:
			log.Debug("worker %d received stop signal", id)
			return

		case job := <-wq.jobChan:
			log.Debug("worker %d got job", id)
			wq.logic(job)
			wq.jobs.Done()
		}
	}
}

func (wq *workQueue) Add(job Job) {
	wq.jobs.Add(1)
	wq.jobChan <- job
}

func (wq *workQueue) Wait() {
	wq.done.Wait()
}

func (wq *workQueue) WaitDone() {
	wq.jobs.Wait()
}

func (wq *workQueue) Stop() {
	close(wq.jobChan)
}
