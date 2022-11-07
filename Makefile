#!/usr/bin/make -f

.ONESHELL:
SHELLFLAGS := -u nounset -ec

THIS_MAKEFILE := $(realpath $(lastword $(MAKEFILE_LIST)))
THIS_DIR      := $(shell dirname $(THIS_MAKEFILE))
THIS_PROJECT  := nq

.PHONY: serve watch

build-env:
	go mod download

build: build-env
	go build

run: build
	./nq -- sh -c 'sleep 2 && echo foo'
	./nq -- sh -c 'sleep 2 && echo bar'

serve:
	./nq -s

watch:
	watcher
