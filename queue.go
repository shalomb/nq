package main

import (
	"container/list"

	log "github.com/sirupsen/logrus"
)

type Queue struct {
	fifo  string
	queue *list.List
	in    chan string
	out   chan string
}

func (q *Queue) setup() {
	log.Printf("Setting up queue\n")

	q.queue = list.New()

	q.in = make(chan string, 1)
	q.out = make(chan string, 1)
}

func (q *Queue) enqueue(j *Job) int {
	log.Debugf("Enqueueing %v", j)
	q.queue.PushBack(j)
	log.Printf("  Enqueue done (%+v). q.queue len: %d", j.uuid, q.queue.Len())
	return 0
}
