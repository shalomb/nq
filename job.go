package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Job struct {
	cmd       []string
	timestamp int64
}

func (j *Job) parse(c []string) {
	j.cmd = c
	j.timestamp = time.Now().UnixNano()
}

func (j *Job) exec() int {
	cmd := exec.Command(j.cmd[0], j.cmd[1:]...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			log.Errorf("%s: %d", err, exiterr.ExitCode())
			return exiterr.ExitCode()
		}
	}

	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	for _, line := range strings.Split(strings.TrimSuffix(outStr, "\n"), "\n") {
		log.Info(line)
	}

	for _, line := range strings.Split(strings.TrimSuffix(errStr, "\n"), "\n") {
		log.Warn(line)
	}

	return 0
}