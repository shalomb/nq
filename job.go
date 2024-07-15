package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/sony/sonyflake"
)

var flake = sonyflake.NewSonyflake(sonyflake.Settings{})

// Job Struct
type Job struct {
	cmd       []string
	timestamp int64
	uuid      uint64
}

func newUUID() uint64 {
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("Error generating NextID(), %+v", err)
	}
	return id
}

func (j *Job) parse(c []string) {
	j.cmd = c
	j.timestamp = time.Now().UnixNano()
	j.uuid = newUUID()
}

func timeSince(start time.Time) time.Duration {
	return time.Since(start)
}

func (j *Job) exec() (int, time.Duration) {
	start := time.Now()
	cmd := exec.Command(j.cmd[0], j.cmd[1:]...)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	if err := cmd.Start(); err != nil {
		log.Errorf("%+v", err)
		return 127, timeSince(start)
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			log.Errorf("  %s: %d", err, exiterr.ExitCode())
			log.Printf("\a")
			return exiterr.ExitCode(), timeSince(start)
		}
	}

	outStr, errStr := stdoutBuf.String(), stderrBuf.String()

	if len(outStr) > 0 {
		for _, line := range strings.Split(strings.TrimSuffix(outStr, "\n"), "\n") {
			log.Info(line)
		}
	}

	if len(errStr) > 0 {
		for _, line := range strings.Split(strings.TrimSuffix(errStr, "\n"), "\n") {
			log.Warn(line)
		}
	}

	exitCode := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	return exitCode, timeSince(start)
}
