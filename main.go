// Package main this is the package content
package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sync"

	log "github.com/sirupsen/logrus"
)

func main() {
	opts = parseOpts()

	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "15:04:05.9999"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	if server {
		var wg sync.WaitGroup
		defer wg.Done()
		log.Printf("Server requested. Starting server")
		startServer()
		wg.Wait()
		os.Exit(0)
	}

	if _, err := os.Stat(namedPipe); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Error finding named pipe: %v (Server not started?)", namedPipe)
	}

	f, err := os.OpenFile(namedPipe, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("Error opening named pipe: %v", err)
	}

	// TODO: Take CLI args or stdin
	makeRequest := func(args []string) {
		b, e := json.Marshal(args)
		if e != nil {
			log.Errorf("Error marshalling arguments to JSON")
		}

		log.Printf("Send job: %s", string(b))
		s, e := f.WriteString(string(b) + "\n")
		if e != nil {
			log.Warn(fmt.Sprintf("Server failed to acknowledge request: %v (%v)", e, s))
			os.Exit(-1)
		}
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		makeRequest(opts.Args())
	} else {
		r, _ := regexp.Compile(*inputPattern)
		log.Printf("Stdin is a pipe! Matching stdin against: %v", r)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			in := scanner.Text()
			if r.MatchString(in) {
				log.Printf("%v matches %v", in, r)
				makeRequest([]string{in})
			}
		}
		if err := scanner.Err(); err != nil {
			log.Fatalf("Error reading in from stdin: %v", err)
		}
	}
	os.Exit(0)
}
