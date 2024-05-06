.PHONY: build test docker-build docker-push all deploy

ANSIBLE_IMAGE_NAME=keyscome/ansible
ANSIBLE_VERSION=latest

build:
	go build -o bin/blossom 

test:
	go test ./...

all: build test

docker-build-ansible:
	docker build -t $(ANSIBLE_IMAGE_NAME):$(ANSIBLE_VERSION) builds -f builds/Dockerfile.ansible

docker-push-ansible:
	docker push $(ANSIBLE_IMAGE_NAME):$(ANSIBLE_VERSION)

deploy-ansible: docker-build-ansible docker-push-ansible
