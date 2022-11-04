package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {

	// log.SetFlags(log.Lshortfile)
	// log.SetPrefix(time.Now().UTC().Format("2006-01-02T15:04:05") + ": ")

	j := new(Job)
	j.parse(os.Args[1:])
	log.Printf("exec")
	e := j.exec()
	os.Exit(e)
}
