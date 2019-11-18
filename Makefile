SHELL = /bin/bash

.PHONY: build
build:
	go build -o bin/octant-helm cmd/octant-helm/main.go

.PHONY: dev
dev:
	scripts/dev.sh
