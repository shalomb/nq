package main

import (
	"bufio"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

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
	log.Info("Starting server")
	sigintHandler()

	err := syscall.Mkfifo(namedPipe, 0666)
	if err != nil {
		log.Fatalf("Error creating named pipe (%s): %v", namedPipe, err)
	}

	log.Printf("Opening named pipe to read: %v", namedPipe)
	file, err := os.OpenFile(namedPipe, os.O_CREATE, os.ModeNamedPipe)
	if err != nil {
		log.Fatal("Open named pipe file error:", err)
	}

	reader := bufio.NewReader(file)

	q := new(Queue)
	q.setup()

	for {
		line, err := reader.ReadBytes('\n')
		if err == nil {
			var s []string
			json.Unmarshal(line, &s)
			log.Printf("Recv job: %v (%v)", s, len(s))

			if len(s) > 0 { // non-empty job
				j := new(Job)
				j.parse(s)

				q.enqueue(j)
				r, e := q.process()

				log.Printf("  Processed job: %+v (in %v)", r, e)
				log.Printf("  ")
			} else {
				log.Info("Ping acked. Empty job request!")
			}
		}
	}
}
