package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/midgarco/fan"
	"github.com/midgarco/fan/worker"
)

func main() {
	log := &MockLog{}

	wfn := func(worker *worker.Details) error {
		log.Infof("worker func: %v\n", worker)
		return nil
	}
	rfn := func(ev *worker.Result) error {
		log.Infof("result func: %v\n", ev)
		return nil
	}

	workers := &fan.Config{
		Workers:    2,
		Logger:     log,
		WorkFunc:   wfn,
		ResultFunc: rfn,
	}
	pool := workers.CreateWokerPool()

	// graceful exit
	done := make(chan bool)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan

		log.Info("graceful exit")

		// exit
		pool.Close()

		done <- true
	}()

	for {
		select {
		case <-done:
			break
		default:
		}

		pool.Work("foo")
		pool.Work("bar")
		pool.Work("baz")

		time.Sleep(5 * time.Second)
	}

	log.Info("done.")
}

type MockLog struct{}

func (*MockLog) Debug(msg string)                         { fmt.Println(msg) }
func (*MockLog) Info(msg string)                          { fmt.Println(msg) }
func (*MockLog) Warn(msg string)                          { fmt.Println(msg) }
func (*MockLog) Error(msg string)                         { fmt.Println(msg) }
func (*MockLog) Fatal(msg string)                         { fmt.Println(msg) }
func (*MockLog) Debugf(msg string, params ...interface{}) { fmt.Printf(msg, params...) }
func (*MockLog) Infof(msg string, params ...interface{})  { fmt.Printf(msg, params...) }
func (*MockLog) Warnf(msg string, params ...interface{})  { fmt.Printf(msg, params...) }
func (*MockLog) Errorf(msg string, params ...interface{}) { fmt.Printf(msg, params...) }
func (*MockLog) Fatalf(msg string, params ...interface{}) { fmt.Printf(msg, params...) }
