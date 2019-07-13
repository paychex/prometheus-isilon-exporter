# Get the location of this makefile.
ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
TARGET_BINARY := prometheus-isilon-exporter
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
RELEASE?=$(shell git describe --abbrev=4 --dirty --always --tags)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
PROJECT_NAME?=github.com/paychex/prometheus-isilon-exporter

all: build

build:
	GO111MODULE=on go build -o bin/${TARGET_BINARY} \
		-ldflags="-X main.Commit=${COMMIT} \
		-X main.BuildTime=${BUILD_TIME} \
		-X main.Release=${RELEASE}" \
		./cmd