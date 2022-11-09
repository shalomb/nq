package main

import (
	"bufio"
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rjeczalik/notify"
	log "github.com/sirupsen/logrus"
)

func sigintHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Printf("Running cleanup routines")
		os.Remove(namedPipe)
		os.Exit(1)
	}()
}

func startServer() {
	sigintHandler()

	err := syscall.Mkfifo(namedPipe, 0666)
	if err != nil {
		log.Fatalf("Error creating named pipe (%s): %v", namedPipe, err)
	}

	q := new(Queue)
	q.setup()

	var wg sync.WaitGroup
	defer wg.Done()
	wg.Add(1)
	log.Printf("Setting up worker")
	go worker(&wg, q)

	log.Printf("Started Server on named pipe: %v", namedPipe)
	file, err := os.OpenFile(namedPipe, os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}

	// https://stackoverflow.startServercom/a/45447384/742600
	var e notify.EventInfo
	c := make(chan notify.EventInfo, 5)
	notify.Watch(namedPipe, c, notify.Write|notify.Remove)
	reader := bufio.NewReader(file)

	for {
		// wait on events
		e = <-c

		switch e.Event() {
		case notify.Write:
			line, err := reader.ReadBytes('\n')
			if err == nil {
				var s []string
				json.Unmarshal(line, &s)
				log.Printf("Recv job: %v (%v)", s, len(s))

				if len(s) > 0 { // non-empty job
					j := new(Job)
					j.parse(s)

					q.enqueue(j)
				} else {
					log.Info("Ping acked. Empty job request!")
				}
			}

		case notify.Remove:
			log.Printf("file removed: %v", file)
		}
	}
}
