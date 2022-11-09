package main

import (
	"os"
	"sync"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// TODO
// Cleanup channels
// Register restart handlers below correctly

var (
	stop = make(chan struct{})
	done = make(chan struct{})
)

func worker(wg *sync.WaitGroup, q *Queue) {
	defer wg.Done()
	log.Printf("Worker started")
LOOP:
	for {
		time.Sleep(175 * time.Millisecond) // this is work to be done by worker.
		select {
		case <-stop:
			break LOOP
		default:
			if q.queue.Len() > 0 {
				// Take only the last item pushed
				e := q.queue.Back()

				j := e.Value.(*Job)
				q.queue.Init()

				r, t := j.exec()

				log.Printf(" Processed job: %+v, %+v (in %v)", j.uuid, r, t)
				log.Printf(" ")
			}
		}
	}
	done <- struct{}{}
}

func termHandler(sig os.Signal) {
	log.Println("terminating...")
	stop <- struct{}{}
	if sig == syscall.SIGQUIT {
		<-done
	}
}

func reloadHandler(sig os.Signal) error {
	log.Println("configuration reloaded")
	return nil
}
