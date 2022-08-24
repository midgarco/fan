package fan

import (
	"sync"
	"time"

	"github.com/midgarco/fan/log"
	"github.com/midgarco/fan/worker"
)

type Workers struct {
	ready   chan bool
	jobs    chan interface{}
	results chan *worker.Result
	pool    []*poolWorker
	done    chan bool

	mu sync.RWMutex

	Logger log.Interface

	WorkFunc   workfn
	ResultFunc resultfn

	ticker   *time.Ticker
	shutdown bool
}

type workfn func(w *worker.Details) error
type resultfn func(res *worker.Result) error

func (m *Workers) result() {
	for ev := range m.results {
		if err := m.ResultFunc(ev); err != nil {
			m.Logger.Errorf("[fan] resultFunc failed: %s\n", err.Error())
		}
		m.ready <- true
	}
}

type Config struct {
	Workers int

	Logger log.Interface

	WorkFunc   workfn
	ResultFunc resultfn
}

func (cfg *Config) CreateWokerPool() *Workers {
	mgr := &Workers{}
	mgr.ready = make(chan bool, cfg.Workers)
	mgr.jobs = make(chan interface{})
	mgr.results = make(chan *worker.Result)
	mgr.done = make(chan bool)
	mgr.Logger = cfg.Logger
	mgr.mu = sync.RWMutex{}
	mgr.WorkFunc = cfg.WorkFunc
	mgr.ResultFunc = cfg.ResultFunc

	go mgr.result()

	mgr.ticker = time.NewTicker(time.Second)

	// tick progress of workers
	go func() {
		for {
			select {
			case <-mgr.ticker.C:
				for _, w := range mgr.pool {
					mgr.Logger.Debugf("[fan] worker %d status: %s\n", w.id, w.status)
				}
			}
		}
	}()

	go func() {
		var wg sync.WaitGroup
		for i := 0; i < cfg.Workers; i++ {
			wg.Add(1)

			w := &poolWorker{id: int64(i + 1)}
			w.status = WorkerStatus_IDLE
			w.workFunc = cfg.WorkFunc

			w.logger = cfg.Logger

			w.readyChan = mgr.ready
			w.resultsChan = mgr.results
			w.jobChan = &mgr.jobs
			w.doneChan = make(chan bool)

			mgr.pool = append(mgr.pool, w)

			go w.doWork(&wg)
		}
		wg.Wait()

		mgr.Logger.Debug("[fan] close results/ready channels")
		close(mgr.results)
		close(mgr.ready)
		close(mgr.jobs)

		mgr.Logger.Debug("[fan] stop progress ticker")
		mgr.ticker.Stop()

		mgr.done <- true
	}()

	return mgr
}

func (mgr *Workers) Work(payload interface{}) {
	if mgr.shutdown {
		return
	}

	mgr.Logger.Debug("[fan] processing payload")

	<-mgr.ready
	mgr.jobs <- payload
}

func (mgr *Workers) Shutdown(done chan<- bool) {
	mgr.Close()
	done <- true
}

func (mgr *Workers) Close() error {
	mgr.shutdown = true

	for _, w := range mgr.pool {
		w.doneChan <- true
	}

	<-mgr.done
	return nil
}
