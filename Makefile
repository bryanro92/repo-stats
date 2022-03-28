SHELL = /bin/bash
COMMIT = $(shell git rev-parse --short HEAD)

install: build 
	chmod 555 ./repo-stats
	mv ./repo-stats /usr/local/bin/repo-stats

build:
	go build . 

run:
	go run . azure aro-rp 1

clean:  
	rm -f ./repo-stats

.PHONY: build install run clean