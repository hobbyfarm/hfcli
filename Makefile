SHELL := /bin/bash

build:
	go build -gcflags="-N -l" -o hfcli main.go

fmt:
	gofmt -l -s -w .

check:
	GOFMT_OUTPUT="$$(gofmt -d -e -l .)"; \
	if [ -n "$$GOFMT_OUTPUT" ]; then \
	  echo "All the following files are not correctly formatted"; \
	  echo "$${GOFMT_OUTPUT}"; \
	  exit 1;  \
	fi;  \
	echo "gofmt-output=Gofmt step succeed"

lint:
	golangci-lint run	