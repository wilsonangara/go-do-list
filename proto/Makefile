SHELL := /usr/bin/env bash

GO_VERSION = 1.19
PROTOC_VERSION = 3.19.4
PROTOC_GEN_GO_VERSION = 1.27.1
PROTOC_GEN_TWIRP_VERSION = 8.1.3+incompatible

MAKEFILE_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))

PROTO_DIR := $(MAKEFILE_DIR)/pb

PROTO_FILES = $(shell find $(PROTO_DIR) -type f -name '*.proto')

GO_OUT_DIR := ./

UID := $(shell id -u)
GID := $(shell id -g)

.PHONY: gen
gen: 
	$(MAKE) build
	docker run \
		--user $(UID):$(GID) \
		--volume ${PWD}:${PWD} \
		--workdir ${PWD} \
		--rm go-do-list-proto \
			/bin/bash -c " \
				protoc \
					-I=$(PROTO_DIR) \
					--go_out=. --go_opt=paths=source_relative \
					--twirp_out=. --twirp_opt=paths=source_relative \
					$(PROTO_FILES); \
			"

.PHONY: build
build: 
	docker build \
		--build-arg GO_VERSION=${GO_VERSION} \
		--build-arg PROTOC_VERSION=${PROTOC_VERSION} \
		--build-arg PROTOC_GEN_GO_VERSION=${PROTOC_GEN_GO_VERSION} \
		--build-arg PROTOC_GEN_TWIRP_VERSION=${PROTOC_GEN_TWIRP_VERSION} \
		-t go-do-list-proto .
