package fan

import (
	"sync"

	"github.com/midgarco/fan/log"
	"github.com/midgarco/fan/worker"
)

type poolWorker struct {
	id       int64
	status   WorkerStatus
	workFunc workfn

	logger log.Interface

	readyChan   chan bool
	jobChan     *chan interface{}
	resultsChan chan *worker.Result
	doneChan    chan bool
}

type WorkerStatus int

const (
	WorkerStatus_IDLE = iota + 1
	WorkerStatus_PROCESSING
	WorkerStatus_COMPLETE
)

func (w *poolWorker) doWork(wg *sync.WaitGroup) {
	defer wg.Done()

	w.readyChan <- true

	for {
		select {
		case payload := <-*w.jobChan:
			done := make(chan bool)

			result := &worker.Result{
				Id:      w.id,
				Payload: payload,
			}

			go func(msg interface{}) {
				defer func() { done <- true }()

				w.status = WorkerStatus_PROCESSING
				w.logger.Debugf("[fan] worker: %d\n", w.id)

				if err := w.workFunc(&worker.Details{
					Id:      w.id,
					Payload: payload,
				}); err != nil {
					result.Error = err
				}

				w.status = WorkerStatus_COMPLETE
			}(payload)

			<-done
			close(done)

			w.resultsChan <- result

			w.status = WorkerStatus_IDLE

		case <-w.doneChan:
			return
		}
	}
}
