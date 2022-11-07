package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pborman/getopt/v2"
)

var (
	server       bool
	inputPattern *string
	namedPipe    = filepath.Join(os.Getenv("TMP"), "nq.fifo")
	timeout      = time.Second * 120
	verbose      bool
	opts         *getopt.Set
)

func init() {
	getopt.Flag(&verbose, 'v', "verbose")
	getopt.Lookup('v').SetOptional()
	getopt.Flag(&server, 's', "server")
	inputPattern = getopt.StringLong("pattern", 'p', ".")
}

func parseOpts() *getopt.Set {
	var opts = getopt.CommandLine
	opts.Parse(os.Args)
	return opts
}
