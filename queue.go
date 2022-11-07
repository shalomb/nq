package main

import (
	"container/list"
	"time"

	log "github.com/sirupsen/logrus"
)

type Queue struct {
	queue   string
	fifo    string
	workers *list.List
}

func (q *Queue) setup() int {
	log.Printf("Setting up queue")
	q.workers = list.New()
	return 0
}

func (q *Queue) enqueue(j *Job) int {
	log.Debugf("Enqueueing %v", j)
	q.workers.PushBack(j)
	return 0
}

func (q *Queue) process() (int, time.Duration) {

	log.Debugf("  q.workers len: %d", q.workers.Len())
	defer q.workers.Init()
	// Take only the last item pushed
	e := q.workers.Back()
	j := e.Value.(*Job)
	r, t := j.exec()

	return r, t
}
