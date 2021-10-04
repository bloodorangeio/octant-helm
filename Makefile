SHELL = /bin/bash

PLUGIN_NAME=octant-helm

ifdef XDG_CONFIG_HOME
	OCTANT_PLUGIN_DIR ?= ${XDG_CONFIG_HOME}/octant/plugins
else ifeq ($(OS),Windows_NT)
	OCTANT_PLUGIN_DIR ?= ${LOCALAPPDATA}/octant/plugins
else
	OCTANT_PLUGIN_DIR ?= ${HOME}/.config/octant/plugins
endif

.PHONY: build
build:
	@go build -o bin/$(PLUGIN_NAME) cmd/octant-helm/main.go
	@mkdir -p $(OCTANT_PLUGIN_DIR)/
	@cp bin/$(PLUGIN_NAME) $(OCTANT_PLUGIN_DIR)/$(PLUGIN_NAME)

.PHONY: dev
dev:
	scripts/dev.sh
