SHELL := /usr/bin/env bash -o pipefail
GOPKG ?= github.com/MrEhbr/media-normalizer
DOCKER_IMAGE ?=	mrehbr/media-normalizer
GOBINS ?= cmd/media-normalizer
GO_APP ?= media-normalizer

include rules.mk
