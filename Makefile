#!/usr/bin/make -f

.ONESHELL:
SHELLFLAGS := -u nounset -ec

THIS_MAKEFILE := $(realpath $(lastword $(MAKEFILE_LIST)))
THIS_DIR      := $(shell dirname $(THIS_MAKEFILE))
THIS_PROJECT  := nq

.PHONY: serve watch

build:
	go build

run: build
	./nq

serve:
	while :; do ./main; done

watch:
	watcher
