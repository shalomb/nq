#!/usr/bin/make -f

.ONESHELL:
SHELLFLAGS := -u nounset -ec

THIS_MAKEFILE := $(realpath $(lastword $(MAKEFILE_LIST)))
THIS_DIR      := $(shell dirname $(THIS_MAKEFILE))
THIS_PROJECT  := nq

.PHONY: serve watch

build: build-env
	go build

build-env:
	go mod download

run: build
	./nq -- sh -c 'sleep 1 && echo foo'
	sleep 0.1
	./nq -- sh -c 'sleep 2 && echo not processed'
	./nq -- sh -c 'sleep 2 && echo not processed'
	./nq -- sh -c 'sleep 2 && echo not processed'
	./nq -- sh -c 'sleep 2 && echo not processed'
	./nq -- sh -c 'sleep 2 && echo not processed'
	./nq -- sh -c 'sleep 2 && echo not processed'
	./nq -- sh -c 'sleep 2 && echo not processed'
	./nq -- sh -c 'sleep 1 && echo bar'
	sleep 1.5
	./nq -- sh -c 'sleep 1 && echo baz'

serve:
	./nq -s

watch:
	watcher
