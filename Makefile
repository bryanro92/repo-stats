SHELL = /bin/bash
COMMIT = $(shell git rev-parse --short HEAD)

install: build 
	chmod 555 ./repo-stats
	mv ./repo-stats /usr/local/bin/repo-stats

build:
	go build -ldflags '-X github.com/bryanro92/repo-stats/pkg/version.GitCommit=$(COMMIT)' . 

run:
	go run -ldflags '-X github.com/bryanro92/repo-stats/pkg/version.GitCommit=$(COMMIT)' . 30

clean:  
	rm -f ./repo-stats

.PHONY: build install run clean