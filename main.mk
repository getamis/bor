DOCKER_REPOSITORY := quay.io/amis
DOCKER_IMAGE := $(DOCKER_REPOSITORY)/indexer-bor
ifeq ($(DOCKER_IMAGE_TAG),)
DOCKER_IMAGE_TAG := $(shell git rev-parse --short HEAD 2> /dev/null)
endif
