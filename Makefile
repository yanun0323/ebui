-include Makefile.env

.PHONY: help test test-race test-bench build-examples run-hello run-playground verify

help:
	@printf "%s\n" \
	"make test              - go test ./..." \
	"make test-race         - go test -race ./..." \
	"make test-bench        - go test ./internal/runtime -run '^$$' -bench ." \
	"make build-examples    - go build ./examples/..." \
	"make run-hello         - go run ./examples/hello" \
	"make run-playground    - go run ./examples/playground" \
	"make verify            - run test, race, and build-examples"

test:
	go test ./...

test-race:
	go test -race ./...

test-bench:
	go test ./internal/runtime -run '^$$' -bench .

build-examples:
	go build ./examples/...

run-hello:
	go run ./examples/hello

run-playground:
	go run ./examples/playground

verify: test test-race build-examples
